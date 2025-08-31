package errors

import "errors"

var (
	ErrAuthentication = errors.New("authentication error")
	ErrRegister       = errors.New("register error")
	ErrUserExists     = errors.New("user already exists")
	ErrTokenExpired   = errors.New("token expired")
	ErrUpdateAccount  = errors.New("update account error")
	ErrUniqueEmail    = errors.New("unique email error")
)
