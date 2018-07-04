package userinfo

import (
	"db"
	"fmt"
	"shared"

	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func AddUserInfo(user shared.UserSignup) shared.Response {
	var response shared.Response
	var userinfo shared.UserSignup
	exception.Try(func() {
		response = ValidateInput(user)
		if response.Success {

			bucket := db.GetDbConnection(shared.BUCKET)

			_, err := bucket.Get("u:"+user.Email, &userinfo)
			if err != nil {
				response = shared.ReturnMessage(false, "Email not found", "E404")

				// err2 := errors.WithStack(err)
				// response.Logs = append(response.Logs, err2)
				// exception.Throw(fmt.Errorf("%+v", err2))
				return

			} else {

				_, err = bucket.MutateIn("u:"+user.Email, 0, 0).
					Upsert("first_name", user.FirstName, true).
					Upsert("last_name", user.LastName, true).
					Upsert("gender", user.Gender, true).
					Upsert("position", user.Position, true).
					Upsert("companyname", user.CompanyName, true).
					Upsert("companywebsite", user.CompanyWebsite, true).
					Upsert("aboutcompany", user.AboutCompany, true).
					Upsert("country", user.Country, true).
					Upsert("state", user.State, true).
					Upsert("city", user.City, true).
					Upsert("phone_number", user.PhoneNumber, true).
					Upsert("skype", user.Skype, true).
					Upsert("nickname", user.Nickname, true).
					Execute()

				//fmt.Println("hello")

				if err != nil {
					err2 := errors.WithStack(err)
					exception.Throw(fmt.Errorf("%+v", err2))

				}
				logrus.WithFields(logrus.Fields{"Email": user.Email}).Info("user info added")
				response = shared.ReturnMessage(true, "User Info Added", "1")
				return

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

func ValidateInput(user shared.UserSignup) shared.Response {
	var response shared.Response
	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "INV"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if validEmail {

	} else {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "INE"}
	}
	return shared.Response{Success: true}
}
