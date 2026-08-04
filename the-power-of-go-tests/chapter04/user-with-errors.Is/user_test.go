package user_test

import (
	"errors"
	"testing"
	"user"
)

func Test_FindUserGivesErrUserNotFoundForMissingUser(t *testing.T) {
	t.Parallel()
	_, err := user.FindUser("bogus user")
	// We use errors.Is to unwrap the original error wraped inside the sentinel
	// error, user.ErrUserNotFound
	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("wrong error: %v", err)
	}
}
