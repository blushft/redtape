package redtape

import "context"

type RequestContext struct {
	Metadata map[string]interface{}
	Context  context.Context
}

func (rc RequestContext) getKey(k string) interface{} {
	if v, ok := rc.Metadata[k]; ok {
		return v
	}

	if v := rc.Context.Value(k); v != nil {
		return v
	}

	return nil
}

type Request struct {
	Resource string         `json:"resource"`
	Action   string         `json:"action"`
	Role     string         `json:"subject"`
	Context  RequestContext `json:"-"`
}

func NewRequest(ctx context.Context, res, action, role string, meta map[string]interface{}) *Request {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Request{
		Resource: res,
		Action:   action,
		Role:     role,
		Context: RequestContext{
			Context:  ctx,
			Metadata: meta,
		},
	}
}
