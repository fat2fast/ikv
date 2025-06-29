package usermodel

import "errors"

var (
	ErrNameIsRequired          = errors.New("user name is required")
	ErrNameIsEmpty             = errors.New("user name is empty")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserStatusInvalid       = errors.New("status must be in (active, inactive, pending, deleted)")
	ErrUserIsDeleted           = errors.New("user is deleted")
	ErrUserBannedOrDeleted     = errors.New("user banned or deleted")
	ErrInvalidEmailAndPassword = errors.New("invalid email or password")
)
