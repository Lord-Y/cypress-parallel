// Package models assemble all struct used by other packages
package models

// Pods struct will be use by kooks package to create pods in kubernetes cluster
type Pods struct {
	GenerateName string            // GenerateName is the prefix that will be use to create the pod name
	Namespace    string            // Namespace in which the pod will be created
	Annotations  map[string]string // Annotations to set to the pod
	Labels       map[string]string // Labels to set to the pod
	Container    container         // Container requirements
}

// EnvironmentVar k/v to set inside of the container
type EnvironmentVar struct {
	Key   string // Variable key
	Value string // Variable value
}

// container hold the configuration that will be use to create pod
type container struct {
	Name            string           // Container name
	Image           string           // Docker image name
	Command         []string         // Command to run inside of the container
	EnvironmentVars []EnvironmentVar // Environments variables to set inside of the container
}
