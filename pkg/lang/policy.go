package lang

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/object"
	"gopkg.in/go-playground/validator.v9"
	"strings"
	"sync"
)

// Policy describes the entire Aptomi policy.
//
// At the highest level, policy consists of namespaces. Namespaces provide isolation for policy objects and access to
// namespaces can be controlled via ACL rules. Thus, different users can have different access rights to different parts
// of Aptomi policy. Namespaces are useful in environments with many users, multiple teams and projects.
//
// Objects get stored in their corresponding namespaces. Names of objects must be unique within a namespace and a given
// object kind.
//
// Once policy is defined, it can be passed to the engine for policy resolution. Policy resolution translates a given
// policy (intent) into actual state (what services/components need to created/updated/deleted, how and where) and the
// corresponding set of actions.
type Policy struct {
	// Namespace is a map from namespace name into a PolicyNamespace
	Namespace map[string]*PolicyNamespace

	// Validator for policy objects
	validator *validator.Validate

	// Access control rules for different policy namespaces
	once        sync.Once
	aclResolver *ACLResolver // lazily initialized value
}

// NewPolicy creates a new Policy
func NewPolicy() *Policy {
	return &Policy{
		Namespace: make(map[string]*PolicyNamespace),
		validator: makeValidator(),
	}
}

// View returns a policy view object, which allows to make all policy operations on behalf of a certain user
// Policy view object will enforce all ACLs, allowing the user to only perform actions which he is allowed to perform
// All ACL rules should be loaded and added to the policy before this method gets called
func (policy *Policy) View(user *User) *PolicyView {
	policy.once.Do(func() {
		systemNamespace := policy.Namespace[object.SystemNS]
		if systemNamespace != nil {
			policy.aclResolver = NewACLResolver(systemNamespace.ACLRules)
		} else {
			policy.aclResolver = NewACLResolver(NewGlobalRules())
		}
	})
	return NewPolicyView(policy, user)
}

// AddObject adds an object into the policy. When you add objects to the policy, they get added to the corresponding
// Namespace. If error occurs (e.g. object validation error, etc) then the error will be returned
func (policy *Policy) AddObject(obj object.Base) error {
	policyNamespace, ok := policy.Namespace[obj.GetNamespace()]
	if !ok {
		policyNamespace = NewPolicyNamespace(obj.GetNamespace(), policy.validator)
		policy.Namespace[obj.GetNamespace()] = policyNamespace
	}
	return policyNamespace.addObject(obj)
}

// GetObjectsByKind returns all objects in a policy with a given kind, across all namespaces
func (policy *Policy) GetObjectsByKind(kind string) []object.Base {
	result := []object.Base{}
	for _, policyNS := range policy.Namespace {
		result = append(result, policyNS.getObjectsByKind(kind)...)
	}
	return result
}

// GetObject looks up and returns an object from the policy, given its kind, locator ([namespace/]name), and current
// namespace relative to which the call is being made
func (policy *Policy) GetObject(kind string, locator string, currentNs string) (object.Base, error) {
	// parse locator: [namespace/]name. we might add [domain/] in the future
	parts := strings.Split(locator, "/")
	var ns, name string
	if len(parts) == 1 {
		ns = currentNs
		name = parts[0]
	} else if len(parts) == 2 {
		ns = parts[0]
		name = parts[1]
	} else {
		return nil, fmt.Errorf("can't parse policy object locator: '%s'", locator)
	}

	policyNS, ok := policy.Namespace[ns]
	if !ok {
		return nil, fmt.Errorf("namespace '%s' doesn't exist, but referenced in locator '%s'", ns, locator)
	}

	return policyNS.getObject(kind, name)
}
