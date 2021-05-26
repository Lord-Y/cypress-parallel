package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel-api/annotations"
	"github.com/Lord-Y/cypress-parallel-api/projects"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationsCreateOrUpdate(t *testing.T) {
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
		annotations string
		statusCode  int
	}{
		{
			annotations: `[{"key": "key", "value": "value"},{"key": "key1", "value": "value1"}]`,
			statusCode:  201,
		},
		{
			annotations: `[{"key": "key", "value": "value1"},{"key": "key1", "value": "value1"}]`,
			statusCode:  201,
		},
		{
			annotations: `[{"key": "key_new1", "value": "value_new1"},{"key": "key_new2", "value": "value_new2"}]`,
			statusCode:  201,
		},
		{
			annotations: "",
			statusCode:  400,
		},
	}

	for _, tc := range tests {
		payload := fmt.Sprintf("projectId=%s", result["project_id"])
		if tc.annotations != "" {
			payload += fmt.Sprintf("&annotations=%s", tc.annotations)
		}
		w, _ := performRequest(router, headers, "POST", "/api/v1/cypress-parallel-api/annotations", payload)
		assert.Equal(tc.statusCode, w.Code)
	}
}

func TestAnnotationsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/annotations/list", "")
	assert.Contains(w.Body.String(), "annotation_id")
}

func TestAnnotationsRead(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := annotations.GetAnnotationIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/annotations/%s", result["annotation_id"]), "")
	assert.Contains(w.Body.String(), "key")
}

func TestAnnotationsDelete(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := annotations.GetAnnotationIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve team id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/cypress-parallel-api/annotations/%s", result["annotation_id"]), "")
	assert.Equal(200, w.Code)
}

func TestAnnotationsSearch(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := annotations.GetAnnotationIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve teams")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/cypress-parallel-api/annotations/search?q=%s", result["key"]), "")
	if len(result) > 0 {
		assert.Contains(w.Body.String(), "key")
		return
	}
	w, _ = performRequest(router, headers, "GET", "/api/v1/cypress-parallel-api/annotations/search?q=", "")
	assert.Equal(400, w.Code)
}
