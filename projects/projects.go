// Package projects will manage all projects requirements
package projects

import (
	"net/http"
	"strconv"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// projects struct handle requirements to create projects
type projects struct {
	TeamID               int    `form:"teamId" json:"teamId" binding:"required"`
	Name                 string `form:"name" json:"name" binding:"required,max=100"`
	Repository           string `form:"repository" json:"repository" binding:"required"`
	Branch               string `form:"branch" json:"branch" binding:"required"`
	Specs                string `form:"specs" json:"specs" binding:"required"`
	Scheduling           string `form:"scheduling" json:"scheduling" binding:"max=15"`
	SchedulingEnabled    bool   `form:"schedulingEnabled" json:"schedulingEnabled"`
	MaxPods              int    `form:"maxPods,default=10" json:"maxPods"`
	CypressDockerVersion string `form:"cypress_docker_version,default=7.2.0-0.0.5" json:"cypress_docker_version"`
	Timeout              int    `form:"timeout,default=10" json:"timeout"`
	Username             string `form:"username" json:"username"`
	Password             string `form:"password" json:"password"`
	Browser              string `form:"browser,default=chrome" json:"browser" binding:"max=100,oneof=chrome firefox"`
	ConfigFile           string `form:"config_file,default=cypress.json" json:"config_file" binding:"max=100"`
}

// getProjects struct handle requirements to get projects
type getProjects struct {
	ProjectID int `form:"projectId" json:"projectId" binding:"required"`
}

// listProjects struct handle requirements to get projects
type listProjects struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// updateProjects struct handle requirements to update projects
type updateProjects struct {
	ProjectID            int    `form:"projectId" json:"projectId" binding:"required"`
	TeamID               int    `form:"teamId" json:"teamId" binding:"required"`
	Name                 string `form:"name" json:"name" binding:"required,max=100"`
	Repository           string `form:"repository" json:"repository" binding:"required"`
	Branch               string `form:"branch" json:"branch" binding:"required"`
	Specs                string `form:"specs" json:"specs" binding:"required"`
	Scheduling           string `form:"scheduling" json:"scheduling" binding:"max=15"`
	SchedulingEnabled    bool   `form:"schedulingEnabled" json:"schedulingEnabled"`
	MaxPods              int    `form:"maxPods,default=10" json:"maxPods"`
	CypressDockerVersion string `form:"cypress_docker_version,default=7.2.0-0.0.5" json:"cypress_docker_version"`
	Timeout              int    `form:"timeout,default=10" json:"timeout"`
	Username             string `form:"username" json:"username" binding:"max=100"`
	Password             string `form:"password" json:"password" binding:"max=100"`
	Browser              string `form:"browser,default=chrome" json:"browser" binding:"max=100,oneof=chrome firefox"`
	ConfigFile           string `form:"config_file,default=cypress.json" json:"config_file" binding:"max=100"`
}

// deleteProject struct handle requirements to delete project
type deleteProject struct {
	ProjectID int `form:"projectId" json:"projectId" binding:"required"`
}

// searchProjects struct handle requirements to get projects
type searchProjects struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// Create handle requirements to create projects with projects struct
func Create(c *gin.Context) {
	var (
		p projects
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// default value not supported yet with json
	if p.CypressDockerVersion == "" {
		p.CypressDockerVersion = "7.2.0-0.0.5"
	}
	if p.Timeout == 0 {
		p.Timeout = 10
	}
	if p.ConfigFile == "" {
		p.ConfigFile = "cypress.json"
	}
	if p.MaxPods == 0 {
		p.MaxPods = 10
	}

	result, err := p.create()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"projectId": result})
	}
}

// Read handle requirements to read projects with getProjects struct
func Read(c *gin.Context) {
	var (
		p getProjects
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

	result, err := p.read()
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

// List handle requirements to read projects with getProjects struct
func List(c *gin.Context) {
	var (
		p listProjects
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

// All handle requirements to return all projects
func All(c *gin.Context) {
	result, err := all()
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

// Update handle requirements to update projects with updateProjects struct
func Update(c *gin.Context) {
	var (
		p updateProjects
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// default value not supported yet with json
	if p.CypressDockerVersion == "" {
		p.CypressDockerVersion = "7.2.0-0.0.5"
	}
	if p.Timeout == 0 {
		p.Timeout = 10
	}
	if p.ConfigFile == "" {
		p.ConfigFile = "cypress.json"
	}
	if p.MaxPods == 0 {
		p.MaxPods = 10
	}

	err := p.update()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusOK, "OK")
	}
}

// Delete handle deletion of project deleteProject struct
func Delete(c *gin.Context) {
	var (
		p deleteProject
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
		p searchProjects
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
