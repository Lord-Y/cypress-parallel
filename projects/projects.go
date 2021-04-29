// Package projects will manage all projects requirements
package projects

import (
	"net/http"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Projects struct handle requirements to create projects
type Projects struct {
	TeamID     int    `form:"teamId" json:"teamId" binding:"required"`
	Name       string `form:"name" json:"name" binding:"required,max=100"`
	Repository string `form:"repository" json:"repository" binding:"required"`
	Branch     string `form:"branch" json:"branch" binding:"required"`
}

// GetProjects struct handle requirements to get projects
type GetProjects struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// Create handle requirements to create projects with Projects struct
func Create(c *gin.Context) {
	var (
		p Projects
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

// Read handle requirements to read projects with GetProjects struct
func Read(c *gin.Context) {
	var (
		p GetProjects
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
