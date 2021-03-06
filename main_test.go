package main

import (
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	pg_uri = "postgres://cypress_parallel:cypress_parallel@127.0.0.1:5432/cypress_parallel?sslmode=disable"
)

// Test example found here: https://github.com/golang/go/issues/21000
func TestMain(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("CYPRESS_PARALLEL_DB_URI", pg_uri)
	defer os.Unsetenv("CYPRESS_PARALLEL_DB_URI")

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		main()
		<-sigc
		signal.Stop(sigc)
	}()

	err = proc.Signal(os.Interrupt)
	assert.NoError(err)
	time.Sleep(1 * time.Second)
}

func TestMain_set_port(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("CYPRESS_PARALLEL_DB_URI", pg_uri)
	defer os.Unsetenv("CYPRESS_PARALLEL_DB_URI")
	os.Setenv("CYPRESS_PARALLEL_PORT", "10000")
	defer os.Unsetenv("CYPRESS_PARALLEL_PORT")
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		main()
		<-sigc
		signal.Stop(sigc)
	}()

	err = proc.Signal(os.Interrupt)
	assert.NoError(err)
	time.Sleep(1 * time.Second)

}
