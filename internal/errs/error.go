package errs

import "errors"

var (
	ErrInvalidUser = errors.New("[kuryr-admin] username or password is invalid")
	ErrUnknownUser = errors.New("[kuryr-admin] unknown user")

	ErrUnauthorized = errors.New("[kuryr-admin] unauthorized")
	ErrLoginExpired = errors.New("[kuryr-admin] login expired")
)
