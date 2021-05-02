// Package hooks will manage all hooks requirements
package hooks

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

// plain struct handle requirements to start unit testing
type plain struct {
	ProjectName          string `form:"project_name" json:"project_name" binding:"required,max=100"`
	Branch               string `form:"branch" json:"branch"`
	Specs                string `form:"specs" json:"specs"`
	ConfigFile           string `form:"config_file" json:"config_file"`
	Group                string `form:"group" json:"group"`
	Browser              string `form:"browser" json:"browser"`
	MaxPods              string `form:"max_pods" json:"max_pods"`
	CypressDockerVersion string `form:"cypress_docker_version,default=7.2.0,max=20" json:"cypress_docker_version,max=20"`
	plain                bool
}

// projects will be use to "mapstructure" data from db
type projects struct {
	Team_id              string
	Project_id           string
	Project_name         string
	Date                 string
	Repository           string
	Branch               string
	Specs                string
	Scheduling           string
	SchedulingEnabled    bool
	Max_pods             string
	CypressDockerVersion string
}

// executions is a slice of execution struct
type executions struct {
	executions []execution
}

// execution handle all requirements to insert execution in DB
type execution struct {
	projectID       int
	uniqID          string
	branch          string
	executionStatus string // must be, NOT_STARTED, SCHEDULED, RUNNING, CANCELLED, FAILED, DONE
	spec            string
	result          string
}

// Plain handle requirements to start unit testing
func Plain(c *gin.Context) {
	var (
		p   plain
		pj  projects
		exs executions
		ex  execution
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := p.getProjectInfos()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	mapstructure.Decode(result, &pj)

	specs, statusCode, err := pj.plainClone(p.Branch, p.Specs)
	if err != nil {
		msg := fmt.Sprintf("Error occured while retrieving specs, error: %s", err.Error())
		if statusCode == 400 {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	for _, spec := range specs {
		uniqID := md5.Sum([]byte(fmt.Sprintf("%s%s%s", pj.Repository, spec, time.Now())))
		runidID := fmt.Sprintf("%x", uniqID)
		projecID, err := strconv.Atoi(pj.Project_id)
		if err != nil {
			log.Error().Err(err).Msg("Error occured while converting string to int")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ex.projectID = projecID
		ex.uniqID = runidID[0:10]
		ex.executionStatus = "NOT_STARTED"
		ex.spec = spec
		ex.result = `{}`
		exs.executions = append(exs.executions, ex)

		err = ex.create()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusCreated, "OK")
}
