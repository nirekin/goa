package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	// RequestIDOption uses a constructor pattern to customize middleware
	RequestIDOption func(*requestIDOption) *requestIDOption

	// requestIDOption is the struct storing all the options.
	requestIDOption struct {
		// useXRequestIDMetadata is true to use incoming "X-Request-Id" metadata,
		// instead of always generating unique IDs, when present in request.
		// Defaults to always-generate.
		useXRequestIDMetadata bool
		// xRequestMetadataLimit is positive to truncate incoming "X-Request-Id"
		// metadata at the specified length. Defaults to no limit.
		xRequestMetadataLimit int
	}
)

const (
	// RequestIDMetadataKey is the key containing the request ID in the gRPC
	// metadata.
	RequestIDMetadataKey = "X-Request-Id"
)

// RequestID returns a middleware, which initializes the metadata with a unique
// value under the RequestIDMetadata key. Optionally uses the incoming
// "X-Request-Id" metadata key, if present, with or without a length limit to
// use as request ID. The default behavior is to always generate a new ID.
//
// examples of use:
//  service.Use(middleware.RequestID())
//
//  // enable options for using "X-Request-Id" metadata key with length limit.
//  service.Use(middleware.RequestID(
//    middleware.UseXRequestIDMetadataOption(true),
//    middleware.XRequestMetadataLimitOption(128)))
func RequestID(options ...RequestIDOption) grpc.UnaryServerInterceptor {
	o := new(requestIDOption)
	for _, option := range options {
		o = option(o)
	}
	return grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		var id string
		{
			if o.useXRequestIDMetadata {
				id = MetadataValue(md, RequestIDMetadataKey)
				if o.xRequestMetadataLimit > 0 && len(id) > o.xRequestMetadataLimit {
					id = id[:o.xRequestMetadataLimit]
				}
			} else {
				id = shortID()
			}
		}
		md.Set(RequestIDMetadataKey, id)
		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	})
}

// UseXRequestIDMetadataOption enables/disables using "X-Request-Id" metadata.
func UseXRequestIDMetadataOption(f bool) RequestIDOption {
	return func(o *requestIDOption) *requestIDOption {
		o.useXRequestIDMetadata = f
		return o
	}
}

// XRequestMetadataLimitOption sets the option for limiting "X-Request-Id"
// metadata length.
func XRequestMetadataLimitOption(limit int) RequestIDOption {
	return func(o *requestIDOption) *requestIDOption {
		o.xRequestMetadataLimit = limit
		return o
	}
}
