package user_signup

import (
	"fmt"
	"shared"
	"testing"
)

func TestUserAdd(t *testing.T) {
	result := AddUserGenerateToken(shared.UserSignup{
		Email:    "mahar.husnain@yahoo.com",
		Password: "hello123",
		Username: "mahar",
		Token:    "eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4R0NNIn0.Iv_GScd3Y4ha0V_lyDjzuQZgzPY1kkvc.FlH4O2RoMRaKtHgg.kKmdeGlYeye_cUR3HjDx2QaRJkhnwQtUG828qGuKTAJ5Tee3CMPZtrFFxV5wnf0YjvupVcqrDqY4lt53sVd5ZayzBDB5TFj5lko1EgyEzaJbhz911dJf.zfYY-mIfXMioP3fG5s2LVg",
	})
	fmt.Println(result)
	//bucket := db.GetDbConnection(shared.ICO_BUCKET)
	//aa := DeleteUser("asimizb@gmail.com", bucket)
	//fmt.Println(aa)
}
