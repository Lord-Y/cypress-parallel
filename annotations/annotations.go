// Package annotations will manage all annotations requirements that will be injected in pods
package annotations

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// annotation struct handle requirements to create annotations
type annotation struct {
	ProjectID int    `form:"projectId" json:"projectId" binding:"required"`
	Key       string `form:"key" json:"key" binding:"required"`
	Value     string `form:"value" json:"value" binding:"required"`
}

// updateAnnotation struct handle requirements to create annotations
type updateAnnotation struct {
	ProjectID    int    `form:"projectId" json:"projectId" binding:"required"`
	AnnotationID int    `form:"annotationId" json:"annotationId" binding:"required"`
	Key          string `form:"key" json:"key" binding:"required"`
	Value        string `form:"value" json:"value" binding:"required"`
}

// getAnnotations struct handle requirements to get annotations
type getAnnotations struct {
	AnnotationID int `form:"annotationId" json:"annotationId" binding:"required"`
}

// deleteAnnotations struct handle requirements to create annotations
type deleteAnnotation struct {
	AnnotationID int `form:"annotationId" json:"annotationId" binding:"required"`
}

// listAnnotations struct handle requirements to get projects
type listAnnotations struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// searchAnnotations struct handle requirements to get projects
type searchAnnotations struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// listAnnotationsByProjectID struct handle requirements to get all annotations from project id
type listAnnotationsByProjectID struct {
	ProjectID int `form:"projectId" json:"projectId" binding:"required"`
}

// dbList struct permit to map data from db
type dbList struct {
	Annotation_id int       `json:"annotation_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Project_id    int       `json:"project_id"`
	Date          time.Time `json:"date"`
	Total         int       `json:"total"`
	Project_name  string    `json:"project_name"`
}

// dbRead struct permit to map data from db
type dbRead struct {
	Annotation_id int       `json:"annotation_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Project_id    int       `json:"project_id"`
	Date          time.Time `json:"date"`
	Project_name  string    `json:"project_name"`
}

// dbCommon struct permit to map data from db
type dbCommon struct {
	Annotation_id int       `json:"annotation_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Project_id    int       `json:"project_id"`
	Date          time.Time `json:"date"`
}

// Create handle requirements to create annotations with annotation struct
func Create(c *gin.Context) {
	var (
		p annotation
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.create()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"projectId": result})
	}
}

// Update handle requirements to update projects with updateAnnotation struct
func Update(c *gin.Context) {
	var (
		p updateAnnotation
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.update()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusOK, "OK")
	}
}

// List handle requirements to read projects with getProjects struct
func List(c *gin.Context) {
	var (
		p listAnnotations
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.StartLimit, p.EndLimit = tools.GetPagination(p.Page, 0, commons.GetRangeLimit(), commons.GetRangeLimit())

	result, err := p.list()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatus(204)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Read handle requirements to read annotations with getAnnotations struct
func Read(c *gin.Context) {
	var (
		p getAnnotations
	)
	id := c.Params.ByName("annotationId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "annotationId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.AnnotationID = vID

	result, err := p.read()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if result == (dbRead{}) {
		c.AbortWithStatus(204)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// Delete handle deletion of annotation viariable with deleteAnnotation struct
func Delete(c *gin.Context) {
	var (
		p deleteAnnotation
	)
	id := c.Params.ByName("annotationId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "annotationId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.AnnotationID = vID

	err = p.delete()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

// Search handle requirements to search projects with searchProjects struct
func Search(c *gin.Context) {
	var (
		p searchAnnotations
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.StartLimit, p.EndLimit = tools.GetPagination(p.Page, 0, commons.GetRangeLimit(), commons.GetRangeLimit())

	result, err := p.search()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatus(204)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// ListByProjectID handle requirements to list annotations by project id
func ListByProjectID(c *gin.Context) {
	var (
		p listAnnotationsByProjectID
	)
	id := c.Params.ByName("projectId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.ProjectID = vID
	result, err := p.listByProjectID()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatus(404)
	} else {
		c.JSON(http.StatusOK, result)
	}
}
