package internal

import "errors"

const (
	errNoDuplicateRecord = 1062
)

var (
	ErrUserAlreadyExist = errors.New("user already exists")
)
