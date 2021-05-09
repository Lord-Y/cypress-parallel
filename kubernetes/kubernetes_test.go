// Package kubernetes will manage all kubernetes pods requirements
package kubernetes

import (
	"os"
	"testing"
	"time"

	"github.com/Lord-Y/cypress-parallel-api/models"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestGetNamespace(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE", "true")
	defer os.Unsetenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE")
	namespace := fake.CharactersN(10)
	client, err := Client()
	err = GetNamespace(client, namespace)
	assert.Error(err)
}

func TestCreateDeleteNamespace(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE", "true")
	defer os.Unsetenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE")
	namespace := fake.CharactersN(10)
	client, err := Client()
	err = CreateNamespace(client, namespace)
	assert.NoError(err)

	err = DeleteNamespace(client, namespace)
	assert.NoError(err)
}

func TestCreatePod(t *testing.T) {
	assert := assert.New(t)
	var (
		pod    models.Pods
		envs   []models.EnvironmentVar
		envVar models.EnvironmentVar
	)

	os.Setenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE", "true")
	defer os.Unsetenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE")
	name := fake.CharactersN(10)

	client, err := Client()
	err = CreateNamespace(client, name)
	assert.NoError(err)

	err = GetServiceAccountName(client, name, name)
	assert.Error(err)

	_, err = CreateServiceAccountName(client, name, name)
	assert.NoError(err)

	pod.GenerateName = fake.CharactersN(5)
	pod.Container.Name = "alpine"
	pod.Container.Image = "alpine:latest"
	pod.Namespace = name
	pod.Container.Command = []string{
		"ls",
	}
	envVar.Key = "key"
	envVar.Value = "value"
	envs = append(envs, envVar)
	pod.Container.EnvironmentVars = envs

	podName, err := CreatePod(client, pod)
	assert.NoError(err)

	time.Sleep(time.Duration(30) * time.Second)
	err = DeletePod(client, name, podName)
	assert.NoError(err)

	time.Sleep(time.Duration(10) * time.Second)
	err = DeleteNamespace(client, name)
	assert.NoError(err)
}

func TestClient_fail_client(t *testing.T) {
	os.Unsetenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE")
	defer func() { recover() }()
	name := fake.CharactersN(10)

	client, _ := Client()
	_ = CreateNamespace(client, name)
	t.Errorf("Code did not panic")
}

func TestClient_fail_client_kubeconfig(t *testing.T) {
	os.Setenv("CYPRESS_PARALLEL_API_K8S_KUBE_CONFIG", os.TempDir())
	defer os.Unsetenv("CYPRESS_PARALLEL_API_K8S_KUBE_CONFIG")
	os.Setenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE", "true")
	defer os.Unsetenv("CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE")

	defer func() { recover() }()
	name := fake.CharactersN(10)

	client, _ := Client()
	_ = CreateNamespace(client, name)
	t.Errorf("Code did not panic")
}
