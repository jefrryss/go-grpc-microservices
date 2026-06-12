package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrPartNotFound  = errors.New("one or more requested parts do not exist in inventory")

	ErrOrderCannotBeCancelled = errors.New("order is already paid or cannot be cancelled")
	ErrOrderAlreadyPaid       = errors.New("order is already paid")
	ErrInvalidOrderStatus     = errors.New("invalid order status for this operation")

	ErrEmptyPartList        = errors.New("order must contain at least one part")
	ErrInvalidPaymentMethod = errors.New("invalid or unspecified payment method")

	ErrNilProtoPart   = errors.New("cannot convert nil proto part")
	ErrNilDomainOrder = errors.New("cannot convert nil domain order")
	ErrNilRepoOrder   = errors.New("cannot convert nil repo order")
)
