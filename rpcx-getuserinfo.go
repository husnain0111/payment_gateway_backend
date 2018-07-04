package main

import (
	"shared"
	"userinfo"

	raven "github.com/getsentry/raven-go"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func main() {
	//reflect.TypeOf(new(user_signup.UserSignup))
	shared.RpcxListener("bkUserInfoGet", new(userinfo.UserInfo), "127.0.0.1:5008")
}
