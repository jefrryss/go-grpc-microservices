package model

import "errors"

var (
	ErrNotFound        = errors.New("item not found")
	ErrInvalidCateogry = errors.New("invalid or unspecified part category")
	ErrNilModelPart    = errors.New("cannot convert nil part to repo part")
	ErrNilRepoPArt     = errors.New("cannot convert nil part to model part")

	ErrInvalidFilter   = errors.New("invalid filter parameters")
	ErrInternalSystem  = errors.New("internal system error")
	ErrInvalidCategory = errors.New("invalid part category")

	ErrNilProtoPart  = errors.New("cannot convert nil proto part")
	ErrNilDomainPart = errors.New("cannot convert nil domain part")
)
