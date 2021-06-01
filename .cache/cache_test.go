// Package cache assemble all functions for redis
package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisSet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri    string
		key    string
		value  string
		expire time.Duration
		fail   bool
	}{
		{
			uri:    "redis://@127.0.0.1:6379",
			key:    "key",
			value:  "value",
			expire: 10,
			fail:   false,
		},
		{
			uri:    "redis://@127.0.0.1:63799",
			key:    "key",
			value:  "value",
			expire: 10,
			fail:   true,
		},
		{
			uri:    "redisss://@127.0.0.1:63799",
			key:    "key",
			value:  "value",
			expire: 10,
			fail:   true,
		},
	}

	for _, tc := range tests {
		z := RedisSet(tc.uri, tc.key, tc.value, tc.expire)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisGet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		key  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key2",
			fail: true,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			key:  "key",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			key:  "key",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key", "value", 2)
	for _, tc := range tests {
		_, z := RedisGet(tc.uri, tc.key)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisDel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		key  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key2",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			key:  "key",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			key:  "key",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key", "value", 60)
	for _, tc := range tests {
		z := RedisDel(tc.uri, tc.key)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisDelWithPrefix(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		key  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key1_",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:6379",
			key:  "key1_",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			key:  "key1_",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			key:  "key1_",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key1_xxx", "value", 30)
	for _, tc := range tests {
		z := RedisDelWithPrefix(tc.uri, tc.key)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisDeleteKeysHasPrefix(t *testing.T) {
	assert := assert.New(t)

	keys := []string{
		"key1_",
		"key2_",
	}
	tests := []struct {
		uri  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:6379",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key1_", "value", 60)
	for _, tc := range tests {
		z := RedisDeleteKeysHasPrefix(tc.uri, keys)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisFlushDB(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key", "value", 60)
	for _, tc := range tests {
		z := RedisFlushDB(tc.uri)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisFlushAll(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			fail: true,
		},
	}

	_ = RedisSet("redis://@127.0.0.1:6379", "key", "value", 60)
	for _, tc := range tests {
		z := RedisFlushAll(tc.uri)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}

func TestRedisPing(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		uri  string
		fail bool
	}{
		{
			uri:  "redis://@127.0.0.1:6379",
			fail: false,
		},
		{
			uri:  "redis://@127.0.0.1:63799",
			fail: true,
		},
		{
			uri:  "redisss://@127.0.0.1:63799",
			fail: true,
		},
	}

	for _, tc := range tests {
		_, z := RedisPing(tc.uri)
		if tc.fail {
			assert.Error(z)
		} else {
			assert.NoError(z)
		}
	}
}
