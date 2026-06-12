package model

import "errors"

var (
	ErrUncorrectUUIDOrder   = errors.New("Uncorrect UUID order")
	ErrUncorrectUUIDUser    = errors.New("Unccoert UUID user")
	ErrInvalidPaymentMethod = errors.New("invalid or unspecified payment method")
)
