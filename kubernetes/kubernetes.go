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
	if commons.GetKubernetesKubeConfig() != "" {
		kconfig = commons.GetKubernetesKubeConfig()
	} else {
		if home := homedir.HomeDir(); home != "" {
			kconfig = filepath.Join(home, ".kube", "config")
		}
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
	_, err = clientset.
		CoreV1().
		Namespaces().
		Get(
			context.TODO(),
			namespace,
			metav1.GetOptions{},
		)
	return
}

// CreateNamespace permit to create namespace
func CreateNamespace(clientset *kubernetes.Clientset, namespace string) (err error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err = clientset.
		CoreV1().
		Namespaces().
		Create(
			context.TODO(),
			ns,
			metav1.CreateOptions{},
		)
	return
}

// DeleteNamespace permit to delete created namespace
func DeleteNamespace(clientset *kubernetes.Clientset, namespace string) (err error) {
	err = clientset.
		CoreV1().
		Namespaces().
		Delete(
			context.TODO(),
			namespace,
			metav1.DeleteOptions{},
		)
	return
}

// GetServiceAccountName permit to get service account in specified namespace
func GetServiceAccountName(clientset *kubernetes.Clientset, namespace string, serviceAccount string) (err error) {
	_, err = clientset.
		CoreV1().
		ServiceAccounts(namespace).
		Get(
			context.TODO(),
			serviceAccount,
			metav1.GetOptions{},
		)
	return
}

// CreateServiceAccountName permit to create service account that will be used while creating the pod
func CreateServiceAccountName(clientset *kubernetes.Clientset, namespace string, serviceAccount string) (serviceAccountName string, err error) {
	sa := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccount,
			Namespace: namespace,
		},
	}
	result, err := clientset.
		CoreV1().
		ServiceAccounts(namespace).
		Create(
			context.TODO(),
			sa,
			metav1.CreateOptions{},
		)
	if err != nil {
		return "", err
	}
	return result.Name, nil
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
			ServiceAccountName:            m.Namespace,
			Affinity: &v1.Affinity{
				PodAntiAffinity: &v1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
						{
							Weight: int32(5),
							PodAffinityTerm: v1.PodAffinityTerm{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"worker": "kubernetes",
									},
								},
								Namespaces: []string{
									m.Namespace,
								},
								TopologyKey: "kubernetes.io/hostname",
							},
						},
					},
				},
			},
		},
	}
	result, err := clientset.
		CoreV1().
		Pods(m.Namespace).
		Create(
			context.TODO(),
			pod,
			metav1.CreateOptions{},
		)
	if err != nil {
		return "", err
	}
	return result.Name, nil
}

// DeletePod permit to delete pod inside of specified namespace
func DeletePod(clientset *kubernetes.Clientset, namespace string, podName string) (err error) {
	err = clientset.
		CoreV1().
		Pods(namespace).
		Delete(
			context.TODO(),
			podName,
			metav1.DeleteOptions{},
		)
	return
}
