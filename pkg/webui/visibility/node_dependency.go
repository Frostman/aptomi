package visibility

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/external/users"
	"github.com/Aptomi/aptomi/pkg/lang"
	"github.com/Aptomi/aptomi/pkg/object"
)

type dependencyNode struct {
	dependency *lang.Dependency
	short      bool
	userLoader users.UserLoader
}

func newDependencyNode(dependency *lang.Dependency, short bool, userLoader users.UserLoader) graphNode {
	return dependencyNode{
		dependency: dependency,
		short:      short,
		userLoader: userLoader,
	}
}

func (n dependencyNode) getIDPrefix() string {
	return "dep-"
}

func (n dependencyNode) getGroup() string {
	if n.short {
		return "dependencyShort"
	}
	/*
		if n.dependency.Resolved {
			return "dependencyLongResolved"
		}
	*/
	return "dependencyLongNotResolved"
}

func (n dependencyNode) getID() string {
	return fmt.Sprintf("%s%s", n.getIDPrefix(), object.GetKey(n.dependency))
}

func (n dependencyNode) isItMyID(id string) string {
	return cutPrefixOrEmpty(id, n.getIDPrefix())
}

func (n dependencyNode) getLabel() string {
	if n.short {
		// for service owner view, don't display much other than a user name
		return n.dependency.User
	}
	// for consumer view - display full dependency info "user name -> contract"
	return fmt.Sprintf("%s \u2192 %s", n.dependency.User, n.dependency.Contract)
}

func (n dependencyNode) getEdgeLabel(dst graphNode) string {
	return ""
}

func (n dependencyNode) getDetails(id string /*, revision *resolve.Revision*/) interface{} {
	return nil //revision.Policy.Dependencies.DependenciesByID[id]
}
