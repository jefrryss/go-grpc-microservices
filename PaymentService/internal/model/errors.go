package model

import "errors"

var (
	ErrInvalidPaymentMethod = errors.New("invalid or unspecified payment method")
	ErrEmptyOrderUUID       = errors.New("order UUID cannot be empty")
	ErrEmptyUserUUID        = errors.New("user UUID cannot be empty")
)
