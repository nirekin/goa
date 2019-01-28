package xray

import (
	"context"
	"net"

	"goa.design/goa/grpc/middleware"
	"goa.design/goa/middleware/xray"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// GRPCSegment represents an AWS X-Ray segment document for gRPC services.
type GRPCSegment struct {
	*xray.Segment
}

// RecordRequest traces a request.
//
// It sets Http.Request & Namespace (ex: "remote")
func (s *GRPCSegment) RecordRequest(ctx context.Context, method, namespace string) {
	s.Lock()
	defer s.Unlock()

	if s.HTTP == nil {
		s.HTTP = &xray.HTTP{}
	}

	s.Namespace = namespace
	s.HTTP.Request = requestData(ctx, method)
}

// RecordResponse traces a response.
//
// It sets Throttle, Fault, Error and HTTP.Response
func (s *GRPCSegment) RecordResponse() {
	s.Lock()
	defer s.Unlock()

	if s.HTTP == nil {
		s.HTTP = &xray.HTTP{}
	}

	s.HTTP.Response = &xray.Response{
		Status: int(codes.OK),
	}
}

// requestData creates a Request from a http.Request.
func requestData(ctx context.Context, method string) *xray.Request {
	var agent string
	{
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			agent = middleware.MetadataValue(md, "user-agent")
		}
	}
	var ip string
	{
		if p, ok := peer.FromContext(ctx); ok {
			ip, _, _ = net.SplitHostPort(p.Addr.String())
		}
	}

	return &xray.Request{
		Method:    method,
		UserAgent: agent,
		ClientIP:  ip,
	}
}
