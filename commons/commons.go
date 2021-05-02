// Package commons assemble all functions used in other packages
package commons

import (
	"os"
	"strings"
)

// BuildDSN permit to retrieve string url to connect to the sql instance
func BuildDSN() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_API_DB_URI"))
}

// GetRedisURI permit to retrieve OS env variable
func GetRedisURI() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_API_REDIS_URI"))
}

// GetRangeLimit return de max range limit for postgresSQL queries
func GetRangeLimit() int {
	return 25
}
