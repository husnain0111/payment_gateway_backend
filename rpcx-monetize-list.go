package main

import (
	"monetiser"
	"shared"

	raven "github.com/getsentry/raven-go"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func main() {
	//reflect.TypeOf(new(user_signup.UserSignup))
	shared.RpcxListener("bkMonetizelist", new(monetiser.MonTRequest), "127.0.0.1:5006")
}
