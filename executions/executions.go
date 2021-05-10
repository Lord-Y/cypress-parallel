// Package executions will manage all executions requirements
package executions

import (
	"net/http"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// listExecutions struct handle requirements to get projects
type listExecutions struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// readExecutions struct handle requirements to get projects
type readExecutions struct {
	ExecutionID int `form:"executionId" json:"executionId" binding:"required"`
}

// updateResultExecution permit to update execution result
type updateResultExecution struct {
	UniqID               string `form:"uniqId" json:"uniqId" binding:"required"`
	Spec                 string `form:"spec" json:"spec" binding:"required"`
	Branch               string `form:"branch" json:"branch"`
	Result               string `form:"result" json:"result" binding:"required"`
	ExecutionStatus      string `form:"executionStatus" json:"executionStatus" binding:"required"`
	ExecutionErrorOutput string `form:"executionErrorOutput" json:"executionErrorOutput"`
}

// List permit to retrieve executions with pagination
func List(c *gin.Context) {
	var (
		p listExecutions
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

// Read permit to read content of specific execution
func Read(c *gin.Context) {
	var (
		p readExecutions
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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

// UpdateResultExecution permit update specific execution result
func UpdateResultExecution(c *gin.Context) {
	var (
		p updateResultExecution
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.updateResult()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusOK, "OK")
	}
}
