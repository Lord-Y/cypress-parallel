// Package commons assemble all functions used in other packages
package commons

import (
	"os"
	"strings"
)

// BuildDSN permit to retrieve string url to connect to the sql instance
func BuildDSN() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_DB_URI"))
}

// GetRedisURI permit to retrieve OS env variable
func GetRedisURI() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_REDIS_URI"))
}
