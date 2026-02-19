package user_const

import (
	"github.com/Fi44er/cloud-store-api/pkg/customerr"
)

var (
	ErrEmailRequired      = customerr.NewError(400, "email is required")
	ErrEmailTooLong       = customerr.NewError(400, "email is too long (max 254 characters)")
	ErrEmailInvalidFormat = customerr.NewError(400, "invalid email format")
)
