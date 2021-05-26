// Package annotations will manage all annotations requirements that will be injected in pods
package annotations

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

// annotations struct handle requirements to create annotations
type annotations struct {
	ProjectID   int    `form:"projectId" json:"projectId" binding:"required"`
	Annotations string `form:"annotations" json:"annotations" binding:"required"`
	annotation  []annotation
}

// annotation struct handle k/v from annotations
type annotation struct {
	Annotation_ID int
	Key           string
	Value         string
}

// getAnnotations struct handle requirements to get annotations
type getAnnotations struct {
	AnnotationID int `form:"annotationId" json:"annotationId" binding:"required"`
}

// deleteAnnotations struct handle requirements to create annotations
type deleteAnnotation struct {
	AnnotationID int `form:"annotationId" json:"annotationId" binding:"required"`
}

// listAnnotations struct handle requirements to get projects
type listAnnotations struct {
	Page       int `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// searchAnnotations struct handle requirements to get projects
type searchAnnotations struct {
	Q          string `form:"q" json:"q" binding:"required"`
	Page       int    `form:"page,default=1" json:"page"`
	RangeLimit int
	StartLimit int
	EndLimit   int
}

// CreateOrUpdate handle requirements to create annotations with annotations struct
func CreateOrUpdate(c *gin.Context) {
	var (
		p annotations
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var m []map[string]interface{}
	err := json.Unmarshal([]byte(p.Annotations), &m)
	if err != nil {
		log.Error().Err(err).Msgf("Error occured while unmarshalling data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	mapstructure.Decode(m, &p.annotation)

	// pg UPSERT cannot be used for our purpose
	for i, k := range p.annotation {
		if k.Annotation_ID == 0 {
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

// List handle requirements to read projects with getProjects struct
func List(c *gin.Context) {
	var (
		p listAnnotations
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

// Read handle requirements to read annotations with getAnnotations struct
func Read(c *gin.Context) {
	var (
		p getAnnotations
	)
	id := c.Params.ByName("annotationId")
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

	p.AnnotationID = vID

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

// Delete handle deletion of annotation viariable with deleteAnnotation struct
func Delete(c *gin.Context) {
	var (
		p deleteAnnotation
	)
	id := c.Params.ByName("annotationId")
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

	p.AnnotationID = vID

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
		p searchAnnotations
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
