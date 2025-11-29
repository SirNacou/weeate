package domain

import (
	"errors"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
)

var (
	ErrOrderAlreadyExists             = errors.New("order already exists")
	ErrCannotCreateOrderWithZeroVotes = errors.New("cannot create order with zero votes")
	ErrPollIDRequired                 = domain.NewError(domain.EInvalid, "poll ID is required")
	ErrOrderResultsRequired           = domain.NewError(domain.EInvalid, "order results are required")
)
