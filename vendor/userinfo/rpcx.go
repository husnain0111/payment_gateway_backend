package userinfo

import (
	"context"
	"encoding/json"
	"shared"
	"strings"
)

type UserInfo shared.Response

func (l *UserInfo) UserInfo(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request shared.UserSignup
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = AddUserInfo(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
func (l *UserInfo) GetUserInfo(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request shared.UserSignup
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = GetUserInfo(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
