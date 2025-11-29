package http

import (
	"log/slog"
	"net/http"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/danielgtaylor/huma/v2"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	// Default values
	code := domain.EInternal
	msg := "Internal Server Error"
	op := domain.Op("unknown")
	var innerErr error

	// 1. Unpack the error
	if e, ok := err.(*domain.Error); ok {
		code = e.Code
		msg = e.Message
		op = e.Op
		innerErr = e.Err
	} else {
		// It's a raw system error (panic, etc)
		innerErr = err
	}

	// 2. Determine HTTP Status
	status := httpStatusFromCode(code)

	// 3. STRUCTURED LOGGING
	// We build a list of attributes. We ALWAYS want Op and Code.
	attrs := []any{
		slog.String("op", string(op)),
		slog.String("code", string(code)),
	}

	// Only log the inner error if it exists
	if innerErr != nil {
		attrs = append(attrs, slog.String("inner_error", innerErr.Error()))
	}

	// 4. Log Level Decision
	// 5xx = ERROR (Alerting)
	// 4xx = WARN (Monitoring)
	if status >= 500 {
		// Sanitize message for 500s to User
		msg = "Internal Server Error"
		slog.Error("Domain Error", attrs...)
	} else {
		slog.Warn("Client Error", attrs...)
	}

	// 5. Return Safe Error to Huma
	return huma.NewError(status, msg)
}

// Helper map function
func httpStatusFromCode(code domain.Code) int {
	switch code {
	case domain.EInvalid:
		return http.StatusBadRequest
	case domain.ENotFound:
		return http.StatusNotFound
	case domain.EConflict:
		return http.StatusConflict
	case domain.EForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
