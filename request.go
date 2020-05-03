package redtape

import "context"

// Request represents a request to be matched against a policy set
type Request struct {
	Resource string          `json:"resource"`
	Action   string          `json:"action"`
	Role     string          `json:"subject"`
	Scope    string          `json:"scope"`
	Context  context.Context `json:"-"`
}

func NewRequest(res, action, role, scope string, meta ...map[string]interface{}) *Request {
	return &Request{
		Resource: res,
		Action:   action,
		Role:     role,
		Scope:    scope,
		Context:  NewRequestContext(nil, meta...),
	}
}

// NewRequestWithContext builds a request from the provided parameters
func NewRequestWithContext(ctx context.Context, res, action, role, scope string, meta ...map[string]interface{}) *Request {
	return &Request{
		Resource: res,
		Action:   action,
		Role:     role,
		Scope:    scope,
		Context:  NewRequestContext(ctx, meta...),
	}
}

// Metadata returns metadata stored in context or an empty set
func (r *Request) Metadata() RequestMetadata {
	return RequestMetadataFromContext(r.Context)
}

// RequestMetadata is a helper type to allow type safe retrieval
type RequestMetadata map[string]interface{}

// RequestMetadataKey is a type to identify RequestMetadata embedded in context
type RequestMetadataKey struct{}

// NewRequestContext builds a context object from an existing context, embedding request metadata. If nil
// values are provided to both arguments, new values are created or returned
func NewRequestContext(ctx context.Context, meta ...map[string]interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	reqmeta := RequestMetadata{}

	for _, md := range meta {
		for k, v := range md {
			reqmeta[k] = v
		}
	}

	ctx = context.WithValue(ctx, RequestMetadataKey{}, reqmeta)

	return ctx
}

// RequestMetadataFromContext extracts RequestMetadata from a given context or returns an empty metadata set
func RequestMetadataFromContext(ctx context.Context) RequestMetadata {
	if ctx == nil {
		return RequestMetadata{}
	}

	md := ctx.Value(RequestMetadataKey{})

	if md == nil {
		return RequestMetadata{}
	}

	return md.(RequestMetadata)
}
