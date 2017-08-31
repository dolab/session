package session

import "errors"

var (
	ErrNotFound       = errors.New("Session not found.")
	ErrDuplicateID    = errors.New("Session ID is duplicated.")
	ErrEmptySecretKey = errors.New("Session secret key is required.")
	ErrDuplicateKey   = errors.New("Session key has existed.")
	ErrCookieExpired  = errors.New("Session has expired.")
)
