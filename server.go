package main

import (
	"fmt"
	"net/http"
	"shared"

	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/smallnest/rpcx/client"
)

func init() {

	raven.SetDSN(shared.Raven)

}

type Services struct {
	ServiceName  string
	FunctionName string
	Address      string
}

var gateway *Gateway
var services map[string]Services

func main() {
	gateway = &Gateway{
		xclients: make(map[string]client.XClient),
	}
	services = map[string]Services{}
	services["signup"] = Services{ServiceName: "bkUserSignup", FunctionName: "SignupUser", Address: "127.0.0.1:5001"}
	services["login"] = Services{ServiceName: "bkUserLogin", FunctionName: "UserSignin", Address: "127.0.0.1:5002"}
	services["forget-password"] = Services{ServiceName: "bkForgetPassword", FunctionName: "ForgetPasswordUser", Address: "127.0.0.1:5003"}
	services["monetiser-get"] = Services{ServiceName: "bkMonetizeget", FunctionName: "MontizeGet", Address: "127.0.0.1:5004"}
	services["monetiser-add"] = Services{ServiceName: "bkMonetizeadd", FunctionName: "MontizeAdd", Address: "127.0.0.1:5005"}
	services["monetiser-list"] = Services{ServiceName: "bkMonetizelist", FunctionName: "MontizeList", Address: "127.0.0.1:5006"}
	//mahar
	services["userinfo"] = Services{ServiceName: "bkUserInfo", FunctionName: "UserInfo", Address: "127.0.0.1:5007"}
	services["userinfo-get"] = Services{ServiceName: "bkUserInfoGet", FunctionName: "GetUserInfo", Address: "127.0.0.1:5008"}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	//services = map[string]Services{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/", handler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func handler(c echo.Context) error {
	u := new(shared.Request)
	if err := c.Bind(u); err != nil {
		//return
	}
	fmt.Println(u.Action)
	fmt.Println(services[u.Action].ServiceName)
	if services[u.Action].ServiceName != "" {
		var result *shared.Response
		exception.Try(func() {
			// something that might go wrong

			result = gateway.Makerequest(services[u.Action].Address, services[u.Action].ServiceName, services[u.Action].FunctionName, shared.Request{Data: u.Data, Action: u.Action})

		}).Catch(func(e error) {
			raven.CaptureErrorAndWait(e.(error), map[string]string{"error": "Fail at Gateway"})

		}).Go()
		return c.JSON(http.StatusOK, shared.Response{
			Success: result.Success,
			Message: result.Message,
			Logs:    result.Logs,
			Code:    result.Code,
		})

	} else {
		return c.JSON(http.StatusOK, shared.Response{Success: false, Message: "Invalid method"})
	}

}
