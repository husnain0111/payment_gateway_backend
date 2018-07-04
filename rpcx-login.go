package main

import (
	"shared"
	"user_login"

	raven "github.com/getsentry/raven-go"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func main() {

	shared.RpcxListener("bkUserLogin", new(user_login.UserSignin), "127.0.0.1:5002")

}
