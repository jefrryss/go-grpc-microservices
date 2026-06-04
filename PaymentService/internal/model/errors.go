package model

import "errors"

var (
	ErrInvalidPaymentMethod = errors.New("invalid or unspecified payment method")
)
