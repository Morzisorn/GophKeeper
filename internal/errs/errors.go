package errs

import "errors"

var (
	//Data validation
	ErrRequiredArgumentIsMissing = errors.New("one of required arguments is missing")

	//User errors
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyRegistered = errors.New("user is already registered")
	ErrIncorrectCredentials  = errors.New("incorrect login or password")

	//Item errors
	ErrIncorrectItemType = errors.New("incorrect item type")
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrItemNotFound = errors.New("item not found")

	//Other errors
	ErrInternalServerError = errors.New("internal server error")
)
