package user_login

import (
	"db"
	"fmt"
	"shared"
	"time"

	"github.com/SermoDigital/jose/jws"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/mssola/user_agent"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func LoginUser(user shared.UserLogin) (response shared.Response) {
	exception.Try(func() {

		response = ValidateInput(user)
		if response.Success {
			var userinfo shared.UserSignup
			bucket := db.GetDbConnection(shared.BUCKET)

			_, err := bucket.Get("u:"+user.Email, &userinfo)
			if err != nil {
				response = shared.ReturnMessage(false, "Email not found", "E404")
				return
			}
			err = bcrypt.CompareHashAndPassword([]byte(userinfo.Password), []byte(user.Password))
			if err != nil {
				response = shared.ReturnMessage(false, "Invalid User", "IU")
				return
			}
			if userinfo.IsActive == false {
				response = shared.ReturnMessage(false, "Email Not Verified", "IU")
				return
			}

			geo, err := shared.GetIpGeoLocation(user.IP)
			if err != nil {
				fmt.Println(err)
			}

			var geolocation string
			geolocation = geo.CountryName + ", " + geo.RegionName + ", " + geo.City
			ip := user.IP
			fmt.Println(ip, geolocation)

			ua := user_agent.New(user.UserAgent)
			os := ua.OS()
			fmt.Printf("%v\n", os) // => "Linux x86_64"
			name, version := ua.Browser()
			browser := name + " " + version
			fmt.Printf("%v\n", browser) // => "Chrome"

			expires := time.Now().Add(time.Duration(24) * time.Hour)

			claims := jws.Claims{}
			claims.SetExpiration(expires)
			claims.SetIssuedAt(time.Now())
			claims.Set("email", user.Email)

			payload, _ := claims.MarshalJSON()
			sharedKey := []byte{54, 104, 215, 160, 189, 218, 171, 9, 214, 26, 111, 221, 170, 199, 125, 71}
			token, err := jose.Encrypt(string(payload), jose.A128KW, jose.A128GCM, sharedKey)
			if err == nil {
				result := map[string]interface{}{}
				result["username"] = userinfo.Username
				result["email"] = userinfo.Email
				result["token"] = token
				//fmt.Println(token)
				response.Message = result
				response.Success = true
				response1 := shared.SendEmail(user.Email,
					`Hi `+userinfo.Username+`,<br/><br/> 

					This email is to notify you of a successful login on your account at Blockipay.com<br/><br/><br/>
					
					
					
					<b>Login:</b>      	`+userinfo.Email+`<br/>    
					<b>IP Address:</b>     `+user.IP+`   <br/>   
					<b>Location:</b>      	`+geolocation+`    <br/>  
					<b>Browser:</b>      	`+browser+`      <br/>
					<b>OS:</b>      	`+os+`<br/>
					<b>Date:</b>      	`+time.Now().UTC().String()+`     <br/><br/><br/> 
					
					
					
					If you did not perform this login, immediately click here to terminate access to your account and then a new generated password will be sent to your email. <br/><br/>
					
					Once you get your new password, login and change it to something you can remember on your profile page. <br/><br/>
					
					
					Blockipay.com<br/>
					Bitcoin Payment Gateway / Url Monetiser Online <br/>
					`,
					"Successful Login")
				fmt.Println(response1)

			} else {
				err2 := errors.WithStack(err)

				exception.Throw(fmt.Errorf("%+v", err2))

			}
			//fmt.Println(userinfo)
		}

	}).CatchAll(func(e interface{}) {
		raven.CaptureErrorAndWait(e.(error), map[string]string{"error": "Fail to Signin"})
		fmt.Println(e.(error))
		response = shared.ReturnMessage(false, "Something goes wrong", "SW")
	}).Finally(func() {

	})

	return response
}
func ValidateInput(user shared.UserLogin) shared.Response {
	var response shared.Response

	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "ER"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(user.Password, "Password Required", "ER"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if !validEmail {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "EINV"}
	}
	return shared.Response{Success: true}
}
