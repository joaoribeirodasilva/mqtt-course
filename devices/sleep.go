package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SleepChannel make a thread sleep for a certain number of
// milliseconds or if an interrupt signal (SIGINT or SIGTERM)
// is received
// returns true if the sleep timeout is achieved and
// false if a interrupt signal is received
func SleepChannel(sleep_time time.Duration) bool {

	sleepChannel := time.After(sleep_time)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sleepChannel:
		return true
	case <-signalChannel:
		return false
	}
}
