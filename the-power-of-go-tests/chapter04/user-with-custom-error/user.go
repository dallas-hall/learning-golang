package user

import (
	"fmt"
)

type User struct {
	Name  string
	Email string
}

var userDB map[string]*User

// Prior to Go 1.13 and errors.Is, you had to create a struct type with the
// name being the error and the fields having the extra sentinel values.
type ErrUserNotFound struct {
	User string
}

// Since error is an interface, so any type can implement it which we do here.
func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("user %q not found", e.User)
}

func FindUser(name string) (*User, error) {
	user, ok := userDB[name]
	if !ok {
		// Create a wrapped sentinel error prior to errors.Is existed.
		return nil, fmt.Errorf("%q: %w", name, ErrUserNotFound{User: name})
	}
	return user, nil
}
