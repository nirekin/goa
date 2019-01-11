package middleware

import (
	"google.golang.org/grpc/metadata"
)

// MetadataHasKey returns true if the given key is found in the metadata
// or false if it doesn't exist.
func MetadataHasKey(md metadata.MD, key string) bool {
	vals := md.Get(key)
	return len(vals) > 0
}

// MetadataValue returns the first value for the given metadata key if
// key exists, else returns an empty string.
func MetadataValue(md metadata.MD, key string) string {
	if MetadataHasKey(md, key) {
		return md.Get(key)[0]
	}
	return ""
}

// WithSpan returns a metadata containing the given trace, span and parent span
// IDs.
func WithSpan(md metadata.MD, traceID, spanID, parentID string) metadata.MD {
	if parentID != "" {
		md.Set(ParentSpanIDMetadataKey, parentID)
	}
	md.Set(TraceIDMetadataKey, traceID)
	md.Set(SpanIDMetadataKey, spanID)
	return md
}
