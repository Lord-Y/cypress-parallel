package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/teams"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestTeamsCreate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))

	router := SetupRouter()
	w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/teams", payload)
	assert.Equal(201, w.Code)
}

func TestTeamsRead(t *testing.T) {
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
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/teams/%s", result["team_id"]), "")
	assert.Contains(w.Body.String(), "name")
}

func TestTeamsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/teams/list", "")
	assert.Contains(w.Body.String(), "name")
}

func TestTeamsCreateMulti(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	for i := 1; i < 5; i++ {
		payload := fmt.Sprintf("name=%s", fake.CharactersN(10))

		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/teams", payload)
		assert.Equal(201, w.Code)
	}
}

func TestTeamsUpdate(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}
	payload := fmt.Sprintf("name=%s", fake.CharactersN(10))
	payload += fmt.Sprintf("&teamId=%s", result["team_id"])

	router := SetupRouter()
	w, _ := performRequest(router, headers, "PUT", "/api/v1/cypress-parallel-api/teams", payload)
	assert.Equal(200, w.Code)
}

func TestTeamsDelete(t *testing.T) {
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
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/cypress-parallel-api/teams/%s", result["team_id"]), "")
	assert.Equal(200, w.Code)
}

func TestTeamsSearch(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := teams.GetTeamIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve teams")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/teams/search?q=%s", result["team_name"]), "")
	if len(result) > 0 {
		assert.Contains(w.Body.String(), "team_name")
		return
	}
	w, _ = performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/teams/search?q=", "")
	assert.Equal(400, w.Code)
}

func TestTeamsAll(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	TestTeamsCreate(t)
	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/teams/all", "")
	assert.Contains(w.Body.String(), "name")
}
