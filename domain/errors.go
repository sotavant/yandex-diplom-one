package domain

import "errors"

var (
	ErrBadParams            = errors.New("params is not valid")
	ErrInternalServerError  = errors.New("internal Server Error")
	ErrLoginExist           = errors.New("login is busy")
	ErrPasswordTooWeak      = errors.New("password too weak")
	ErrBadUserData          = errors.New("wrong login/password")
	ErrBadOrderNum          = errors.New("bad order number")
	ErrOrderAlreadyUploaded = errors.New("order uploaded by another user")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserNotAuthorized    = errors.New("user not authorized")
)
