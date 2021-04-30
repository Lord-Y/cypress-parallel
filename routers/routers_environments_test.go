package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/environments"
	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentsCreateOrUpdate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestProjectsCreate(t)
	resultTeamID, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	tests := []struct {
		environments string
		statusCode   int
	}{
		{
			environments: `[{"key": "key", "value": "value"},{"key": "key1", "value": "value1"}]`,
			statusCode:   201,
		},
		{
			environments: `[{"key": "key", "value": "value1"},{"key": "key1", "value": "value1"}]`,
			statusCode:   201,
		},
		{
			environments: `[{"key": "key_new1", "value": "value_new1"},{"key": "key_new2", "value": "value_new2"}]`,
			statusCode:   201,
		},
		{
			environments: "",
			statusCode:   400,
		},
	}

	for _, tc := range tests {
		payload := fmt.Sprintf("projectId=%s", resultTeamID["project_id"])
		if tc.environments != "" {
			payload += fmt.Sprintf("&environments=%s", tc.environments)
		}
		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/environments", payload)
		assert.Equal(tc.statusCode, w.Code)
	}
}

func TestEnvironmentsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/environments", "")
	assert.Contains(w.Body.String(), "key")
}

func TestEnvironmentsDelete(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resultEnvironmentID, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/cypress-parallel-api/environments/%s", resultEnvironmentID["environment_id"]), "")
	assert.Equal(200, w.Code)
}
