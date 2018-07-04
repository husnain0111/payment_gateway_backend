package db

import (
	"fmt"
	shared "shared"
	"time"

	"github.com/ftloc/exception"
	"github.com/pkg/errors"
	gocb "gopkg.in/couchbase/gocb.v1"
)

func GetDbConnection(bucketname string) *gocb.Bucket {
	cluster, err := gocb.Connect(shared.DB_URL)
	if err != nil {
		fmt.Println("Error in url ", err)
	}

	cluster.SetConnectTimeout(5 * time.Second)
	cluster.SetServerConnectTimeout(5 * time.Second)
	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: shared.DB_USERNAME,
		Password: shared.DB_PASSWORD,
	})

	if err != nil {
		fmt.Println("Error in Authentication ", err)
	}
	bucket, err := cluster.OpenBucket(bucketname, "")
	if err != nil {
		err2 := errors.WithStack(err)
		exception.Throw(fmt.Errorf("%+v", err2))
	}
	//bucket.Manager("", "").CreatePrimaryIndex("", true, false)
	go func() {
		time.Sleep(5 * time.Second)
		bucket.Close()
	}()
	return bucket
}

func GetDbConnection2(bucketname string) *gocb.Bucket {
	cluster, err := gocb.Connect(shared.DB_URL)
	if err != nil {
		fmt.Println("Error in url ", err)
	}

	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: shared.DB_USERNAME,
		Password: shared.DB_PASSWORD,
	})

	if err != nil {
		fmt.Println("Error in Authentication ", err)
	}
	bucket, err := cluster.OpenBucket(bucketname, "")
	if err != nil {
		err2 := errors.WithStack(err)
		exception.Throw(fmt.Errorf("%+v", err2))
	}
	//bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	return bucket
}
