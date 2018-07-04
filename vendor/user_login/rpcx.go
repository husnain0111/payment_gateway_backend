package user_login

import (
	"context"
	"encoding/json"
	"shared"
	"strings"
)

type UserSignin shared.Response

func (l *UserSignin) UserSignin(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request shared.UserLogin
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = LoginUser(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
