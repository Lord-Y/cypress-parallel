package routers

import (
	"fmt"
	"testing"

	"github.com/Lord-Y/cypress-parallel/annotations"
	"github.com/Lord-Y/cypress-parallel/projects"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationsCreate(t *testing.T) {
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
		payload := fmt.Sprintf("projectId=%d", result.Project_id)
		payload += fmt.Sprintf("&key=%s", tc.key)
		payload += fmt.Sprintf("&value=%s", tc.value)
		w, _ := performRequest(router, headers, "POST", "/api/v1/annotations", payload)
		assert.Equal(tc.statusCode, w.Code)
	}
}

func TestAnnotationsUpdate(t *testing.T) {
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

	payload := fmt.Sprintf("projectId=%d", result.Project_id)
	payload += fmt.Sprintf("&annotationId=%d", result.Annotation_id)
	payload += fmt.Sprintf("&key=%s", fake.CharactersN(5))
	payload += fmt.Sprintf("&value=%s", fake.CharactersN(5))
	w, _ := performRequest(router, headers, "PUT", "/api/v1/annotations", payload)
	assert.Equal(200, w.Code)
}

func TestAnnotationsUpdate_fail(t *testing.T) {
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
	payload := fmt.Sprintf("projectId=%d", result.Project_id)
	payload += fmt.Sprintf("&annotationId=%d", result.Annotation_id)
	w, _ := performRequest(router, headers, "PUT", "/api/v1/annotations", payload)
	assert.Equal(400, w.Code)
}

func TestAnnotationsList(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", "/api/v1/annotations/list", "")
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
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/annotations/%d", result.Annotation_id), "")
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
	w, _ := performRequest(router, headers, "DELETE", fmt.Sprintf("/api/v1/annotations/%d", result.Annotation_id), "")
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
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/annotations/search?q=%s", result.Key), "")
	assert.Contains(w.Body.String(), "key")
	w, _ = performRequest(router, headers, "GET", "/api/v1/annotations/search?q=", "")
	assert.Equal(400, w.Code)
}

func TestAnnotationsByProjectID(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	result, err := annotations.GetAnnotationIDForUnitTesting()
	if err != nil {
		log.Err(err).Msgf("Fail to retrieve environment id")
		t.Fail()
		return
	}

	router := SetupRouter()
	w, _ := performRequest(router, headers, "GET", fmt.Sprintf("/api/v1/annotations/list/by/projectid/%d", result.Project_id), "")
	assert.Equal(200, w.Code)
}
