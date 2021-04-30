// Package projects will manage all projects requirements
package projects

import (
	"net/http"
	"strconv"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// projects struct handle requirements to create projects
type projects struct {
	TeamID     int    `form:"teamId" json:"teamId" binding:"required"`
	Name       string `form:"name" json:"name" binding:"required,max=100"`
	Repository string `form:"repository" json:"repository" binding:"required"`
	Branch     string `form:"branch" json:"branch" binding:"required"`
}

// getProjects struct handle requirements to get projects
type getProjects struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// updateProjects struct handle requirements to update projects
type updateProjects struct {
	ProjectID         int    `form:"projectId" json:"projectId" binding:"required"`
	TeamID            int    `form:"teamId" json:"teamId" binding:"required"`
	Name              string `form:"name" json:"name" binding:"required,max=100"`
	Repository        string `form:"repository" json:"repository" binding:"required"`
	Branch            string `form:"branch" json:"branch" binding:"required"`
	Specs             string `form:"specs" json:"specs" binding:"required"`
	Scheduling        string `form:"scheduling" json:"scheduling" binding:"max=15"`
	SchedulingEnabled bool   `form:"schedulingEnabled" json:"schedulingEnabled"`
	MaxPods           int    `form:"maxPods,default=10" json:"maxPods"`
}

// deleteProject struct handle requirements to delete project
type deleteProject struct {
	ProjectID int `form:"projectId" json:"projectId" binding:"required"`
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
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.StartLimit, p.EndLimit = tools.GetPagination(p.Page, 0, commons.GetRangeLimit(), commons.GetRangeLimit())

	result, err := p.read()
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

// Update handle requirements to update projects with Updateprojects struct
func Update(c *gin.Context) {
	var (
		p updateProjects
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
	convID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.ProjectID = convID

	err = p.delete()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}
