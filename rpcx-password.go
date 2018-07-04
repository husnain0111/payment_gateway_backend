package main

import (
	"shared"
	"user_forget"

	raven "github.com/getsentry/raven-go"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func main() {

	shared.RpcxListener("bkForgetPassword", new(user_forget.UserForgetPassword), "127.0.0.1:5003")

}
