package redtape

import "context"

type RequestContextKey struct{}

type RequestOption func(*Request)

func RequestResource(res string) RequestOption {
	return func(r *Request) {
		r.Resource = res
	}
}

func RequestAction(action string) RequestOption {
	return func(r *Request) {
		r.Action = action
	}
}

func RequestSubject(sub Subject) RequestOption {
	return func(r *Request) {
		r.Subject = sub
	}
}

func RequestScope(s string) RequestOption {
	return func(r *Request) {
		r.Scope = s
	}
}

func RequestContext(ctx context.Context, meta ...map[string]interface{}) RequestOption {
	return func(r *Request) {
		r.Context = NewRequestContext(ctx, meta...)
	}
}

func WithMetadata(meta ...map[string]interface{}) RequestOption {
	return func(r *Request) {
		if r.Context == nil {
			r.Context = NewRequestContext(context.Background(), meta...)
			return
		}

		r.AddMetadata(meta...)
	}
}

// Request represents a request to be matched against a policy set.
type Request struct {
	Resource string          `json:"resource"`
	Action   string          `json:"action"`
	Subject  Subject         `json:"subject"`
	Scope    string          `json:"scope"`
	Context  context.Context `json:"-"`
}

// NewRequestWithContext builds a request from the provided options.
func NewRequest(opts ...RequestOption) *Request {
	req := &Request{
		Context: NewRequestContext(context.Background()),
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// Metadata returns metadata stored in context or an empty set.
func (r *Request) Metadata() RequestMetadata {
	return RequestMetadataFromContext(r.Context)
}

func (r *Request) AddMetadata(meta ...map[string]interface{}) {
	rmd := RequestMetadataFromContext(r.Context)
	for _, md := range meta {
		for k, v := range md {
			rmd[k] = v
		}
	}

	r.Context = context.WithValue(r.Context, RequestMetadataKey{}, rmd)
}

// RequestMetadata is a helper type to allow type safe retrieval.
type RequestMetadata map[string]interface{}

// RequestMetadataKey is a type to identify RequestMetadata embedded in context.
type RequestMetadataKey struct{}

// NewRequestContext builds a context object from an existing context, embedding request metadata. If nil
// values are provided to both arguments, new values are created or returned.
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

// RequestMetadataFromContext extracts RequestMetadata from a given context or returns an empty metadata set.
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
