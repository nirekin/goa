// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// chatter gRPC server types
//
// Command:
// $ goa gen goa.design/goa/examples/chatter/design -o
// $(GOPATH)/src/goa.design/goa/examples/chatter

package server

import (
	goa "goa.design/goa"
	chattersvc "goa.design/goa/examples/chatter/gen/chatter"
	chattersvcviews "goa.design/goa/examples/chatter/gen/chatter/views"
	"goa.design/goa/examples/chatter/gen/grpc/chatter/pb"
)

// NewLoginPayload builds the payload of the "login" endpoint of the "chatter"
// service from the gRPC request type.
func NewLoginPayload(user string, password string) *chattersvc.LoginPayload {
	payload := &chattersvc.LoginPayload{}
	payload.User = user
	payload.Password = password
	return payload
}

// NewLoginResponse builds the gRPC response type from the result of the
// "login" endpoint of the "chatter" service.
func NewLoginResponse(result string) *pb.LoginResponse {
	message := &pb.LoginResponse{}
	message.Field = result
	return message
}

// NewEchoerPayload builds the payload of the "echoer" endpoint of the
// "chatter" service from the gRPC request type.
func NewEchoerPayload(token string) *chattersvc.EchoerPayload {
	payload := &chattersvc.EchoerPayload{}
	payload.Token = token
	return payload
}

func NewEchoerResponse(result string) *pb.EchoerResponse {
	v := &pb.EchoerResponse{}
	v.Field = result
	return v
}

func NewEchoerStreamingRequest(v *pb.EchoerStreamingRequest) string {
	spayload := v.Field
	return spayload
}

// NewListenerPayload builds the payload of the "listener" endpoint of the
// "chatter" service from the gRPC request type.
func NewListenerPayload(token string) *chattersvc.ListenerPayload {
	payload := &chattersvc.ListenerPayload{}
	payload.Token = token
	return payload
}

func NewListenerStreamingRequest(v *pb.ListenerStreamingRequest) string {
	spayload := v.Field
	return spayload
}

// NewSummaryPayload builds the payload of the "summary" endpoint of the
// "chatter" service from the gRPC request type.
func NewSummaryPayload(token string) *chattersvc.SummaryPayload {
	payload := &chattersvc.SummaryPayload{}
	payload.Token = token
	return payload
}

// NewChatSummaryCollection builds the gRPC response type from the result of
// the "summary" endpoint of the "chatter" service using the "default" view.
func NewChatSummaryCollection(vresult chattersvcviews.ChatSummaryCollectionView) *pb.ChatSummaryCollection {
	v := &pb.ChatSummaryCollection{}
	v.Field = make([]*pb.ChatSummary, len(vresult))
	for i, val := range vresult {
		v.Field[i] = &pb.ChatSummary{
			Message_: *val.Message,
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

func NewSummaryStreamingRequest(v *pb.SummaryStreamingRequest) string {
	spayload := v.Field
	return spayload
}

// NewHistoryPayload builds the payload of the "history" endpoint of the
// "chatter" service from the gRPC request type.
func NewHistoryPayload(view string, token string) *chattersvc.HistoryPayload {
	payload := &chattersvc.HistoryPayload{}
	payload.View = &view
	payload.Token = token
	return payload
}

// NewHistoryResponseTiny builds the gRPC response type from the result of the
// "history" endpoint of the "chatter" service using the "tiny" view.
func NewHistoryResponseTiny(vresult *chattersvcviews.ChatSummaryView) *pb.HistoryResponse {
	v := &pb.HistoryResponse{
		Message_: *vresult.Message,
	}
	return v
}

// NewHistoryResponse builds the gRPC response type from the result of the
// "history" endpoint of the "chatter" service using the "default" view.
func NewHistoryResponse(vresult *chattersvcviews.ChatSummaryView) *pb.HistoryResponse {
	v := &pb.HistoryResponse{
		Message_: *vresult.Message,
	}
	if vresult.Length != nil {
		v.Length = int32(*vresult.Length)
	}
	if vresult.SentAt != nil {
		v.SentAt = *vresult.SentAt
	}
	return v
}

// ValidateChatSummaryCollection runs the validations defined on
// ChatSummaryCollection.
func ValidateChatSummaryCollection(message *pb.ChatSummaryCollection) (err error) {
	for _, e := range message.Field {
		if e != nil {
			if err2 := ValidateChatSummary(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateChatSummary runs the validations defined on ChatSummary.
func ValidateChatSummary(message *pb.ChatSummary) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.sent_at", message.SentAt, goa.FormatDateTime))

	return
}

// ValidateHistoryResponse runs the validations defined on HistoryResponse.
func ValidateHistoryResponse(message *pb.HistoryResponse) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.sent_at", message.SentAt, goa.FormatDateTime))

	return
}
