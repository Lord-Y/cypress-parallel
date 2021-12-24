// Package environments will manage all environments requirements that will be injected in pods
package environments

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// environment struct handle requirements to create environments
type environment struct {
	ProjectID int    `form:"projectId" json:"projectId" binding:"required"`
	Key       string `form:"key" json:"key" binding:"required"`
	Value     string `form:"value" json:"value" binding:"required"`
}

// updateEnvironment struct handle requirements to create annotations
type updateEnvironment struct {
	ProjectID     int    `form:"projectId" json:"projectId" binding:"required"`
	EnvironmentID int    `form:"environmentId" json:"environmentId" binding:"required"`
	Key           string `form:"key" json:"key" binding:"required"`
	Value         string `form:"value" json:"value" binding:"required"`
}

// getEnvironments struct handle requirements to get environments
type getEnvironments struct {
	EnvironmentID int `form:"environmentId" json:"environmentId" binding:"required"`
}

// deleteEnvironment struct handle requirements to delete environment
type deleteEnvironment struct {
	EnvironmentID int `form:"environmentId" json:"environmentId" binding:"required"`
}

// listEnvironments struct handle requirements to get projects
type listEnvironments struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// searchEnvironments struct handle requirements to get projects
type searchEnvironments struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// listEnvironmentsByProjectID struct handle requirements to get all environment from project id
type listEnvironmentsByProjectID struct {
	ProjectID int `form:"projectId" json:"projectId" binding:"required"`
}

// dbList struct permit to map data from db
type dbList struct {
	Environment_id int       `json:"environment_id"`
	Key            string    `json:"key"`
	Value          string    `json:"value"`
	Project_id     int       `json:"project_id"`
	Date           time.Time `json:"date"`
	Total          int       `json:"total"`
	Project_name   string    `json:"project_name"`
}

// dbRead struct permit to map data from db
type dbRead struct {
	Environment_id int       `json:"environment_id"`
	Key            string    `json:"key"`
	Value          string    `json:"value"`
	Project_id     int       `json:"project_id"`
	Date           time.Time `json:"date"`
	Project_name   string    `json:"project_name"`
}

// dbCommon struct permit to map data from db
type dbCommon struct {
	Environment_id int       `json:"environment_id"`
	Key            string    `json:"key"`
	Value          string    `json:"value"`
	Project_id     int       `json:"project_id"`
	Date           time.Time `json:"date"`
}

// Create handle requirements to create environments with environments struct
func Create(c *gin.Context) {
	var (
		p environment
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

// Update handle requirements to update environment with updateEnvironment struct
func Update(c *gin.Context) {
	var (
		p updateEnvironment
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

// List handle requirements to list environments
func List(c *gin.Context) {
	var (
		p listEnvironments
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

// Read handle requirements to read environments with getEnvironments struct
func Read(c *gin.Context) {
	var (
		p getEnvironments
	)
	id := c.Params.ByName("environmentId")
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

	p.EnvironmentID = vID

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

// Delete handle deletion of environment viariable with deleteEnvironment struct
func Delete(c *gin.Context) {
	var (
		p deleteEnvironment
	)
	id := c.Params.ByName("environmentId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "environmentId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.EnvironmentID = vID

	err = p.delete()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

// Search handle requirements to search environments by key or value
func Search(c *gin.Context) {
	var (
		p searchEnvironments
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

// ListByProjectID handle requirements to list environments by project id
func ListByProjectID(c *gin.Context) {
	var (
		p listEnvironmentsByProjectID
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
