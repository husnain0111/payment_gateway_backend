package user_forget

import (
	"context"
	"encoding/json"
	"shared"
	"strings"
)

type UserForgetPassword shared.Response

func (l *UserForgetPassword) ForgetPasswordUser(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request shared.ForgetPasswordRequest
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		if request.Token == "" {
			response = ForgetPasswordGenerateToken(request)
		} else {
			response = ForegetPasswordUpdate(request)
		}
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
