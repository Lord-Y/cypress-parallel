package tools

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value  int
		actual int
	}{
		{
			value:  0,
			actual: 1,
		},
		{
			value:  -20,
			actual: 1,
		},
		{
			value:  1,
			actual: 1,
		},
		{
			value:  4,
			actual: 4,
		},
		{
			value:  500,
			actual: 500,
		},
	}

	for _, tc := range tests {
		z := RandString(tc.value)
		assert.Equal(len(z), tc.actual)
	}
}

func TestRandStringInt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value  int
		actual int
	}{
		{
			value:  0,
			actual: 1,
		},
		{
			value:  -20,
			actual: 1,
		},
		{
			value:  1,
			actual: 1,
		},
		{
			value:  4,
			actual: 4,
		},
		{
			value:  500,
			actual: 500,
		},
	}

	for _, tc := range tests {
		z := RandStringInt(tc.value)
		assert.Equal(len(z), tc.actual)
	}
}

func TestCaptureOutput(t *testing.T) {
	assert := assert.New(t)

	f1 := func() {
		log.Print("f1")
	}

	output := CaptureOutput(f1)
	assert.Contains(output, "f1")
}

func TestGetPagination(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		page             int
		start            int
		end              int
		rangeLimit       int
		actualStartLimit int
		actualEndLimit   int
	}{
		{
			page:             0,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             1,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             2,
			start:            2,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 20,
			actualEndLimit:   20,
		},
		{
			page:             -1,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             -1,
			start:            -1,
			end:              -20,
			rangeLimit:       -20,
			actualStartLimit: 0,
			actualEndLimit:   1,
		},
	}

	for _, tc := range tests {
		sl, el := GetPagination(tc.page, tc.start, tc.end, tc.rangeLimit)
		assert.Equal(sl, tc.actualStartLimit)
		assert.Equal(el, tc.actualEndLimit)
	}
}
