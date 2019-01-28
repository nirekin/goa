package xray

import (
	"context"
	"fmt"
	"net"
	"time"

	grpcm "goa.design/goa/grpc/middleware"
	"goa.design/goa/middleware"
	"goa.design/goa/middleware/xray"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewUnaryServer returns a server middleware that sends AWS X-Ray segments
// to the daemon running at the given address. It stores the request segment
// in the context. User code can further configure the segment for example to
// set a service version or record an error. It extracts the trace information
// from the incoming unary request metadata using the tracing middleware
// package. The tracing middleware must be mounted on the service.
//
// service is the name of the service reported to X-Ray. daemon is the hostname
// (including port) of the X-Ray daemon collecting the segments.
//
// User code may create child segments using the Segment NewSubsegment method
// for tracing requests to external services. Such segments should be closed via
// the Close method once the request completes. The middleware takes care of
// closing the top level segment. Typical usage:
//
//     if s := ctx.Value(SegKey); s != nil {
//       segment := s.(*xray.Segment)
//     }
//     sub := segment.NewSubsegment("external-service")
//     defer sub.Close()
//     err := client.MakeRequest()
//     if err != nil {
//         sub.Error = xray.Wrap(err)
//     }
//     return
//
func NewUnaryServer(service, daemon string) (grpc.UnaryServerInterceptor, error) {
	connection, err := xray.Connect(context.Background(), time.Minute, func() (net.Conn, error) {
		return net.Dial("udp", daemon)
	})
	if err != nil {
		return nil, fmt.Errorf("xray: failed to connect to daemon - %s", err)
	}
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var s *GRPCSegment
		ctx, s = withSegment(ctx, service, info.FullMethod, connection)
		if s == nil {
			return handler(ctx, req)
		}
		defer s.Close()
		return handler(ctx, req)
	}), nil
}

// NewStreamServer is similar to NewUnaryServer except it is used for
// streaming endpoints.
func NewStreamServer(service, daemon string) (grpc.StreamServerInterceptor, error) {
	connection, err := xray.Connect(context.Background(), time.Minute, func() (net.Conn, error) {
		return net.Dial("udp", daemon)
	})
	if err != nil {
		return nil, fmt.Errorf("xray: failed to connect to daemon - %s", err)
	}
	return grpc.StreamServerInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		var s *GRPCSegment
		ctx, s = withSegment(ctx, service, info.FullMethod, connection)
		wss := grpcm.NewWrappedServerStream(ctx)
		if s == nil {
			return handler(srv, wss)
		}
		defer s.Close()
		return handler(srv, wss)
	}), nil
}

// UnaryClient middleware creates XRay subsegments if a segment is found in
// the context and stores the subsegment to the context. It also sets the
// trace information in the context which is used by the tracing middleware.
// This middleware must be mounted before the tracing middleware.
func UnaryClient(host string) grpc.UnaryClientInterceptor {
	return grpc.UnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		seg := ctx.Value(xray.SegKey)
		if seg == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		s := seg.(*xray.Segment)
		sub := GRPCSegment{s.NewSubsegment(host)}
		defer sub.Close()

		// update the context with the latest segment
		ctx = middleware.WithSpan(ctx, sub.TraceID, sub.ID, sub.ParentID)
		sub.RecordRequest(ctx, method, "remote")
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			sub.RecordError(err)
		} else {
			sub.RecordResponse()
		}
		return nil
	})
}

// StreamClient is the streaming endpoint middleware equivalent for UnaryClient.
func StreamClient(host string) grpc.StreamClientInterceptor {
	return grpc.StreamClientInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		seg := ctx.Value(xray.SegKey)
		if seg == nil {
			return streamer(ctx, desc, cc, method, opts...)
		}
		s := seg.(*xray.Segment)
		sub := GRPCSegment{s.NewSubsegment(host)}
		defer sub.Close()

		// update the context with the latest segment
		ctx = middleware.WithSpan(ctx, sub.TraceID, sub.ID, sub.ParentID)
		sub.RecordRequest(ctx, method, "remote")
		cs, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			sub.RecordError(err)
		} else {
			sub.RecordResponse()
		}
		return cs, nil
	})
}

// withSegment creates a new X-Ray segment and stores it in the context.
// It also returns the newly created segment.
func withSegment(ctx context.Context, service, method string, connection func() net.Conn) (context.Context, *GRPCSegment) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// incoming metadata does not exist. Probably trace middleware is not
		// loaded before this one.
		return ctx, nil
	}

	var (
		traceID  string
		spanID   string
		parentID string
	)
	{
		spanID = ctx.Value(middleware.TraceSpanIDKey)
		traceID = ctx.Value(middleware.TraceIDKey)
		parentID = ctx.Value(middleware.TraceParentSpanIDKey)
	}
	if traceID == nil || spanID == nil {
		return ctx, nil
	}
	s := &GRPCSegment{xray.NewSegment(service, traceID.(string), spanID.(string), connection())}
	s.RecordRequest(ctx, method, "")
	if parentID != "" {
		s.ParentID = parentID.(string)
	}
	return context.WithValue(ctx, xray.SegKey, s.Segment), s
}
