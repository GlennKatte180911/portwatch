// Package portpolicy enforces allow/deny rules against observed ports.
// Rules are evaluated in order; the first matching rule wins.
package portpolicy

import (
	"fmt"
	"sync"
)

// Action is the outcome of a policy evaluation.
type Action string

const (
	ActionAllow Action = "allow"
	ActionDeny  Action = "deny"
)

// Rule pairs a port predicate with an action.
type Rule struct {
	Ports  []int
	Action Action
	Reason string
}

// matches returns true when port is listed in the rule.
func (r Rule) matches(port int) bool {
	for _, p := range r.Ports {
		if p == port {
			return true
		}
	}
	return false
}

// Result holds the outcome for a single port evaluation.
type Result struct {
	Port   int
	Action Action
	Reason string
}

// Policy holds an ordered set of rules.
type Policy struct {
	mu      sync.RWMutex
	rules   []Rule
	default_ Action
}

// New creates a Policy with the given default action applied when no rule matches.
func New(defaultAction Action) (*Policy, error) {
	if defaultAction != ActionAllow && defaultAction != ActionDeny {
		return nil, fmt.Errorf("portpolicy: invalid default action %q", defaultAction)
	}
	return &Policy{default_: defaultAction}, nil
}

// AddRule appends a rule to the policy.
func (p *Policy) AddRule(r Rule) error {
	if r.Action != ActionAllow && r.Action != ActionDeny {
		return fmt.Errorf("portpolicy: invalid action %q", r.Action)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, r)
	return nil
}

// Evaluate returns the Result for a single port.
func (p *Policy) Evaluate(port int) Result {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if r.matches(port) {
			return Result{Port: port, Action: r.Action, Reason: r.Reason}
		}
	}
	return Result{Port: port, Action: p.default_, Reason: "default"}
}

// Apply evaluates every port and returns only those that are denied.
func (p *Policy) Apply(ports []int) []Result {
	var denied []Result
	for _, port := range ports {
		if r := p.Evaluate(port); r.Action == ActionDeny {
			denied = append(denied, r)
		}
	}
	return denied
}
