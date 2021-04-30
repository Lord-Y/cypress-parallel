// Package teams will manage all teams requirements
package teams

import (
	"net/http"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// teams struct handle requirements to create teams
type teams struct {
	Name string `form:"name" json:"name" binding:"required,max=100"`
}

// getTeams struct handle requirements to get teams
type getTeams struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// Create handle requirements to create teams with teams struct
func Create(c *gin.Context) {
	var (
		p teams
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

// Read handle requirements to read teams with getTeams struct
func Read(c *gin.Context) {
	var (
		p getTeams
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
