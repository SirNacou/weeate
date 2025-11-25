package domain

import "errors"

var (
	ErrOrderAlreadyExists             = errors.New("order already exists")
	ErrCannotCreateOrderWithZeroVotes = errors.New("cannot create order with zero votes")
)
