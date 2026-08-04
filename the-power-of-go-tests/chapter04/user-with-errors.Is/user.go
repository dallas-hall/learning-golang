package user

import (
	"errors"
	"fmt"
)

type User struct {
	Name  string
	Email string
}

var userDB map[string]*User

var ErrUserNotFound = errors.New("user not found")

func FindUser(name string) (*User, error) {
	user, ok := userDB[name]
	if !ok {
		// Create a wrapped sentinel error, it may be unwrapped by the caller with
		// errors.Is
		return nil, fmt.Errorf("%q: %w", name, ErrUserNotFound)
	}
	return user, nil
}
