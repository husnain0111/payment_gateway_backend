package main

import (
	"shared"
	"user_signup"

	raven "github.com/getsentry/raven-go"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func main() {
	//reflect.TypeOf(new(user_signup.UserSignup))
	shared.RpcxListener("bkUserSignup", new(user_signup.UserSignup), "127.0.0.1:5001")
}
