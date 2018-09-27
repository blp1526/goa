// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// chatter gRPC server types
//
// Command:
// $ goa gen goa.design/goa/examples/streaming/design -o
// $(GOPATH)/src/goa.design/goa/examples/streaming

package server

import (
	chattersvc "goa.design/goa/examples/streaming/gen/chatter"
	chatterpb "goa.design/goa/examples/streaming/gen/grpc/chatter"
)

// NewLoginPayload builds the payload of the "login" endpoint of the "chatter"
// service from the gRPC request type.
func NewLoginPayload(message *chatterpb.LoginRequest, user string, password string) *chattersvc.LoginPayload {
	v := &chattersvc.LoginPayload{}
	v.User = user
	v.Password = password
	return v
}

// NewLoginResponse builds the gRPC response type from the result of the
// "login" endpoint of the "chatter" service.
func NewLoginResponse(res string) *chatterpb.LoginResponse {
	v := &chatterpb.LoginResponse{}
	v.Field = res
	return v
}

// NewEchoerPayload builds the payload of the "echoer" endpoint of the
// "chatter" service from the gRPC request type.
func NewEchoerPayload(token string) *chattersvc.EchoerPayload {
	v := &chattersvc.EchoerPayload{}
	v.Token = token
	return v
}

func NewEchoerResponse(res string) *chatterpb.EchoerResponse {
	v := &chatterpb.EchoerResponse{}
	v.Field = res
	return v
}

func NewEchoerStreamingRequest(v *chatterpb.EchoerStreamingRequest) string {
	p := v.Field
	return p
}

// NewListenerPayload builds the payload of the "listener" endpoint of the
// "chatter" service from the gRPC request type.
func NewListenerPayload(token string) *chattersvc.ListenerPayload {
	v := &chattersvc.ListenerPayload{}
	v.Token = token
	return v
}

func NewListenerStreamingRequest(v *chatterpb.ListenerStreamingRequest) string {
	p := v.Field
	return p
}

// NewSummaryPayload builds the payload of the "summary" endpoint of the
// "chatter" service from the gRPC request type.
func NewSummaryPayload(token string) *chattersvc.SummaryPayload {
	v := &chattersvc.SummaryPayload{}
	v.Token = token
	return v
}

func NewChatSummaryCollection(res chattersvc.ChatSummaryCollection) *chatterpb.ChatSummaryCollection {
	v := &chatterpb.ChatSummaryCollection{}
	v.Field = make([]*chatterpb.ChatSummary, len(res))
	for i, val := range res {
		v.Field[i] = &chatterpb.ChatSummary{
			Message_: val.Message,
		}
		if val.Length != nil {
			v.Field[i].Length = int32(*val.Length)
		}
		if val.SentAt != nil {
			v.Field[i].SentAt = *val.SentAt
		}
	}
	return v
}

func NewSummaryStreamingRequest(v *chatterpb.SummaryStreamingRequest) string {
	p := v.Field
	return p
}

// NewHistoryPayload builds the payload of the "history" endpoint of the
// "chatter" service from the gRPC request type.
func NewHistoryPayload(message *chatterpb.HistoryRequest, view string, token string) *chattersvc.HistoryPayload {
	v := &chattersvc.HistoryPayload{}
	v.View = &view
	v.Token = token
	return v
}

func NewHistoryResponse(res *chattersvc.ChatSummary) *chatterpb.HistoryResponse {
	v := &chatterpb.HistoryResponse{
		Message_: res.Message,
	}
	if res.Length != nil {
		v.Length = int32(*res.Length)
	}
	if res.SentAt != nil {
		v.SentAt = *res.SentAt
	}
	return v
}