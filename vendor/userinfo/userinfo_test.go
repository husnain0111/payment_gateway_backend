package userinfo

import (
	"fmt"
	"shared"
	"testing"
)

func TestUserInfo(t *testing.T) {
	result := GetUserInfo(shared.UserSignup{
		// FirstName:      "mahar",
		// LastName:       "husnain",
		Email: "farina0@mailinator.com",
		// Country:        "pakistan",
		// PhoneNumber:    "1234567890",
		// Gender:         "male",
		// Position:       "engineer",
		// CompanyName:    "BCD",
		// CompanyWebsite: "bcd.io",
		// AboutCompany:   "blockchain company",
		// State:          "islamabad",
		// City:           "pakistan",
		// Skype:          "maharhusnain",
		// Nickname:       "mahar",
	})
	fmt.Println(result)
	//bucket := db.GetDbConnection(shared.ICO_BUCKET)
	//aa := DeleteUser("asimizb@gmail.com", bucket)
	//fmt.Println(aa)
}
