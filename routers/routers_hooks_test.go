package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestHooksPlainCreate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestProjectsCreate(t)
	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	tests := []struct {
		branch     string
		browser    string
		statusCode int
	}{
		{
			branch:     "",
			browser:    "chrome",
			statusCode: 201,
		},
		{
			branch:     "master",
			browser:    "chrome",
			statusCode: 201,
		},
		{
			branch:     "test",
			statusCode: 400,
		},
	}

	for _, tc := range tests {
		payload := fmt.Sprintf("project_name=%s", result["project_name"])
		if tc.branch != "" {
			payload += fmt.Sprintf("&branch=%s", tc.branch)
			payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
		}
		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/hooks/launch/plain", payload)
		assert.Equal(tc.statusCode, w.Code)
	}

	payload := fmt.Sprintf("project_name=%s", result["project_name"])
	payload += fmt.Sprintf("&branch=master")
	payload += fmt.Sprintf("&specs=bad_%s", tools.RandomValueFromSlice(specs))
	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/hooks/launch/plain", payload)
	assert.Equal(400, w.Code)
}
