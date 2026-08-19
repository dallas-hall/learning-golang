package server

import (
	"math/rand"
	"net"
	"time"
)

func ListenAsync(address string) {
	// "go" starts this function running concurrently in a goroutine. ListenAsync
	// returns immediately here. It does NOT wait for the code below to run. The
	// server starts "in the background".
	go func() {
		// Add a random delay to simulate a flakey connection.
		time.Sleep(time.Duration(rand.Intn(175)) * time.Millisecond)

		// Start listening on address. Ignoring the error here (the underscore) is
		// only OK because this is demo/test code.
		l, _ := net.Listen("tcp", address)
		defer l.Close()

		// Block forever so the listener stays open. An empty select with no cases
		// waits on nothing, forever. This keeps the goroutine (and the listener)
		// alive for demo purposes.
		select {}
	}()
}
