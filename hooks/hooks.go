// Package hooks will manage all hooks requirements
package hooks

import (
	"crypto/md5"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/Lord-Y/cypress-parallel/git"
	"github.com/Lord-Y/cypress-parallel/kubernetes"
	"github.com/Lord-Y/cypress-parallel/models"
	"github.com/Lord-Y/cypress-parallel/tools"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// plain struct handle requirements to start unit testing
type plain struct {
	ProjectName          string `form:"project_name" json:"project_name" binding:"required,max=100"`
	Branch               string `form:"branch" json:"branch"`
	Specs                string `form:"specs" json:"specs"`
	ConfigFile           string `form:"config_file,default=cypress.config.js" json:"config_file" binding:"max=100"`
	Browser              string `form:"browser,default=chrome" json:"browser" binding:"max=100,oneof=chrome firefox"`
	MaxPods              int    `form:"maxPods,default=10" json:"maxPods"`
	CypressDockerVersion string `form:"cypress_docker_version,default=10.10.0-0.3.0,max=20" json:"cypress_docker_version"`
}

// execution handle all requirements to insert execution in DB
type execution struct {
	projectID       int
	uniqID          string
	branch          string
	executionStatus string // must be, NOT_STARTED, QUEUED, SCHEDULED, RUNNING, CANCELLED, FAILED, DONE
	spec            string
	result          string
}

// updatePodName will be used to update pod name in DB
type updatePodName struct {
	podName string
	uniqID  string
	spec    string
}

// dbProject struct permit to map data from db
type dbProject struct {
	Project_id             int       `json:"project_id"`
	Project_name           string    `json:"project_name"`
	Date                   time.Time `json:"date"`
	Team_id                int       `json:"team_id"`
	Repository             string    `json:"repository"`
	Branch                 string    `json:"branch"`
	Specs                  string    `json:"specs"`
	Scheduling             string    `json:"scheduling"`
	Scheduling_enabled     bool      `json:"scheduling_enabled"`
	Max_pods               int       `json:"max_pods"`
	Cypress_docker_version string    `json:"cypress_docker_version"`
	Timeout                int       `json:"timeout"`
	Username               string    `json:"username"`
	Password               string    `json:"password"`
	Browser                string    `json:"browser"`
	Config_file            string    `json:"config_file"`
}

// dbProjectAnnotation struct permit to map data from db
type dbProjectAnnotation struct {
	Annotation_id int       `json:"annotation_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	Project_id    int       `json:"project_id"`
	Date          time.Time `json:"date"`
}

// dbProjectEnvironment struct permit to map data from db
type dbProjectEnvironment struct {
	Environment_id int       `json:"environment_id"`
	Key            string    `json:"key"`
	Value          string    `json:"value"`
	Project_id     int       `json:"project_id"`
	Date           time.Time `json:"date"`
}

// dbProjectUniqIDExecution struct permit to map data from db
type dbProjectUniqIDExecution struct {
	Uniq_id string `json:"uniq_id"`
}

// dbPGQueue struct permit to map data from db
type dbPGQueue struct {
	Execution_id           int       `json:"execution_id"`
	Project_id             int       `json:"project_id"`
	Branch                 string    `json:"branch"`
	Execution_status       string    `json:"execution_status"`
	Uniq_id                int       `json:"uniq_id"`
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
	Count int `json:"count"`
}

var (
	commonLabels = map[string]string{
		"worker": "kubernetes",
		"app":    "cypress-parallel-jobs",
	}
)

const (
	ghr = "ghcr.io/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images"
)

// Plain handle requirements to start unit testing
func Plain(c *gin.Context) {
	var (
		p           plain
		ex          execution
		pod         models.Pods
		gitc        git.Repository
		targetSpecs string
		branch      string
		specs       []string
		tmpSecs     []string
		finalSecs   []string
		nbSpec      int
		reset       bool
	)
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if p.CypressDockerVersion == "" {
		p.CypressDockerVersion = "10.10.0-0.3.0"
	}
	if p.ConfigFile == "" {
		p.ConfigFile = "cypress.config.js"
	}
	if p.MaxPods == 0 {
		p.MaxPods = 10
	}

	result, err := p.getProjectInfos()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// original branch must remain for POST "executions" with update db otherwise, pod will stay up forever
	if p.Branch != "" {
		if p.Branch == "master" {
			gitc.Ref = ""
			branch = "master"
		} else {
			gitc.Ref = p.Branch
			branch = p.Branch
		}
	} else {
		if result.Branch == "master" {
			gitc.Ref = ""
			branch = "master"
		} else {
			gitc.Ref = result.Branch
			branch = result.Branch
		}
	}
	gitc.Repository = result.Repository
	gitc.Username = result.Username
	gitc.Password = result.Password

	gitdir, statusCode, err := gitc.Clone()
	defer os.RemoveAll(gitdir)
	if err != nil {
		msg := fmt.Sprintf("Error occured while cloning git repository, error: %s", err.Error())
		if statusCode == 400 {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if p.Specs != "" {
		targetSpecs = p.Specs
	} else {
		targetSpecs = result.Specs
	}

	if strings.HasSuffix(targetSpecs, ".spec.js") || strings.HasSuffix(targetSpecs, ".cy.js") || strings.HasSuffix(targetSpecs, ".spec.ts") || strings.HasSuffix(targetSpecs, ".cy.ts") {
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
			if info.Mode().IsRegular() && (strings.HasSuffix(file, ".spec.js") || strings.HasSuffix(file, ".cy.js") || strings.HasSuffix(file, ".spec.ts") || strings.HasSuffix(file, ".cy.ts")) {
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

	for k, v := range specs {
		tmpSecs = append(tmpSecs, v)
		if (nbSpec + 1) == commons.GetMaxSpecs() {
			finalSecs = append(finalSecs, strings.Join(tmpSecs, ","))
			tmpSecs = nil
			nbSpec = 0
			reset = true
		}
		if !reset {
			nbSpec++
		}
		if nbSpec == 0 {
			reset = false
		}
		if k == len(specs)-1 && commons.GetMaxSpecs()%2 == 1 {
			finalSecs = append(finalSecs, strings.Join(tmpSecs, ","))
		}
	}

	uniqID := md5.Sum([]byte(fmt.Sprintf("%s%s%s", result.Repository, finalSecs, time.Now())))
	runidID := fmt.Sprintf("%x", uniqID)
	uniqID_ := runidID[0:10]

	for count, spec := range finalSecs {
		var (
			pdn     updatePodName
			tag     string
			command []string
		)
		if count < p.MaxPods {
			for _, splittedSpec := range strings.Split(spec, ",") {
				ex.projectID = result.Project_id
				ex.uniqID = uniqID_
				ex.executionStatus = "NOT_STARTED"
				ex.spec = splittedSpec
				ex.result = `{}`
				ex.branch = branch

				_, err = ex.create()
				if err != nil {
					log.Error().Err(err).Msg("Error occured while performing db query")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					return
				}
			}

			annotations, err := result.getProjectAnnotations()
			if err != nil {
				log.Error().Err(err).Msg("Error occured while performing db query")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			if len(annotations) > 0 {
				annotation := make(map[string]string)
				for _, k := range annotations {
					annotation[k.Key] = k.Value
				}
				pod.Annotations = annotation
			}

			envVars, err := result.getProjectEnvironments()
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
					envVar.Key = fmt.Sprintf("CYPRESS_%s", k.Key)
					envVar.Value = k.Value
					envs = append(envs, envVar)
				}
				if strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL")) != "" {
					envVar.Key = "CYPRESS_PARALLEL_CLI_LOG_LEVEL"
					envVar.Value = strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL"))
					envs = append(envs, envVar)
				}
				envVar.Key = "NO_COLOR"
				envVar.Value = "1"
				envs = append(envs, envVar)
				pod.Container.EnvironmentVars = envs
			}
			pod.Namespace = commons.GetKubernetesJobsNamespace()
			pod.GenerateName = "cypress-parallel-jobs-"
			pod.Labels = commonLabels

			command = append(command, "cypress-parallel-cli")
			command = append(command, "cypress")
			command = append(command, "--browser")
			command = append(command, p.Browser)
			command = append(command, "--config-file")
			command = append(command, p.ConfigFile)
			command = append(command, "--specs")
			command = append(command, spec)
			command = append(command, "--uid")
			command = append(command, uniqID_)
			command = append(command, "--branch")
			command = append(command, branch)
			command = append(command, "--repository")
			command = append(command, gitc.Repository)
			command = append(command, "--api-url")
			command = append(command, commons.GetAPIUrl())
			command = append(command, "--report-back")
			command = append(command, "--timeout")
			command = append(command, strconv.Itoa(result.Timeout))
			if result.Username != "" {
				command = append(command, "--username")
				command = append(command, result.Username)
			}
			if result.Password != "" {
				command = append(command, "--password")
				command = append(command, result.Password)
			}
			if p.CypressDockerVersion != "" {
				tag = p.CypressDockerVersion
			} else {
				tag = result.Cypress_docker_version
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

			for _, splittedSpec := range strings.Split(spec, ",") {
				pdn.podName = podName
				pdn.uniqID = uniqID_
				pdn.spec = splittedSpec

				err = pdn.update()
				if err != nil {
					log.Error().Err(err).Msg("Error occured while performing update db query")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					return
				}
			}
		} else {
			for _, splittedSpec := range strings.Split(spec, ",") {
				ex.projectID = result.Project_id
				ex.uniqID = uniqID_
				ex.executionStatus = "QUEUED"
				ex.spec = splittedSpec
				ex.result = `{}`
				ex.branch = branch

				_, err = ex.create()
				if err != nil {
					log.Error().Err(err).Msg("Error occured while performing db query")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					return
				}
			}
		}
	}
	c.JSON(http.StatusCreated, "OK")
}

func Queued() {
	log.Debug().Msg("Checking QUEUED execution list")
	status, err := executionStatus("RUNNING")
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		return
	}

	if len(status) > 0 {
		wg := sync.WaitGroup{}
		log.Debug().Msgf("execution status running %+v length %d", status, len(status))
		for _, run := range status {
			wg.Add(1)
			go func(run dbProjectUniqIDExecution) {
				defer wg.Done()
				queuing(run)
			}(run)
		}
	} else {
		status, err := executionStatus("QUEUED")
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			return
		}
		if len(status) > 0 {
			wg := sync.WaitGroup{}
			log.Debug().Msgf("execution status queued %+v length %d", status, len(status))
			for _, run := range status {
				wg.Add(1)
				go func(run dbProjectUniqIDExecution) {
					defer wg.Done()
					queuing(run)
				}(run)
			}
		}
	}
}

func queuing(run dbProjectUniqIDExecution) {
	var (
		p plain
	)

	clientset, err := kubernetes.Client()
	if err != nil {
		log.Error().Err(err).Msg("Error occured while initializing kubernetes client")
		return
	}
	err = kubernetes.GetNamespace(clientset, commons.GetKubernetesJobsNamespace())
	if err != nil {
		log.Warn().Err(err).Msg("Error occured while getting kubernetes namespace")
		err = kubernetes.CreateNamespace(clientset, commons.GetKubernetesJobsNamespace())
		if err != nil {
			log.Error().Err(err).Msg("Error occured while creating kubernetes namespace")
			return
		}
	}
	err = kubernetes.GetServiceAccountName(clientset, commons.GetKubernetesJobsNamespace(), commons.GetKubernetesJobsNamespace())
	if err != nil {
		log.Warn().Err(err).Msgf("Error occured while getting kubernetes service account %s", commons.GetKubernetesJobsNamespace())
		_, err = kubernetes.CreateServiceAccountName(clientset, commons.GetKubernetesJobsNamespace(), commons.GetKubernetesJobsNamespace())
		if err != nil {
			log.Error().Err(err).Msgf("Error occured while creating kubernetes service account %s", commons.GetKubernetesJobsNamespace())
			return
		}
	}

	uniqID := run.Uniq_id
	countExecution, err := countExecutions(uniqID)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		return
	}

	queued, err := pgqueued(uniqID)
	if err != nil {
		log.Error().Err(err).Msg("Error occured while performing db query")
		return
	}

	log.Debug().Msgf("queued %s length length %d", uniqID, len(queued))
	if len(queued) > 0 {
		var (
			tmpSecs   []string
			finalSecs []string
			nbSpec    int
			reset     bool
			pdn       updatePodName
			command   []string
			pod       models.Pods
		)
		p.ProjectName = queued[0].Project_name
		p.Branch = queued[0].Branch

		for k, v := range queued {
			tmpSecs = append(tmpSecs, v.Spec)
			if (nbSpec + 1) == commons.GetMaxSpecs() {
				finalSecs = append(finalSecs, strings.Join(tmpSecs, ","))
				tmpSecs = nil
				nbSpec = 0
				reset = true
			}
			if !reset {
				nbSpec++
			}
			if nbSpec == 0 {
				reset = false
			}
			if k == len(queued)-1 && commons.GetMaxSpecs()%2 == 1 {
				finalSecs = append(finalSecs, strings.Join(tmpSecs, ","))
			}
		}
		log.Debug().Msgf("queued %s finalSecs %s", uniqID, finalSecs)

		result, err := p.getProjectInfos()
		if err != nil {
			log.Error().Err(err).Msg("Error occured while performing db query")
			return
		}

		count := int(math.Floor(float64(countExecution.Count) / float64(commons.GetMaxSpecs())))

		log.Debug().Msgf("queued %s running pods count %d VS max pods %d", uniqID, count, result.Max_pods)
		if count < result.Max_pods {
			annotations, err := result.getProjectAnnotations()
			if err != nil {
				log.Error().Err(err).Msg("Error occured while performing db query")
				return
			}
			if len(annotations) > 0 {
				annotation := make(map[string]string)
				for _, k := range annotations {
					annotation[k.Key] = k.Value
				}
				pod.Annotations = annotation
			}

			envVars, err := result.getProjectEnvironments()
			if err != nil {
				log.Error().Err(err).Msg("Error occured while performing db query")
				return
			}
			if len(envVars) > 0 {
				var (
					envs   []models.EnvironmentVar
					envVar models.EnvironmentVar
				)
				for _, k := range envVars {
					envVar.Key = fmt.Sprintf("CYPRESS_%s", k.Key)
					envVar.Value = k.Value
					envs = append(envs, envVar)
				}
				if strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL")) != "" {
					envVar.Key = "CYPRESS_PARALLEL_CLI_LOG_LEVEL"
					envVar.Value = strings.TrimSpace(os.Getenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL"))
					envs = append(envs, envVar)
				}
				envVar.Key = "NO_COLOR"
				envVar.Value = "1"
				envs = append(envs, envVar)
				pod.Container.EnvironmentVars = envs
			}
			pod.Namespace = commons.GetKubernetesJobsNamespace()
			pod.GenerateName = "cypress-parallel-jobs-"
			pod.Labels = commonLabels

			command = append(command, "cypress-parallel-cli")
			command = append(command, "cypress")
			command = append(command, "--browser")
			command = append(command, result.Browser)
			command = append(command, "--config-file")
			command = append(command, result.Config_file)
			command = append(command, "--specs")
			command = append(command, finalSecs[0])
			command = append(command, "--uid")
			command = append(command, uniqID)
			command = append(command, "--branch")
			command = append(command, p.Branch)
			command = append(command, "--repository")
			command = append(command, result.Repository)
			command = append(command, "--api-url")
			command = append(command, commons.GetAPIUrl())
			command = append(command, "--report-back")
			command = append(command, "--timeout")
			command = append(command, strconv.Itoa(result.Timeout))
			if result.Username != "" {
				command = append(command, "--username")
				command = append(command, result.Username)
			}
			if result.Password != "" {
				command = append(command, "--password")
				command = append(command, result.Password)
			}
			tag := result.Cypress_docker_version

			pod.Container.Command = command
			pod.Container.Name = "cypress-parallel-jobs"
			pod.Container.Image = fmt.Sprintf("%s:%s", ghr, tag)

			podName, err := kubernetes.CreatePod(clientset, pod)
			if err != nil {
				log.Error().Err(err).Msg("Error occured while creating pod")
				return
			}
			log.Debug().Msgf("Pod name %s created for specs %s", podName, finalSecs[0])

			for _, splittedSpec := range strings.Split(finalSecs[0], ",") {
				pdn.podName = podName
				pdn.uniqID = uniqID
				pdn.spec = splittedSpec

				err = pdn.update()
				if err != nil {
					log.Error().Err(err).Msg("Error occured while performing update db query")
					return
				}
			}
		}
	}
}
