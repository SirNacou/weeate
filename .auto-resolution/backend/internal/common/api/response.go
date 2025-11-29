package api

type Response[T any] struct {
	Body T
}

func NewResponse[T any](body *T) *Response[T] {
	if body == nil {
		var empty T
		return &Response[T]{Body: empty}
	}
	return &Response[T]{Body: *body}
}
