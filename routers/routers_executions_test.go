package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel/executions"
	"github.com/Lord-Y/cypress-parallel/projects"
	"github.com/Lord-Y/cypress-parallel/tools"
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
	payload := fmt.Sprintf("project_name=%s", result.Project_name)
	payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
	w, _ := performRequest(router, headers, "POST", "/api/v1/hooks/launch/plain", payload)
	assert.Equal(201, w.Code)

	w, _ = performRequest(router, headers, "GET", "/api/v1/executions/list", "")
	assert.Equal(200, w.Code)
}

func TestExecutionsUpdateResult(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve executions")
		t.Fail()
		return
	}
	payload := `result={"key": "key", "value": "value", "environment_id": 35}`
	payload += "&executionStatus=DONE"
	payload += fmt.Sprintf("&branch=%s", resultEx.Branch)
	payload += fmt.Sprintf("&spec=%s", resultEx.Spec)
	payload += fmt.Sprintf("&uniqId=%s", resultEx.Uniq_id)

	w, _ := performRequest(router, headers, "POST", "/api/v1/executions/update", payload)
	if resultEx == (executions.DBRead{}) {
		assert.Equal(400, w.Code)
		return
	}
	assert.Equal(200, w.Code)

	// rollback
	payload = `result={}`
	payload += "&executionStatus=NOT_STARTED"
	payload += fmt.Sprintf("&branch=%s", resultEx.Branch)
	payload += fmt.Sprintf("&spec=%s", resultEx.Spec)
	payload += fmt.Sprintf("&uniqId=%s", resultEx.Uniq_id)

	w, _ = performRequest(router, headers, "POST", "/api/v1/executions/update", payload)
	assert.Equal(200, w.Code)
}

func TestExecutionsUpdateResult_fail(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve executions")
		t.Fail()
		return
	}
	payload := `result={"key": "key", "value": "value", "environment_id": 35}`
	payload += "&executionStatus=DONE"
	payload += fmt.Sprintf("&branch=%s", resultEx.Branch)
	payload += fmt.Sprintf("&spec=%s", resultEx.Spec)

	w, _ := performRequest(router, headers, "POST", "/api/v1/executions/update", payload)
	assert.Equal(400, w.Code)
}

func TestExecutionsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	resultEx, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		assert.Fail("Fail to retrieve executions")
		return
	}

	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/executions/%d", resultEx.Execution_id), "")
	if resultEx == (executions.DBRead{}) {
		assert.Equal(404, w.Code)
		return
	}
	assert.Equal(200, w.Code)
}

func TestExecutionsSearch(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		assert.Fail("Fail to retrieve executions")
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/executions/search?q=%s", result.Branch), "")
	if result != (executions.DBRead{}) {
		assert.Contains(w.Body.String(), "branch")
		return
	}
	w, _ = performRequest(router, headers, "GET", "/api/v1/executions/search?q=", "")
	assert.Equal(400, w.Code)
}

func TestExecutionsUniqID(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestHooksPlainCreate(t)

	result, err := executions.GetExecutionIDForUnitTesting()
	if err != nil {
		assert.Fail("Fail to retrieve executions")
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/executions/list/by/uniqid/%s", result.Uniq_id), "")
	if result != (executions.DBRead{}) {
		assert.Contains(w.Body.String(), "branch")
		return
	}
	w, _ = performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/executions/list/by/uniqid/%s", "404"), "")
	assert.Equal(404, w.Code)
}
