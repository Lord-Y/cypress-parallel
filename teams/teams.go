// Package teams will manage all teams requirements
package teams

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// teams struct handle requirements to create teams
type teams struct {
	Name string `form:"name" json:"name" binding:"required,max=100"`
}

// getTeams struct handle requirements to get teams
type getTeams struct {
	TeamID int `form:"teamId" json:"teamId" binding:"required"`
}

// listTeams struct handle requirements to get teams
type listTeams struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// updateTeam struct handle requirements to update teams
type updateTeam struct {
	TeamID int    `form:"teamId" json:"teamId" binding:"required"`
	Name   string `form:"name" json:"name" binding:"required,max=100"`
}

// deleteTeam struct handle requirements to delete team
type deleteTeam struct {
	TeamID int `form:"teamId" json:"teamId" binding:"required"`
}

// searchTeams struct handle requirements to get teams
type searchTeams struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// dbList struct permit to map data from db
type dbList struct {
	Team_id   int       `json:"team_id"`
	Team_name string    `json:"team_name"`
	Date      time.Time `json:"date"`
	Total     int       `json:"total"`
}

// dbCommon struct permit to map data from db
type dbCommon struct {
	Team_id   int       `json:"team_id"`
	Team_name string    `json:"team_name"`
	Date      time.Time `json:"date"`
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
		c.JSON(http.StatusCreated, gin.H{"teamId": result})
	}
}

// Read handle requirements to read teams with getTeams struct
func Read(c *gin.Context) {
	var (
		p getTeams
	)
	id := c.Params.ByName("teamId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.TeamID = vID

	result, err := p.read()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if result == (dbCommon{}) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// List handle requirements to read teams with listTeams struct
func List(c *gin.Context) {
	var (
		p listTeams
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

// Update handle requirements to update teams with updateTeam struct
func Update(c *gin.Context) {
	var (
		p updateTeam
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

// Delete handle deletion of project deleteTeam struct
func Delete(c *gin.Context) {
	var (
		p deleteTeam
	)
	id := c.Params.ByName("teamId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "teamId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.TeamID = vID

	err = p.delete()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

// Search handle requirements to search teams with searchTeams struct
func Search(c *gin.Context) {
	var (
		p searchTeams
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
