package user_signup

import (
	"context"
	"encoding/json"
	"shared"
	"strings"
)

type UserSignup shared.Response

func (l *UserSignup) SignupUser(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request shared.UserSignup
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		if request.Token == "" {
			response = AddUserGenerateToken(request)
		} else {
			response = AddUser(request)
		}
		//response = AddUSer(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
