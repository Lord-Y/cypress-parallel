// Package kubernetes will manage all kubernetes pods requirements
package kubernetes

import (
	"context"
	"path/filepath"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	"github.com/Lord-Y/cypress-parallel-api/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client return requirements to be able to connect to kubernetes cluster
// with the program running inside or outside of the cluster
func Client() (c *kubernetes.Clientset, err error) {
	if commons.GetKubernetesMode() == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return c, err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return c, err
		}
		return clientset, nil
	}

	var kconfig string
	if home := homedir.HomeDir(); home != "" {
		kconfig = filepath.Join(home, ".kube", "config")
	} else {
		kconfig = commons.GetKubernetesKubeConfig()
	}

	kubeconfig := &kconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return c, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return c, err
	}
	return clientset, nil
}

// GetNamespace check if namespace defined for jobs exist or not
func GetNamespace(clientset *kubernetes.Clientset, namespace string) (err error) {
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), commons.GetKubernetesJobsNamespace(), metav1.GetOptions{})
	return err
}

// CreateNamespace permit to create namespace
func CreateNamespace(clientset *kubernetes.Clientset) (err error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: commons.GetKubernetesJobsNamespace(),
		},
	}
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return
}

// CreatePod permit to create pod inside of specified namespace
func CreatePod(clientset *kubernetes.Clientset, m models.Pods) (podName string, err error) {
	var (
		env  v1.EnvVar
		envs []v1.EnvVar
	)
	terminationGracePeriodSeconds := int64(300)

	if len(m.Container.EnvironmentVars) > 0 {
		for _, k := range m.Container.EnvironmentVars {
			env.Name = k.Key
			env.Value = k.Value
			envs = append(envs, env)
		}
	}
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: m.GenerateName,
			Namespace:    m.Namespace,
			Labels:       m.Labels,
			Annotations:  m.Annotations,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            m.Container.Name,
					Image:           m.Container.Image,
					Command:         m.Container.Command,
					Env:             envs,
					ImagePullPolicy: v1.PullIfNotPresent,
				},
			},
			RestartPolicy:                 v1.RestartPolicyNever,
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
		},
	}
	result, err := clientset.CoreV1().Pods(m.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return result.Name, nil
}
