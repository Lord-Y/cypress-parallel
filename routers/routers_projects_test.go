package routers

import (
	"fmt"
	"testing"

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
	payload += "&repository=https://github.com/cypress-io/cypress.git"
	payload += "&branch=master"

	log.Info().Msgf("payload %s", payload)

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
