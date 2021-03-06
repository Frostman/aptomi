package lang

import (
	"github.com/Aptomi/aptomi/pkg/object"
)

// ClusterObject is an informational data structure with Kind and Constructor for Cluster
var ClusterObject = &object.Info{
	Kind:        "cluster",
	Versioned:   true,
	Constructor: func() object.Base { return &Cluster{} },
}

// Cluster defines an individual cluster where containers get deployed.
// Various cloud providers are supported via setting a cluster type (k8s, Amazon ECS, GKE, etc).
type Cluster struct {
	Metadata

	// Type is a cluster type. Based on its type, the appropriate deployment plugin will be called to deploy containers.
	Type string `validate:"clustertype"`

	// Labels is a set of labels assigned to the cluster
	Labels map[string]string

	// Config for a given cluster type
	Config struct {
		KubeContext     string
		TillerNamespace string
		Namespace       string
	}
}
