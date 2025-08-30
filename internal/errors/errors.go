package errors

import "errors"

var (
	ErrAuthentication = errors.New("authentication error")
	ErrRegister       = errors.New("register error")
	ErrUserExists     = errors.New("user already exists")
)
