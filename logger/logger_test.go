package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggerLogLevel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		actual   string
		errorMsg string
	}{
		{
			logLevel: "info",
			actual:   "info",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "warn",
			actual:   "warn",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "debug",
			actual:   "debug",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "error",
			actual:   "error",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "fatal",
			actual:   "fatal",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "trace",
			actual:   "trace",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "panic",
			actual:   "panic",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "plop",
			actual:   "info",
			errorMsg: "Fail to get expected logLevel",
		},
	}

	for i, tc := range tests {
		os.Setenv("CYPRESS_PARALLEL_API_LOG_LEVEL", tc.logLevel)
		SetLoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(z, tc.actual, fmt.Sprintf("%s on test number %d", tc.errorMsg, i))
		os.Unsetenv("CYPRESS_PARALLEL_API_LOG_LEVEL")
	}
}
