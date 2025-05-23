package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {

	if os.Getenv("RUN_MAIN_TEST") == "1" {
		go main()
		time.Sleep(10 * time.Millisecond)
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("main() panicked: %v", r)
			}
			done <- true
		}()
		main()
	}()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
	}
}

func TestMainPackageExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main package test in short mode")
	}

}
