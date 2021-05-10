package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/executions"
	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestExecutionsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestProjectsCreate(t)
	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve project and team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	payload := fmt.Sprintf("project_name=%s", result["project_name"])
	payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/hooks/launch/plain", payload)
	assert.Equal(201, w.Code)

	w, _ = performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/executions/list", "")
	assert.Equal(200, w.Code)
}

func TestExecutionsUpdateResult(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve project and team id")
		t.Fail()
		return
	}
	payload := `result={"key": "key", "value": "value", "environment_id": 35}`
	payload += "&executionStatus=DONE"
	payload += fmt.Sprintf("&branch=%s", resultEx["branch"])
	payload += fmt.Sprintf("&spec=%s", resultEx["spec"])
	payload += fmt.Sprintf("&uniqId=%s", resultEx["spec"])

	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/executions/update", payload)
	assert.Equal(200, w.Code)

	// rollback
	payload = `result={}`
	payload += "&executionStatus=NOT_STARTED"
	payload += fmt.Sprintf("&branch=%s", resultEx["branch"])
	payload += fmt.Sprintf("&spec=%s", resultEx["spec"])
	payload += fmt.Sprintf("&uniqId=%s", resultEx["spec"])

	w, _ = performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/executions/update", payload)
	assert.Equal(200, w.Code)
}

func TestExecutionsUpdateResult_fail(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve project and team id")
		t.Fail()
		return
	}
	payload := `result={"key": "key", "value": "value", "environment_id": 35}`
	payload += "&executionStatus=DONE"
	payload += fmt.Sprintf("&branch=%s", resultEx["branch"])
	payload += fmt.Sprintf("&spec=%s", resultEx["spec"])

	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/executions/update", payload)
	assert.Equal(400, w.Code)
}

func TestExecutionsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		assert.Fail("Fail to retrieve project and team id")
		return
	}

	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/executions?executionId=%s", resultEx["execution_id"]), "")
	assert.Equal(200, w.Code)
}
