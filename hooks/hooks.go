// Package hooks will manage all hooks requirements
package hooks

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/git"
	"github.com/Lord-Y/cypress-parallel-api/kubernetes"
	"github.com/Lord-Y/cypress-parallel-api/models"
	"github.com/Lord-Y/cypress-parallel-api/tools"
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
	CypressDockerVersion string `form:"cypress_docker_version,default=7.2.0-0.0.1,max=20" json:"cypress_docker_version,max=20"`
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

// environmentVar k/v to set inside of the container
type environmentVar struct {
	Key   string // Variable key
	Value string // Variable value
}

// updatePodName will be used to update pod name in DB
type updatePodName struct {
	podName string
	uniqID  string
	spec    string
}

const (
	ghr = "docker.pkg.github.com/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images"
)

// Plain handle requirements to start unit testing
func Plain(c *gin.Context) {
	var (
		p           plain
		pj          projects
		exs         executions
		ex          execution
		pod         models.Pods
		gitc        git.Repository
		targetSpecs string
		specs       []string
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

	if p.Branch != "" {
		if p.Branch == "master" {
			gitc.Branch = ""
		} else {
			gitc.Branch = p.Branch
		}
	} else {
		if pj.Branch == "master" {
			gitc.Branch = ""
		} else {
			gitc.Branch = pj.Branch
		}
	}
	gitc.Repository = pj.Repository

	gitdir, statusCode, err := gitc.Clone()
	defer os.RemoveAll(gitdir)
	if err != nil {
		msg := fmt.Sprintf("Error occured while cloning git repository, error: %s", err.Error())
		if statusCode == 400 {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	if p.Specs != "" {
		targetSpecs = p.Specs
	} else {
		targetSpecs = pj.Specs
	}

	if strings.HasSuffix(targetSpecs, ".spec.js") || strings.HasSuffix(targetSpecs, ".ts") {
		err := tools.CheckIsFile(fmt.Sprintf("%s/%s", gitdir, targetSpecs))
		if err != nil {
			msg := fmt.Sprintf("Error occured while retrieving specs, error: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		specs = append(specs, targetSpecs)
	} else {
		err := filepath.Walk(fmt.Sprintf("%s/%s", gitdir, targetSpecs), func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() && (strings.HasSuffix(file, ".spec.js") || strings.HasSuffix(file, ".ts")) {
				specs = append(specs, strings.ReplaceAll(file, fmt.Sprintf("%s/", gitdir), ""))
			}
			return nil
		})
		if err != nil {
			msg := fmt.Sprintf("Error occured while retrieving specs, error: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
	}

	clientset, err := kubernetes.Client()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while initializing kubernetes client")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	err = kubernetes.GetNamespace(clientset, commons.GetKubernetesJobsNamespace())
	if err != nil {
		log.Warn().Err(err).Msg("Error occured while getting kubernetes namespace")
		err = kubernetes.CreateNamespace(clientset, commons.GetKubernetesJobsNamespace())
		if err != nil {
			log.Error().Err(err).Msg("Error occured while creating kubernetes namespace")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	err = kubernetes.GetServiceAccountName(clientset, commons.GetKubernetesJobsNamespace(), commons.GetKubernetesJobsNamespace())
	if err != nil {
		log.Warn().Err(err).Msgf("Error occured while getting kubernetes service account %s", commons.GetKubernetesJobsNamespace())
		_, err = kubernetes.CreateServiceAccountName(clientset, commons.GetKubernetesJobsNamespace(), commons.GetKubernetesJobsNamespace())
		if err != nil {
			log.Error().Err(err).Msgf("Error occured while creating kubernetes service account %s", commons.GetKubernetesJobsNamespace())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}

	for _, spec := range specs {
		var (
			pdn     updatePodName
			tag     string
			command []string
		)
		uniqID := md5.Sum([]byte(fmt.Sprintf("%s%s%s", pj.Repository, spec, time.Now())))
		runidID := fmt.Sprintf("%x", uniqID)
		uniqID_ := runidID[0:10]
		projecID, err := strconv.Atoi(pj.Project_id)
		if err != nil {
			log.Error().Err(err).Msg("Error occured while converting string to int")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ex.projectID = projecID
		ex.uniqID = uniqID_
		ex.executionStatus = "NOT_STARTED"
		ex.spec = spec
		ex.result = `{}`
		exs.executions = append(exs.executions, ex)

		_, err = ex.create()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		annotations, err := pj.getProjectAnnotations()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		if len(annotations) > 0 {
			annotation := make(map[string]string)
			for _, k := range annotations {
				annotation[fmt.Sprintf("%s", k["key"])] = fmt.Sprintf("%s", k["value"])
			}
			pod.Annotations = annotation
		}

		envVars, err := pj.getProjectEnvironments()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		if len(envVars) > 0 {
			var (
				envs   []models.EnvironmentVar
				envVar models.EnvironmentVar
			)
			for _, k := range envVars {
				envVar.Key = fmt.Sprintf("%s", k["key"])
				envVar.Value = fmt.Sprintf("%s", k["value"])
				envs = append(envs, envVar)
			}
			pod.Container.EnvironmentVars = envs
		}
		pod.Namespace = commons.GetKubernetesJobsNamespace()
		pod.GenerateName = "cypress-parallel-jobs-"
		pod.Labels = map[string]string{
			"worker": "kubernetes",
			"app":    "cypress-parallel-jobs",
		}

		command = append(command, "cypress-parallel-cli")
		command = append(command, "cypress")
		if p.Browser != "" {
			command = append(command, "--browser")
			command = append(command, p.Browser)
		}
		if p.ConfigFile != "" {
			command = append(command, "--cf")
			command = append(command, p.ConfigFile)
		}
		command = append(command, "--specs")
		command = append(command, spec)
		command = append(command, "--uid")
		command = append(command, uniqID_)
		command = append(command, "--branch")
		command = append(command, gitc.Branch)
		command = append(command, "--repository")
		command = append(command, gitc.Repository)
		command = append(command, "--api-url")
		if strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_API_URL")) == "" {
			command = append(command, "http://127.0.0.1:8080")
		} else {
			command = append(command, strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_API_URL")))
		}
		command = append(command, "--report-back")
		if p.CypressDockerVersion != "" {
			tag = p.CypressDockerVersion
		} else {
			tag = pj.CypressDockerVersion
		}

		pod.Container.Command = command
		pod.Container.Name = "cypress-parallel-jobs"
		pod.Container.Image = fmt.Sprintf("%s:%s", ghr, tag)

		podName, err := kubernetes.CreatePod(clientset, pod)
		if err != nil {
			log.Error().Err(err).Msg("Error occured while creating pod")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		log.Debug().Msgf("Pod name %s created", podName)

		pdn.podName = podName
		pdn.uniqID = uniqID_
		pdn.spec = spec

		err = pdn.update()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing update db query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusCreated, "OK")
}
