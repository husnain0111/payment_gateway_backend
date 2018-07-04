package userinfo

import (
	"db"
	"fmt"
	"shared"
	"time"

	"github.com/SermoDigital/jose/jws"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func GetUserInfo(user shared.UserSignup) (response shared.Response) {
	exception.Try(func() {

		response = ValidateInput1(user)
		if response.Success {
			var userinfo shared.UserSignup
			bucket := db.GetDbConnection(shared.BUCKET)

			_, err := bucket.Get("u:"+user.Email, &userinfo)
			if err != nil {
				response = shared.ReturnMessage(false, "Email not found", "E404")
				return
			}
			if userinfo.IsActive == false {
				response = shared.ReturnMessage(false, "Email Not Verified", "IU")
				return
			}

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
				if userinfo.FirstName != "" {
					result["first_name"] = userinfo.FirstName
				}
				if userinfo.LastName != "" {
					result["last_name"] = userinfo.LastName
				}
				if userinfo.Country != "" {
					result["country"] = userinfo.Country
				}
				if userinfo.PhoneNumber != "" {
					result["phone_number"] = userinfo.PhoneNumber
				}
				if userinfo.Gender != "" {
					result["gender"] = userinfo.Gender
				}
				if userinfo.Position != "" {
					result["position"] = userinfo.Position
				}
				if userinfo.CompanyName != "" {
					result["companyname"] = userinfo.CompanyName
				}
				if userinfo.CompanyWebsite != "" {
					result["companywebsite"] = userinfo.CompanyWebsite
				}
				if userinfo.AboutCompany != "" {
					result["aboutcompany"] = userinfo.AboutCompany
				}
				if userinfo.State != "" {
					result["state"] = userinfo.State
				}
				if userinfo.City != "" {
					result["city"] = userinfo.City
				}
				if userinfo.Skype != "" {
					result["skype"] = userinfo.Skype
				}
				if userinfo.Nickname != "" {
					result["nickname"] = userinfo.Nickname
				}

				result["token"] = token
				//fmt.Println(token)
				response.Message = result
				response.Success = true

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
func ValidateInput1(user shared.UserSignup) shared.Response {
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
