package user_forget

import (
	"fmt"
	"shared"
	"testing"
)

/*
func TestForgetPassword(t *testing.T) {
	result := ForgetPasswordGenerateToken(shared.ForgetPasswordRequest{
		Email: "asimizb@gmail.com",
	})
	fmt.Println(result)
}*/

func TestForgetPasswordWithToken(t *testing.T) {
	result := ForegetPasswordUpdate(shared.ForgetPasswordRequestWithToken{
		Email:    "asimizb@gmail.com",
		Password: "asim1234",
		Token:    "eyJhbGciOiJBMTI4S1ciLCJlbmMiOiJBMTI4R0NNIn0.giAwwQhQ8uvu5ssW03JOhK1xS1sl1md1.U9fChX8K15NpUOrP.om6bWo1UqmhU_kIoKdyc0OUGJeG2soJJaAU2HekYaQjIrf0Sx4hwD7z_K4Np_Gtq3k_HIW4xQc4SLCPXBn05KEVUqo9G0Bt_kXdfIkcqmqhsEuNo1J7b7g.I6dFp9vWkfUr4ZRBvRP2lw",
	})
	fmt.Println(result)
}
