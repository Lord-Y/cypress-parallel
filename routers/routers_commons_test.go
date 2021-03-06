package routers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/health", "")

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealth_no_xrequestid(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/health", "")

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealth_prometheus(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	os.Setenv("CYPRESS_PARALLEL_PROMETHEUS", "1")
	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/health", "")

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
	os.Unsetenv("CYPRESS_PARALLEL_PROMETHEUS")
}
