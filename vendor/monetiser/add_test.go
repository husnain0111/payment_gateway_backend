package monetiser

import (
	"fmt"
	"testing"
)

func TestAddMonetiserUrl(t *testing.T) {
	result := CreateMonetiserUrl(MonetiserRequest{
		PrivateURL:    "http://google.com",
		PrivateText:   "this is sample text",
		Price:         0.1,
		WalletAddress: "XNKANKANKJANSJAS",
		Currency:      "BTC",
		PublicTitle:   "Sample test",
	})
	/*result := GetMonetizeInfo(MonetiserRequest{
		ID:   "OAKLvJduBcjRPQz5cJx19wlEjJqlWQzHnnrjfYVobnF0dIIukT",
		Mode: "test",
	})*/
	fmt.Println(result)
}
