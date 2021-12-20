package routers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Lord-Y/cypress-parallel/projects"
	"github.com/Lord-Y/cypress-parallel/teams"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestProjectsCreate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestTeamsCreate(t)
	result, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
	payload += fmt.Sprintf("&teamId=%s", result["team_id"])
	payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
	payload += "&branch=master"
	payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
	payload += "&scheduling=* */10 * * *"
	payload += "&schedulingEnabled=false"
	payload += "&maxPods=10"
	payload += fmt.Sprintf("&cypress_docker_version=%s", tools.RandomValueFromSlice(cypressVersions))
	payload += "&browser=chrome"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "POST", "/api/v1/projects", payload)
	assert.Equal(201, w.Code)
}

func TestProjectsCreate_fail(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestTeamsCreate(t)
	result, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
	payload += fmt.Sprintf("&teamId=%s", result["team_id"])
	payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
	payload += "&branch=master"
	payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
	payload += "&scheduling=* */10 * * *"
	payload += "&schedulingEnabled=false"
	payload += "&maxPods=10"
	payload += fmt.Sprintf("&cypress_docker_version=%s", tools.RandomValueFromSlice(cypressVersions))

	router := SetupRouter()
	w, _ := performRequest(router, headers, "POST", "/api/v1/projects", payload)
	assert.Equal(201, w.Code)
}

func TestProjectsRead(t *testing.T) {
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
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/projects/%s", result["project_id"]), "")
	assert.Contains(w.Body.String(), "name")
}

func TestProjectsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/projects/list", "")
	assert.Contains(w.Body.String(), "name")
}

func TestProjectsCreateMulti(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	router := SetupRouter()
	for i := 1; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
		payload += fmt.Sprintf("&teamId=%s", result["team_id"])
		payload += "&repository=https://github.com/cypress-io/cypress-example-kitchensink.git"
		payload += "&branch=master"
		payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
		payload += "&scheduling="
		payload += "&schedulingEnabled=false"
		payload += fmt.Sprintf("&maxPods=%d", fake.MonthNum())
		payload += "&browser=chrome"

		w, _ := performRequest(router, headers, "POST", "/api/v1/projects", payload)
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
	payload += fmt.Sprintf("&specs=%s", tools.RandomValueFromSlice(specs))
	payload += "&scheduling="
	payload += "&schedulingEnabled=false"
	payload += fmt.Sprintf("&maxPods=%d", fake.MonthNum())
	payload += fmt.Sprintf("&browser=%s", result["browser"])

	router := SetupRouter()
	w, _ := performRequest(router, headers, "PUT", "/api/v1/projects", payload)
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
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/projects/%s", result["project_id"]), "")
	assert.Equal(200, w.Code)
}

func TestProjectsSearch(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := projects.GetProjectIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve teams")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/projects/search?q=%s", result["project_name"]), "")
	if len(result) > 0 {
		assert.Contains(w.Body.String(), "project_name")
		return
	}
	w, _ = performRequest(router, headers, "GET", "/api/v1/projects/search?q=", "")
	assert.Equal(400, w.Code)
}

func TestProjectsAll(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestTeamsCreate(t)
	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/projects/all", "")
	assert.Contains(w.Body.String(), "name")
}
