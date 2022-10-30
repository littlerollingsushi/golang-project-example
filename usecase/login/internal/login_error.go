package internal

import "errors"

var (
	ErrEmptyEmail        = errors.New("login email can not be empty")
	ErrEmptyPassword     = errors.New("login password can not be empty")
	ErrInvalidPassword   = errors.New("login password is not valid")
	ErrInvalidPrivateKey = errors.New("login usecase private key is not valid")
	ErrUserNotFound      = errors.New("user with given email is not found")
)
