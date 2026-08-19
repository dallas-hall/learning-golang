package server_test

import (
	"net"
	"server"
	"testing"
	"time"
)

func randomLocalAddr(t *testing.T) string {
	t.Helper()
	// Get a random free port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	return listener.Addr().String()
}

// waitForServer polls address until a TCP connection succeeds or the timeout
// elapses, failing the test via t.Fatal if the server never comes up in time.
func waitForServer(t *testing.T, address string) {
	t.Helper()
	// NewTimer starts a timer and returns a Timer struct. Its field C is a
	// channel: when the timer expires, the current time gets sent on it.
	timeout := time.NewTimer(100 * time.Millisecond)
	defer timeout.Stop()

	// Try connecting once before entering the retry loop.
	_, err := net.Dial("tcp", address)
	for err != nil {
		// select lets us wait on multiple channels. With a default case, it
		// doesn't block: it just checks whether timeout.C already has a value
		// waiting, and falls through to default if not.
		select {
		case <-timeout.C:
			// Something arrived on timeout.C, meaning the timer fired. The <- reads
			// (receives) the value off the channel; we don't care about the value
			// itself, just that it happened.
			t.Fatal("timed out")
		default:
			// Nothing on timeout.C yet, so the timer hasn't fired. Wait briefly and
			// try connecting again.
			t.Log("retrying...")
			time.Sleep(time.Millisecond)
			_, err = net.Dial("tcp", address)
		}
	}
	t.Log("connected")
}

func TestListenAsyncListensOnGivenAddress(t *testing.T) {
	t.Parallel()
	address := randomLocalAddr(t)
	server.ListenAsync(address)
	waitForServer(t, address)
	// Test whatever when the server is ready.
}
