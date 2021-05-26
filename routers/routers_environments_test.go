package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/environments"
	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentsCreate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestProjectsCreate(t)
	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	tests := []struct {
		key        string
		value      string
		statusCode int
	}{
		{
			key:        "key",
			value:      "value",
			statusCode: 201,
		},
		{
			key:        "key1",
			value:      "value1",
			statusCode: 201,
		},
		{
			key:        "",
			value:      "",
			statusCode: 400,
		},
	}

	for _, tc := range tests {
		payload := fmt.Sprintf("projectId=%s", result["project_id"])
		payload += fmt.Sprintf("&key=%s", tc.key)
		payload += fmt.Sprintf("&value=%s", tc.value)
		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/environments", payload)
		assert.Equal(tc.statusCode, w.Code)
	}
}

func TestEnvironmentsUpdate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	payload := fmt.Sprintf("projectId=%s", result["project_id"])
	payload += fmt.Sprintf("&environmentId=%s", result["environment_id"])
	payload += fmt.Sprintf("&key=%s", fake.CharactersN(5))
	payload += fmt.Sprintf("&value=%s", fake.CharactersN(5))
	w, _ := performRequest(router, headers, "PUT", "/api/v1/cypress-parallel-api/environments", payload)
	assert.Equal(200, w.Code)
}

func TestEnvironmentsUpdate_fail(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	payload := fmt.Sprintf("projectId=%s", result["project_id"])
	payload += fmt.Sprintf("&environmentId=%s", result["environment_id"])
	w, _ := performRequest(router, headers, "PUT", "/api/v1/cypress-parallel-api/environments", payload)
	assert.Equal(400, w.Code)
}

func TestEnvironmentsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/environments/list", "")
	assert.Contains(w.Body.String(), "environment_id")
}

func TestEnvironmentsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/environments/%s", result["environment_id"]), "")
	assert.Contains(w.Body.String(), "key")
}

func TestEnvironmentsDelete(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/cypress-parallel-api/environments/%s", result["environment_id"]), "")
	assert.Equal(200, w.Code)
}

func TestEnvironmentsSearch(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := environments.GetEnvironmentIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve teams")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/environments/search?q=%s", result["key"]), "")
	if len(result) > 0 {
		assert.Contains(w.Body.String(), "key")
		return
	}
	w, _ = performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/environments/search?q=", "")
	assert.Equal(400, w.Code)
}
