package monetiser

import (
	"context"
	"encoding/json"
	"shared"
	"strings"
)

type MonTRequest MonetiserRequest

func (l *MonTRequest) MontizeAdd(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request MonetiserRequest
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = CreateMonetiserUrl(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
func (l *MonTRequest) MontizeGet(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request MonetiserRequest
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = GetMonetizeInfo(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
func (l *MonTRequest) MontizeList(ctx context.Context, args *shared.Request, reply *shared.Response) error {
	var response shared.Response
	var request ListMontiserRequest
	err := json.NewDecoder(strings.NewReader(string(args.Data))).Decode(&request)
	if err != nil {
		reply.Success = false
		reply.Message = "Invalid Json Data"
	} else {
		response = ListAllUrls(request)
		reply.Success = response.Success
		reply.Code = response.Code
		reply.Message = response.Message
		reply.Logs = response.Logs

	}
	return nil
}
