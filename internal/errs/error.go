package errs

import "errors"

var (
	ErrUnauthorized = errors.New("[kuryr-admin] unauthorized")

	ErrRecordNotFound = errors.New("[kuryr-admin] record not found")

	ErrInvalidAccountType = errors.New("[kuryr-admin] invalid account type")
	ErrInvalidVerifyType  = errors.New("[kuryr-admin] invalid verify type")

	ErrInvalidUser = errors.New("[kuryr-admin] account or login credential is invalid")
	ErrUnknownUser = errors.New("[kuryr-admin] unknown user")

	ErrDuplicateKey   = errors.New("[kuryr-admin] duplicate key violation")
	ErrBizKeyConflict = errors.New("[kuryr-admin] biz key already exists")
)
