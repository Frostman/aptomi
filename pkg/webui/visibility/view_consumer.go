package visibility

// ConsumerView represents a view from a particular consumer(s) (service consumer point of view)
// TODO: UI may be broken now because of userID and dependencyID (lint forced changing Id -> ID)
type ConsumerView struct {
	userID       string
	dependencyID string
	//	revision     *resolve.Revision
	g *graph
}

// NewConsumerView creates a new ConsumerView
func NewConsumerView(userID string, dependencyID string) ConsumerView {
	return ConsumerView{
		userID:       userID,
		dependencyID: dependencyID,
		//		revision:     resolve.LoadRevision(),
		g: newGraph(),
	}
}

// GetData returns graph for a given view
func (view ConsumerView) GetData() interface{} {
	/*
		// go over all dependencies of a given user
		for _, dependency := range view.revision.Policy.Dependencies.DependenciesByID {
			if filterMatches(dependency.UserID, view.userID) && filterMatches(dependency.GetID(), view.dependencyID) {
				// Step 1 - add a node for every matching dependency found
				dependencyNode := newDependencyNode(dependency, false, view.revision.UserLoader)
				view.g.addNode(dependencyNode, 0)

				// Step 2 - process subgraph (doesn't matter whether it's resolved successfully or not)
				view.addResolvedDependencies(dependency.ServiceKey, dependencyNode, 1)
			}
		}
	*/

	return view.g.GetData()
}

// Returns if value is good with respect to filterValue
func filterMatches(value string, filterValue string) bool {
	return len(filterValue) <= 0 || filterValue == value
}

// Adds to the graph nodes/edges which are triggered by usage of a given dependency
func (view ConsumerView) addResolvedDependencies(key string, nodePrev graphNode, nextLevel int) {
	/*
		// try to get this component instance from resolved data
		v := view.revision.Resolution.ComponentInstanceMap[key]

		// okay, this component likely failed to resolved, so let's look it up from unresolved pool
		if v == nil {
			v = view.revision.Resolution.Unresolved.ComponentInstanceMap[key]
		}

		// if it's a service, add node and connect with previous
		if v.Key.IsService() {
			// add service instance node
			svcInstanceNode := newServiceInstanceNode(key, view.revision.Policy.Services[v.Key.ServiceName], v.Key.ContextName, v.Key.ContextNameWithKeys, v, nextLevel <= 1)
			view.g.addNode(svcInstanceNode, nextLevel)

			// connect service instance nodes
			view.g.addEdge(nodePrev, svcInstanceNode)

			// update prev
			nodePrev = svcInstanceNode
		}

		// go over all outgoing edges
		for k := range v.EdgesOut {
			// proceed further with updated service instance node
			view.addResolvedDependencies(k, nodePrev, nextLevel+1)
		}
	*/
}
