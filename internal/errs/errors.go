package errs

import "errors"

var (
	//Data validation
	ErrRequiredArgumentIsMissing = errors.New("one of required arguments is missing")

	//User errors
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyRegistered = errors.New("user is already registered")
	ErrIncorrectCredentials  = errors.New("incorrect login or password")
	
	//Other errors
	ErrInternalServerError   = errors.New("internal server error")
)