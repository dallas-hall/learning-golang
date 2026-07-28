package service_test

import (
	"service"
	"testing"
)

func Test_NewReturnsValidService(t *testing.T) {

}

// We followed the Red, Green, Refactor principle here. We started by writing
// a test with the code that we wanted to call. That helped us with the design
// of the "service" package. We implemented the "service" package functions
// just enough so the tests failed. The next steps would then be to make them
// work. And finally to refactor to make the code better. We stopped at Red.
func Test_RunningIsTrueWhenServiceIsRunning(t *testing.T) {
	t.Parallel()
	service.Start()
	if !service.Running() {
		t.Error(false)
	}
}
