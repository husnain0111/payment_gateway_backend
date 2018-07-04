package user_signup

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
	gocb "gopkg.in/couchbase/gocb.v1"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func AddUserGenerateToken(user shared.UserSignup) shared.Response {
	var response shared.Response
	exception.Try(func() {
		response = ValidateInput(user)
		if response.Success {
			user.Email = strings.ToLower(user.Email)
			bucket := db.GetDbConnection(shared.BUCKET)
			isExist := IsEmailUsernameExist(user.Email, bucket)
			if isExist {
				response.Success = false
				response.Message = "Email/Username already exists"
				response.Code = "EF"

				return
			}
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashedPassword)
			user.IsActive = false

			expires := time.Now().Add(time.Duration(30) * time.Minute)
			claims := jws.Claims{}
			claims.SetExpiration(expires)
			claims.SetIssuedAt(time.Now())
			claims.Set("email", user.Email)
			claims.Set("action", "signup")
			payload, _ := claims.MarshalJSON()
			sharedKey := []byte{54, 184, 135, 160, 180, 228, 101, 19, 240, 16, 110, 220, 13, 199, 15, 131}
			token, err := jose.Encrypt(string(payload), jose.A128KW, jose.A128GCM, sharedKey)
			fmt.Println(token)
			if err != nil {
				response.Message = "Unable to generate token"
				response.Logs = append(response.Logs, err)
			}

			response = shared.SendEmail(user.Email,
				`Dear `+user.Username+`,<br/><br/>
				Thank you for your registration on Bitcoin Payment Gateway / Url Monetiser Online.<br/>
				To complete your Registration please <a href="http://172.25.33.37:1323/resetpassword?token=`+token+`">Click Here</a> <br/><br/><br/>
				Best Regards,<br/><br/> 

				Blockipay.com Team <br/>
				Bitcoin Payment Gateway / Url Monetiser Online<br/><br/> 
				Website: Blockipay.com<br/>
				`,
				"Blockipay Confirmation of Email Address")
			bucket1 := db.GetDbConnection(shared.BUCKET)
			_, err1 := bucket1.Upsert("u:"+user.Email, user, 0)

			if err1 != nil {
				err2 := errors.WithStack(err1)
				response = shared.ReturnMessage(false, "Unable to Insert", "DIN")
				response.Logs = append(response.Logs, err2)
				exception.Throw(fmt.Errorf("%+v", err2))

			} else {
				response.Success = true
				response.Message = "Registration Email Send Successfully"
				response.Code = "1"

			}
		}
	}).CatchAll(func(e interface{}) {
		raven.CaptureErrorAndWait(e.(error), map[string]string{"error": "Fail to Signup"})
		fmt.Println(e.(error))
		response = shared.ReturnMessage(false, "Something goes wrong", "SW")
	}).Finally(func() {

	})
	return response
}
func AddUser(user shared.UserSignup) (response shared.Response) {

	exception.Try(func() {

		response = ValidateInputWithToken(user)
		if response.Success {

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
				if pay.Action == "signup" {
					bucket := db.GetDbConnection(shared.BUCKET)
					_, err = bucket.MutateIn("u:"+pay.Email, 0, 0).
						Replace("is_active", true).Execute()
					if err != nil {
						err2 := errors.WithStack(err)
						exception.Throw(fmt.Errorf("%+v", err2))

					}
					response = shared.ReturnMessage(true, "User Successfully Added", "1")
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

func IsEmailUsernameExist(email string, bucket *gocb.Bucket) bool {
	var user shared.UserSignup
	_, err := bucket.Get("u:"+email, &user)
	if err != nil {
		return false
	}
	return true
}
func ValidateInput(user shared.UserSignup) shared.Response {
	var response shared.Response
	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "INV"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if validEmail {

		if response = shared.ValidateRequiredWithMessage(user.Username, "Username Required", "INV"); !response.Success {
			return response
		}
		if response = shared.ValidateRequiredWithMessage(user.Password, "Password Required", "INV"); !response.Success {
			return response
		}

	} else {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "INE"}
	}
	return shared.Response{Success: true}
}

func ValidateInputWithToken(user shared.UserSignup) shared.Response {
	var response shared.Response

	if response = shared.ValidateRequiredWithMessage(user.Token, "Token Required", "ER"); !response.Success {
		return response
	}

	return shared.Response{Success: true}
}
