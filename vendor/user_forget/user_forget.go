package user_forget

import (
	"db"
	"encoding/json"
	"fmt"
	"shared"
	"strings"
	"time"

	"github.com/SermoDigital/jose/jws"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func ForgetPasswordGenerateToken(user shared.ForgetPasswordRequest) (response shared.Response) {
	exception.Try(func() {

		response = ValidateInput(user)
		if response.Success {
			var userinfo shared.UserSignup
			bucket := db.GetDbConnection(shared.BUCKET)
			_, err := bucket.Get("u:"+user.Email, &userinfo)
			if err != nil {
				response = shared.ReturnMessage(true, "Email Sent", "E404")
				return
			}
			expires := time.Now().Add(time.Duration(30) * time.Minute)
			claims := jws.Claims{}
			claims.SetExpiration(expires)
			claims.SetIssuedAt(time.Now())
			claims.Set("email", user.Email)
			claims.Set("action", "resetpassword")
			payload, _ := claims.MarshalJSON()
			sharedKey := []byte{54, 184, 135, 160, 180, 228, 101, 19, 240, 16, 110, 220, 13, 199, 15, 131}
			token, err := jose.Encrypt(string(payload), jose.A128KW, jose.A128GCM, sharedKey)
			//fmt.Println(token)
			if err != nil {
				response.Message = "Unable to generate token"
				response.Logs = append(response.Logs, err)
			}

			response = shared.SendEmail(user.Email,
				`Hi `+userinfo.Username+`,<br/><br/> 

				This email is to notify you to change password for your account at Blockipay.com<br/><br/><br/>
				
				To change your password please <a href="http://172.25.33.37:1323/resetpassword?token=`+token+`">Click Here</a> <br/><br/><br/>
				Best Regards,<br/><br/>
				
				
				If you did not perform this action, immediately click here to terminate access to your account and then a new generated password will be sent to your email. <br/><br/>
				
				Once you get your new password, login and change it to something you can remember on your profile page. <br/><br/>
				
				
				Blockipay.com<br/>
				Bitcoin Payment Gateway / Url Monetiser Online <br/>
				`,
				"Successful Login")

			response = shared.ReturnMessage(true, "Email Sent", "1")
			return

		}
	}).CatchAll(func(e interface{}) {
		raven.CaptureErrorAndWait(e.(error), map[string]string{"error": "Forget  Password Failed"})
		fmt.Println(e.(error))
		response = shared.ReturnMessage(false, "Something goes wrong", "SW")
	}).Finally(func() {
	})

	return response
}

func ForegetPasswordUpdate(user shared.ForgetPasswordRequest) (response shared.Response) {

	exception.Try(func() {

		response = ValidateInputWithToken(user)
		if response.Success {
			password := []byte(user.Password)
			hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

			sharedKey := []byte{54, 184, 135, 160, 180, 228, 101, 19, 240, 16, 110, 220, 13, 199, 15, 131}
			var pay shared.Payload
			payload, _, _ := jose.Decode(user.Token, sharedKey)
			err := json.NewDecoder(strings.NewReader(string(payload))).Decode(&pay)
			if err != nil {
				err2 := errors.WithStack(err)
				exception.Throw(fmt.Errorf("%+v", err2))

			}
			timeIAT := time.Now()
			fmt.Println(timeIAT.Unix())
			timeExp := time.Unix(pay.Exp, 0)
			timeex := timeIAT.Before(timeExp)
			if timeex == true {
				if pay.Action == "resetpassword" {
					bucket := db.GetDbConnection(shared.BUCKET)
					_, err = bucket.MutateIn("u:"+pay.Email, 0, 0).
						Replace("password", string(hashedPassword)).Execute()
					if err != nil {
						err2 := errors.WithStack(err)
						exception.Throw(fmt.Errorf("%+v", err2))

					}
					response = shared.ReturnMessage(true, "Password Updated", "1")
					response = shared.SendEmail(user.Email, `Your Password Updated`,
						"Password change successfuly")
					return
				} else {
					response = shared.ReturnMessage(false, "Invalid Token ", "TINV")
				}
			} else {
				response = shared.ReturnMessage(false, "Token Expire", "TE")
			}

		}
	}).CatchAll(func(e interface{}) {
		raven.CaptureErrorAndWait(e.(error), map[string]string{"error": "Fail While Forget Password"})
		fmt.Println(e.(error))
		response = shared.ReturnMessage(false, "Something goes wrong", "SW")
	}).Finally(func() {
	})

	return response
}
func ValidateInput(user shared.ForgetPasswordRequest) shared.Response {
	var response shared.Response

	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "ER"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if !validEmail {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "EINV"}
	}
	return shared.Response{Success: true}
}

func ValidateInputWithToken(user shared.ForgetPasswordRequest) shared.Response {
	var response shared.Response

	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "ER"); !response.Success {
		return response
	}

	if response = shared.ValidateRequiredWithMessage(user.Password, "Password Required", "ER"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(user.Token, "Token Required", "ER"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if !validEmail {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "EINV"}
	}
	return shared.Response{Success: true}
}
