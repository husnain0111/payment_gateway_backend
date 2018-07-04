package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"

	"shared"

	raven "github.com/getsentry/raven-go"
	"github.com/smallnest/rpcx/client"

	"sync"
	"time"
)

func init() {

	raven.SetDSN(shared.Raven)

}

type Gateway struct {
	mu               sync.RWMutex
	xclients         map[string]client.XClient
	serviceDiscovery client.ServiceDiscovery
	FailMode         client.FailMode
	SelectMode       client.SelectMode
	Option           client.Option
}

func (g *Gateway) Makerequest(address string, servicePath string, method string, args shared.Request) *shared.Response {
	flag.Parse()

	option := client.DefaultOption
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	option.TLSConfig = conf

	var xc client.XClient
	g.Option = option
	//
	//	xc = client.NewXClient(service, client.Failtry, client.ConsistentHash, d, option)
	//	defer xc.Close()

	g.mu.Lock()
	xc, err := getXClient(g, servicePath, address)
	if err != nil {
		g.xclients[servicePath] = nil
		xc, err = getXClient(g, servicePath, address)
	}
	g.mu.Unlock()

	reply := &shared.Response{}
	call, err := xc.Go(context.Background(), method, args, reply, nil)
	if err != nil {
		raven.CaptureErrorAndWait(err.(error), map[string]string{"error": "Fail at Gateway" + servicePath})
		g.xclients[servicePath] = nil
		reply.Success = false
	} else {
		replyCall := <-call.Done
		if replyCall.Error != nil {
			reply.Success = false
			g.xclients[servicePath] = nil
			return reply
		}

	}
	return reply
}
func getXClient(g *Gateway, servicePath string, address string) (xc client.XClient, err error) {
	defer func() {
		if e := recover(); e != nil {
			if ee, ok := e.(error); ok {
				err = ee
				return
			}

			err = fmt.Errorf("failed to get xclient: %v", e)
		}
	}()

	if g.xclients[servicePath] == nil {
		d := client.NewPeer2PeerDiscovery("tcp@"+address, "")
		g.xclients[servicePath] = client.NewXClient(servicePath, g.FailMode, g.SelectMode, d, g.Option)
	}
	xc = g.xclients[servicePath]

	return xc, err
}
