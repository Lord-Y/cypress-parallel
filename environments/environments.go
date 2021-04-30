// Package environments will manage all environments requirements that will be injected in pods
package environments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/tools"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// Environments struct handle requirements to create environments
type Environments struct {
	ProjectID    int    `form:"projectId" json:"projectId" binding:"required"`
	Environments string `form:"environments" json:"environments" binding:"required"`
	environment  []environment
}

// environment struct handle k/v from environments
type environment struct {
	Environment_ID int
	Key            string
	Value          string
}

// GetEnvironments struct handle requirements to get environments
type GetEnvironments struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// DeleteEnvironment struct handle requirements to delete environment
type DeleteEnvironment struct {
	EnvironmentID int `form:"environmentId" json:"environmentId" binding:"required"`
}

// CreateOrUpdate handle requirements to create environments with Environments struct
func CreateOrUpdate(c *gin.Context) {
	var (
		p Environments
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var m []map[string]interface{}
	err := json.Unmarshal([]byte(p.Environments), &m)
	if err != nil {
		log.Error().Err(err).Msgf("Error occured while unmarshalling data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	mapstructure.Decode(m, &p.environment)

	// pg UPSERT cannot be used for our purpose
	for i, k := range p.environment {
		if k.Environment_ID == 0 {
			resultSelect, err := p.selectBeforeAct(i)
			if err != nil {
				log.Error().Err(err).Msg("Error occured while performing select db query")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			if resultSelect["total"] == "0" {
				err = p.create(i)
				if err != nil {
					log.Error().Err(err).Msg("Error occured while performing create db query")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					return
				}
			}
		} else {
			err = p.update(i)
			if err != nil {
				log.Error().Err(err).Msg("Error occured while performing update db query")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
		}
	}
	c.JSON(http.StatusCreated, "OK")
}

// Read handle requirements to read environments with GetEnvironments struct
func Read(c *gin.Context) {
	var (
		p GetEnvironments
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

// Delete handle deletion of environment viariable with DeleteEnvironment struct
func Delete(c *gin.Context) {
	var (
		p DeleteEnvironment
	)
	id := c.Params.ByName("environmentId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "environmentId is missing in uri"})
		return
	}
	convID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.EnvironmentID = convID

	err = p.delete()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}
