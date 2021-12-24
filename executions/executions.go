// Package executions will manage all executions requirements
package executions

import (
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/kubernetes"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// listExecutions struct handle requirements to get executions
type listExecutions struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// readExecutions struct handle requirements to get executions
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
	Encoded              bool   `form:"encoded" json:"encoded"`
}

// searchExecutions struct handle requirements to get executions
type searchExecutions struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// uniqIDExecutions struct handle requirements to get uniq id executions
type uniqIDExecutions struct {
	UniqID string `form:"uniqId" json:"uniqId" binding:"required"`
}

// dbList struct permit to map data from db
type dbList struct {
	Execution_id           int       `json:"execution_id"`
	Project_id             int       `json:"project_id"`
	Branch                 string    `json:"branch"`
	Execution_status       string    `json:"execution_status"`
	Uniq_id                string    `json:"uniq_id"`
	Spec                   string    `json:"spec"`
	Result                 string    `json:"result"`
	Date                   time.Time `json:"date"`
	Execution_error_output string    `json:"execution_error_output"`
	Pod_name               string    `json:"pod_name"`
	Pod_cleaned            bool      `json:"pod_cleaned"`
	Total                  int       `json:"total"`
	Project_name           string    `json:"project_name"`
}

// DBRead struct permit to map data from db
type DBRead struct {
	Execution_id           int       `json:"execution_id"`
	Project_id             int       `json:"project_id"`
	Branch                 string    `json:"branch"`
	Execution_status       string    `json:"execution_status"`
	Uniq_id                string    `json:"uniq_id"`
	Spec                   string    `json:"spec"`
	Result                 string    `json:"result"`
	Date                   time.Time `json:"date"`
	Execution_error_output string    `json:"execution_error_output"`
	Pod_name               string    `json:"pod_name"`
	Pod_cleaned            bool      `json:"pod_cleaned"`
	Project_name           string    `json:"project_name"`
}

// dbCountExecutions struct permit to map data from db
type dbCountExecutions struct {
	Pod_name         string `json:"pod_name"`
	Execution_status string `json:"execution_status"`
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
	id := c.Params.ByName("executionId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "executionId is missing in uri"})
		return
	}
	vID, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while converting string to int")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	p.ExecutionID = vID
	result, err := p.read()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if result == (DBRead{}) {
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

	if p.Encoded {
		decoded, err := hex.DecodeString(p.Result)
		if err != nil {
			log.Error().Err(err).Msg("Error occured while decoding result")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		p.Result = string(decoded)
	}

	_, err := p.updateResult()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		c.JSON(http.StatusOK, "OK")
	}
	log.Debug().Msgf("POST body %+v", p)

	remaining, err := p.countExecutions()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	log.Debug().Msgf("remaining %+v", remaining)

	if remaining == (dbCountExecutions{}) {
		pod, err := p.countExecutionsInverted()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		if pod != (dbCountExecutions{}) {
			clientset, err := kubernetes.Client()
			if err != nil {
				log.Error().Err(err).Msg("Error occured while initializing kubernetes client")
				return
			}
			err = kubernetes.DeletePod(clientset, commons.GetKubernetesJobsNamespace(), pod.Pod_name)
			if err != nil {
				log.Error().Err(err).Msgf("Error occured while trying to delete pod name: %s", pod.Pod_name)
				return
			}
		}
	}
}

// Search handle requirements to search projects with searchExecutions struct
func Search(c *gin.Context) {
	var (
		p searchExecutions
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

// UniqID permit to get content of all uniq id executions
func UniqID(c *gin.Context) {
	var (
		p uniqIDExecutions
	)
	id := c.Params.ByName("uniqId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uniqId is missing in uri"})
		return
	}

	p.UniqID = id
	result, err := p.uniqId()
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
