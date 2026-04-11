package errs

import "errors"

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrInvalidCredentials  = errors.New("invalid credential")
	ErrTokenNotFound       = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
	ErrTweetNotFound       = errors.New("tweet not found")
	ErrTweetMediaNotFound  = errors.New("tweet media not found")
	ErrUnauthorizedUpdate  = errors.New("user is not authorized to update this tweet")

	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidmediaType = errors.New("invalid media type")

	ErrCacheKeyNotFound = errors.New("key not found")
)
