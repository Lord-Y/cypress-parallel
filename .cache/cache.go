// Package cache assemble all functions for redis
package cache

import (
	"fmt"
	"regexp"
	"time"

	"context"

	customLogger "github.com/Lord-Y/cypress-parallel-api/logger"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// init func
func init() {
	customLogger.SetLoggerLogLevel()
}

// RedisSet permit to put k/v in redis
func RedisSet(redisURI string, keyName string, value string, expire time.Duration) (err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	err = rdb.SetNX(ctx, keyName, value, expire*time.Second).Err()
	if err != nil {
		return fmt.Errorf("Failed to set keyName on redis address %s, error message: %s", opt.Addr, err.Error())
	}
	return
}

// RedisGet permit to get k/v in redis
func RedisGet(redisURI string, keyName string) (z string, err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return "", fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	z, err = rdb.Get(ctx, keyName).Result()
	return
}

// RedisGet permit to get k/v in redis
func RedisDel(redisURI string, keyName string) (err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	err = rdb.Del(ctx, keyName).Err()
	if err != nil {
		return fmt.Errorf("Error occured deleting keyName %s on redis address %s, error message: %s", keyName, opt.Addr, err.Error())
	}
	return
}

// RedisGet permit to get k/v in redis
func RedisKeys(redisURI string, keyPrefix string) (z []string, err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return z, fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	z, err = rdb.Keys(ctx, keyPrefix).Result()
	if err != nil {
		return z, fmt.Errorf("Error occured while finding key with prefix %s address %s, error message: %s", keyPrefix, opt.Addr, err.Error())
	}
	return
}

// RedisDelWithPrefix permit to list keys keys has prefix like keyName* and then, delete keys
func RedisDelWithPrefix(redisURI string, keyPrefix string) (err error) {
	matched, _ := regexp.MatchString(`\*$`, keyPrefix)
	if !matched {
		log.Debug().Msgf("keyPrefix %s does not finish with wildcard, so let's set it", keyPrefix)
		keyPrefix = fmt.Sprintf("%s*", keyPrefix)
	}
	keys, err := RedisKeys(redisURI, keyPrefix)
	if err != nil {
		return fmt.Errorf("Error occured while getting redis keys %s", err)
	}
	if len(keys) > 0 {
		for _, key := range keys {
			RedisDel(redisURI, key)
		}
	}
	return
}

// RedisDeleteKeysHasPrefix permit to list keys keys has prefix like keyName* and then, delete keys
func RedisDeleteKeysHasPrefix(redisURI string, prefixes []string) (err error) {
	if len(prefixes) > 0 {
		for _, keyPrefix := range prefixes {
			matched, _ := regexp.MatchString(`\*$`, keyPrefix)
			if !matched {
				log.Debug().Msgf("keyPrefix %s does not finish with wildcard, so let's set it", keyPrefix)
				keyPrefix = fmt.Sprintf("%s*", keyPrefix)
			}
			keys, err := RedisKeys(redisURI, keyPrefix)
			if err != nil {
				return fmt.Errorf("Error occured while getting redis keys %s", err)
			}
			if len(keys) > 0 {
				for _, key := range keys {
					RedisDel(redisURI, key)
				}
			}
		}
	}
	return
}

// RedisFlushDB permit to flush actual used db in redis
func RedisFlushDB(redisURI string) (err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	err = rdb.FlushDB(ctx).Err()
	if err != nil {
		return fmt.Errorf("Error occured while while flushing DB on address %s, error message: %s", opt.Addr, err.Error())
	}
	return
}

// RedisFlushAll permit to flush all db in redis
func RedisFlushAll(redisURI string) (err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	err = rdb.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("Error occured while while flushing all DBs on address %s, error message: %s", opt.Addr, err.Error())
	}
	return
}

// RedisPing permit to get redis status
func RedisPing(redisURI string) (b bool, err error) {
	var (
		ctx = context.Background()
	)
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return false, fmt.Errorf("Error occured while parsing to redis address you provided, error message: %s", err.Error())
	}

	rdb := redis.NewClient(opt)
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return false, fmt.Errorf("Error occured while pinging redis address %s, error message: %s", opt.Addr, err.Error())
	}
	return true, nil
}
