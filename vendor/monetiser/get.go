package monetiser

import (
	"db"
	"fmt"
	"shared"
	"time"

	gocb "gopkg.in/couchbase/gocb.v1"

	"github.com/ftloc/exception"
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

func init() {

	raven.SetDSN(shared.Raven)

}
func ValidateInputID(data MonetiserRequest) (response shared.Response) {
	if response = shared.ValidateRequiredWithMessage(data.ID, "ID Required", "INV"); !response.Success {
		return response
	}

	return response
}

func GetMonetizeInfo(data MonetiserRequest) (response shared.Response) {
	exception.Try(func() {
		response = ValidateInputID(data)
		if response.Success {

			bucket := db.GetDbConnection(shared.BUCKET)
			_, err := bucket.Get("mont:"+data.ID, &data)

			if err != nil {
				err2 := errors.WithStack(err)
				response = shared.ReturnMessage(false, "Unable to Find ID", "DIN")
				response.Logs = append(response.Logs, err2)
				exception.Throw(fmt.Errorf("%+v", err2))

			} else if data.PrivateURL != "" {
				data.Address = GetAddress(data, bucket)
				response.Success = true
				response.Message = data
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

func GetAddress(contributionAddress MonetiserRequest, bucket *gocb.Bucket) string {

	query := gocb.NewN1qlQuery("SELECT * FROM " + shared.BUCKET + "  use index (`idx_address`)   where inuse=false AND currency=\"" + contributionAddress.Currency + "\" limit 1")
	if contributionAddress.Mode == "test" {
		query = gocb.NewN1qlQuery("SELECT * FROM " + shared.BUCKET + "  use index (`idx_address`)   where currency=\"" + contributionAddress.Currency + "\" limit 1")
	}
	contAddress := ""
	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{})
	shared.ThrowErrorIFNotNilWithMessage(err, "Error in getting the address")
	var row interface{}
	rows.Close()
	if rows.Metrics().ResultCount > 0 {
		for rows.Next(&row) {
		}
		rowBucket := row.(map[string]interface{})
		addressDta := rowBucket[shared.BUCKET].(map[string]interface{})
		var row interface{}
		addresskey := (contributionAddress.Currency) + addressDta["address"].(string)
		contAddress = addressDta["address"].(string)
		cas, _ := bucket.GetAndLock(addresskey,
			uint32(time.Unix(100, 10000).Unix()), &row)

		if cas != 0 {
			_, err := bucket.Replace(addresskey,
				shared.Address{
					Address:  addressDta["address"].(string),
					Inuse:    true,
					Currency: contributionAddress.Currency,
					ID:       contributionAddress.ID,
				}, cas, 0)

			if err != nil {

				return GetAddress(contributionAddress, bucket)
			}
			return addressDta["address"].(string)

		} else {
			fmt.Println("Unable to get the address")
		}

		return GetAddress(contributionAddress, bucket)

	} else {
		return contAddress
	}

}
