package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggerLogLevel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		expected string
		errorMsg string
	}{
		{
			logLevel: "info",
			expected: "info",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "warn",
			expected: "warn",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "debug",
			expected: "debug",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "error",
			expected: "error",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "fatal",
			expected: "fatal",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "trace",
			expected: "trace",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "panic",
			expected: "panic",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "plop",
			expected: "info",
			errorMsg: "Fail to get expected logLevel",
		},
	}

	for _, tc := range tests {
		os.Setenv("CYPRESS_PARALLEL_API_LOG_LEVEL", tc.logLevel)
		SetLoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(tc.expected, z)
		os.Unsetenv("CYPRESS_PARALLEL_API_LOG_LEVEL")
	}
}
