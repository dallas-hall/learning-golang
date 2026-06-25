package main

import "fmt"

func main() {
	/*
		By default channels are unbuffered, meaning that they will only accept sends (chan <-) if there is a corresponding receive (<- chan) ready to receive the sent value. Buffered channels accept a limited number of values without a corresponding receiver for those values.

		Here we make a channel of strings buffering up to 2 values.
	*/
	msgs := make(chan string, 2)

	// Because this channel is buffered, we can send these values into the channel without a corresponding concurrent receive.
	msgs <- "buffered"
	msgs <- "channel"

	// Later we can receive these two values as usual.
	fmt.Println(<-msgs)
	fmt.Println(<-msgs)
}
