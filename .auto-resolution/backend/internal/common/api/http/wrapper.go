package http

import (
	"context"
)

// HandlerFunc defines the standard signature of a Huma handler.
// I = Input struct (e.g. RegisterRequest)
// O = Output struct (e.g. RegisterResponse)
type HandlerFunc[I any, O any] func(context.Context, *I) (*O, error)

// Handle wraps a controller function to automatically handle errors.
// It acts as a middleware specifically for the return values.
func Handle[I any, O any](logic HandlerFunc[I, O]) func(context.Context, *I) (*O, error) {
	return func(ctx context.Context, input *I) (*O, error) {
		// 1. Execute the business logic
		out, err := logic(ctx, input)

		// 2. Intercept the error
		if err != nil {
			// 3. Map it to the correct HTTP/Huma response using our existing MapError
			return nil, MapError(err)
		}

		// 4. Return success
		return out, nil
	}
}