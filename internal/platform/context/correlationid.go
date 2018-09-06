package context

import "context"

// CorrelationID returns a  correlation ID from a context (if any).
func CorrelationID(ctx context.Context) string {
	if cid, ok := ctx.Value(contextKeyCorrelationId).(string); ok {
		return cid
	}

	return ""
}

// WithCorrelationID returns a new context with the current correlation ID.
func WithCorrelationID(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, contextKeyCorrelationId, cid)
}
