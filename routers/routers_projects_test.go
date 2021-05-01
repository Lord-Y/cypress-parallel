package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/Lord-Y/cypress-parallel-api/teams"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestProjectsCreate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestTeamsCreate(t)
	resultTeamID, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
	payload += fmt.Sprintf("&teamId=%s", resultTeamID["team_id"])
	payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
	payload += "&branch=master"
	payload += "&specs=cypress/integration/examples/actions.spec.js"
	payload += "&scheduling=* */10 * * *"
	payload += "&schedulingEnabled=false"
	payload += "&maxPods=10"
	payload += "&cypress_docker_version=7.1.0"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/projects", payload)
	assert.Equal(201, w.Code)
}

func TestProjectsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/projects", "")
	assert.Contains(w.Body.String(), "name")
}

func TestProjectsCreateMulti(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resultTeamID, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	router := SetupRouter()
	for i := 1; i < 5; i++ {
		payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
		payload += fmt.Sprintf("&teamId=%s", resultTeamID["team_id"])
		payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
		payload += "&branch=master"
		payload += "&specs=cypress/integration/examples/actions.spec.js"
		payload += "&scheduling="
		payload += "&schedulingEnabled=false"
		payload += fmt.Sprintf("&maxPods=%d", fake.MonthNum())

		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/projects", payload)
		assert.Equal(201, w.Code)
	}
}

func TestProjectsUpdate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve project and team id")
		t.Fail()
		return
	}
	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
	payload += fmt.Sprintf("&projectId=%s", result["project_id"])
	payload += fmt.Sprintf("&teamId=%s", result["team_id"])
	payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
	payload += "&branch=master"
	payload += "&specs=cypress/integration/examples/actions.spec.js"
	payload += "&scheduling="
	payload += "&schedulingEnabled=false"
	payload += fmt.Sprintf("&maxPods=%d", fake.MonthNum())

	router := SetupRouter()
	w, _ := performRequest(router, headers, "PUT", "/api/v1/cypress-parallel-api/projects", payload)
	assert.Equal(200, w.Code)
}

func TestProjectsDelete(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve project and team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/cypress-parallel-api/projects/%s", result["project_id"]), "")
	assert.Equal(200, w.Code)
}
