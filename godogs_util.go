package main

import "time"

func someMsLater(ms int) func() error {
	return func() error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}
}
