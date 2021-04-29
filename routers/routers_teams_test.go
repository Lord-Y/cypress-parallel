package routers

import (
	"fmt"
	"testing"

	"github.com/icrowley/fake"
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

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/teams", "")
	assert.Contains(w.Body.String(), "name")
}
