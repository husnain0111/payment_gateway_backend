package monetiser

import (
	"db"
	"encoding/json"
	"fmt"
	"shared"

	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	gocb "gopkg.in/couchbase/gocb.v1"
)

func ListAllUrls(data ListMontiserRequest) (response shared.Response) {

	exception.Try(func() {
		response = ValidateInputListAll(data)
		if response.Success {
			email := shared.ValidateToken(data.Token)
			if email == "" {
				response = shared.ReturnMessage(false, "Token Expire, Logout", "TEX")
				return
			}
			if email == data.Email {
				bucket := db.GetDbConnection(shared.BUCKET)
				defer bucket.Close()
				query := gocb.NewN1qlQuery("SELECT * FROM " + shared.BUCKET + " WHERE private_url is not missing ")
				rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{})
				shared.ThrowErrorIFNotNil(err)
				rows.Close()
				isfound := rows.Metrics().ResultCount
				if isfound > 0 {
					var row map[string]interface{}
					monetizelinks := []MonetiserRequest{}
					for rows.Next(&row) {
						var monetiser MonetiserRequest
						txRow := row[shared.BUCKET].(map[string]interface{})
						jsonString, _ := json.Marshal(txRow)
						if err := json.Unmarshal(jsonString, &monetiser); err != nil {
							shared.ThrowErrorIFNotNilWithMessage(err, "montiserlist")
						}

						monetizelinks = append(monetizelinks, monetiser)
						response.Message = monetizelinks
						response.Success = true
						response.Code = "1"
					}
				} else {
					response = shared.ReturnMessage(false, "No Data", "ND")

				}
			} else {
				response = shared.ReturnMessage(false, "Forbidden", "403")
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

func ValidateInputListAll(user ListMontiserRequest) (response shared.Response) {

	if response = shared.ValidateRequiredWithMessage(user.Email, "Email Required", "INV"); !response.Success {
		return response
	}
	validEmail := shared.ValidateEmail(user.Email)
	if !validEmail {
		return shared.Response{Success: false, Message: "Invalid Email", Code: "INE"}
	}
	if response = shared.ValidateRequiredWithMessage(user.Token, "Token Required", "INV"); !response.Success {
		return response
	}
	return shared.Response{Success: true}
}
