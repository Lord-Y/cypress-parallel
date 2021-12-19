// Package commons assemble all functions used in other packages
package commons

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// BuildDSN permit to retrieve string url to connect to the sql instance
func BuildDSN() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_DB_URI"))
}

// GetRedisURI permit to retrieve OS env variable
func GetRedisURI() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_REDIS_URI"))
}

// GetRangeLimit return de max range limit for postgresSQL queries
func GetRangeLimit() int {
	return 25
}

// GetKubernetesMode permit to retrieve OS env variable
func GetKubernetesMode() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_K8S_CLIENT_OUTSIDE"))
}

// GetKubernetesKubeConfig permit to retrieve OS env variable
func GetKubernetesKubeConfig() string {
	return strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_K8S_KUBE_CONFIG"))
}

// GetKubernetesJobsNamespace permit to retrieve OS env variable
func GetKubernetesJobsNamespace() string {
	z := strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_JOBS_NAMESPACE"))
	if z == "" {
		return "cypress-parallel-jobs"
	} else {
		return z
	}
}

// GetMaxSpecs permit to retrieve OS env variable
func GetMaxSpecs() int {
	harcoded := 3
	max := strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_MAX_SPECS"))
	if max == "" {
		return harcoded
	} else {
		m, err := strconv.Atoi(max)
		if err != nil {
			log.Error().Err(err).Msgf("Error occured while converting string to int so let's set it to %d anyway", harcoded)
			return harcoded
		}
		return m
	}
}

// GetAPIUrl permit to retrieve OS env variable
func GetAPIUrl() string {
	z := strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_URL"))
	if z == "" {
		return "http://127.0.0.1:8080"
	} else {
		return z
	}
}
