package lang

import (
	"github.com/Aptomi/aptomi/pkg/object"
)

// DependencyObject is an informational data structure with Kind and Constructor for Dependency
var DependencyObject = &object.Info{
	Kind:        "dependency",
	Versioned:   true,
	Constructor: func() object.Base { return &Dependency{} },
}

// Dependency is a declaration of use, defined in a form <User> needs an instance of <Contract> with
// specified set of <Labels>. It allows users to request contracts, which will translate into instantiation of
// service instances (and their dependencies) in the cloud
type Dependency struct {
	Metadata

	User     string `validate:"required"`
	Contract string `validate:"required"`
	Labels   map[string]string
}

// GlobalDependencies represents the list of global dependencies (see the definition above)
type GlobalDependencies struct {
	// DependencyMap is a map[name] -> *Dependency
	DependencyMap map[string]*Dependency

	// DependenciesByContract contains dependency map <contractName> -> list of dependencies
	DependenciesByContract map[string][]*Dependency
}

// NewGlobalDependencies creates and initializes a new empty list of global dependencies
func NewGlobalDependencies() *GlobalDependencies {
	return &GlobalDependencies{
		DependencyMap:          make(map[string]*Dependency),
		DependenciesByContract: make(map[string][]*Dependency),
	}
}

// addDependency appends a single dependency to an existing object
func (globalDependencies GlobalDependencies) addDependency(dependency *Dependency) {
	globalDependencies.DependencyMap[dependency.GetName()] = dependency
	globalDependencies.DependenciesByContract[dependency.Contract] = append(globalDependencies.DependenciesByContract[dependency.Contract], dependency)
}
