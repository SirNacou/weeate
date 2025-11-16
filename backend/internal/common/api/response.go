package api

type Response[T any] struct {
	Body T
}

func NewResponse[T any](body T) *Response[T] {
	return &Response[T]{Body: body}
}
