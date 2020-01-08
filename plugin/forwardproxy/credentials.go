package forwardproxy

import "errors"

var ErrInvalidCredentials = errors.New("Invalid credentials")

type CredentialsRepository interface {
	GetPassword(login string) (string, error)
}
