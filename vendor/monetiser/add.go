package monetiser

import (
	"db"
	"fmt"
	"shared"

	gocb "gopkg.in/couchbase/gocb.v1"

	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/zemirco/uid"
)

func init() {

	raven.SetDSN(shared.Raven)

}

func CreateMonetiserUrl(data MonetiserRequest) (response shared.Response) {

	exception.Try(func() {
		response = ValidateInput(data)
		if response.Success {
			data.ID = uid.New(50)
			bucket := db.GetDbConnection(shared.BUCKET)
			isExist := true
			for isExist {
				data.ID = uid.New(50)
				isExist = IsIDExist(data.ID, bucket)

			}
			data.Address = GetAddress(data, bucket)
			_, err := bucket.Insert("mont:"+data.ID, data, 0)
			if err != nil {
				err2 := errors.WithStack(err)
				response = shared.ReturnMessage(false, "Unable to Insert", "DIN")
				response.Logs = append(response.Logs, err2)
				exception.Throw(fmt.Errorf("%+v", err2))

			} else {
				response.Success = true
				response.Message = `{"id":"` + data.ID + `"}`
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
func IsIDExist(uid string, bucket *gocb.Bucket) bool {
	var data MonetiserRequest
	_, err := bucket.Get("mont:"+uid, &data)
	if err != nil {
		return false
	}
	return true
}
func ValidateInput(data MonetiserRequest) (response shared.Response) {
	if response = shared.ValidateRequiredWithMessage(data.Currency, "Currency Required", "INV"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(data.PrivateURL, "Private URL Required", "INV"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(data.PublicTitle, "Public Title Required", "INV"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(data.Price, "Price Required", "INV"); !response.Success {
		return response
	}
	if response = shared.ValidateRequiredWithMessage(data.WalletAddress, "Wallet Address Required", "INV"); !response.Success {
		return response
	}

	return response
}
