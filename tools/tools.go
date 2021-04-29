// Package tools assemble useful functions used by other packages
package tools

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"time"
)

// RandString function is explicit
func RandString(n int) string {
	if n < 1 {
		n = 1
	}
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// RandStringInt function is explicit
func RandStringInt(n int) string {
	if n < 1 {
		n = 1
	}
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// CaptureOutput permit to catch os stdout/stderr
func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

// GetPagination permit to set start limit and end limit for pagination
func GetPagination(page int, start int, end int, rangeLimit int) (startLimit int, EndLimit int) {
	if page <= 0 {
		page = 1
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 1
	}
	if page == 1 {
		start = 0
	}
	if rangeLimit < 1 {
		rangeLimit = 1
	}
	switch page {
	case 1:
		return start, end
	default:
		return rangeLimit * (page - 1), rangeLimit
	}
}
