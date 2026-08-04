package user_test

import (
	"errors"
	"testing"
	"user"
)

func Test_FindUserGivesErrUserNotFoundForMissingUserUsingTypeAssertion(t *testing.T) {
	t.Parallel()
	_, err := user.FindUser("bogus user")
	// Before the errors package was implemented, type assertion was used to
	// inspect errors and their types. An interface in Go has 2 components: the
	// concrete type and concrete value. This is asking, does the concrete type
	// stored in this interface value match user.ErrUserNotFound. If yes, ok is
	// true and the first value is err re-typed as user.ErrUserNotFound. The
	// concrete value of the error interface is discarded.
	//
	// This is a comma-ok type assertion inside an if-with-init-statement.
	// comma-ok is `_, ok := ...` and all if statements are:
	// if OpitionalInitStatement; Condition {
	// ...
	// }
	// OptionalInitStatement runs first, then Condition is evaluated, then the
	// block executes if Condition is true. That's all the semicolon form is
	// doing - running one statement, then checking a boolean.
	if _, ok := err.(user.ErrUserNotFound); ok {
		t.Errorf("wrong error: %v", err)
	}

	// The above statement is the equivalent of
	_, ok := err.(user.ErrUserNotFound)
	if ok {
		t.Errorf("wrong error: %v", err)
	}
}

func Test_FindUserGivesErrUserNotFoundForMissingUserUsingErrorsAs(t *testing.T) {
	t.Parallel()
	_, err := user.FindUser("bogus user")
	// The previous type assertion can be replaced with errors.As.
	// errors.Is tells you if err is a specified error value.
	// errors.As tells you if err is a specified error type.
	if !errors.As(err, &user.ErrUserNotFound{}) {
		t.Errorf("wrong error: %v", err)
	}
}
