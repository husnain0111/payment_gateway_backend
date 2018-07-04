package user_login

import (
	"fmt"
	"shared"
	"testing"
)

func TestUserLogin(t *testing.T) {
	result := LoginUser(shared.UserLogin{
		Email:     "mahar.husnain@yahoo.com",
		Password:  "hello123",
		IP:        "182.180.59.151",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36",
	})
	fmt.Println(result)
}
