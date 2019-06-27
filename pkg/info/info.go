package info

import (
	"fmt"
	"strings"
)

//PolicyInfo defines policy information
type PolicyInfo struct {
	// Name is policy name
	Name string
	// Resource is resource name
	Resource string
	// Kind represents the resource kind
	Kind string
	// Namespace is the ns of resource
	// empty on non-namespaced resources
	Namespace string
	success   bool
	Rules     []*RuleInfo
}

//NewPolicyInfo returns a new policy info
func NewPolicyInfo(policyName, resource, kind, ns string) *PolicyInfo {
	return &PolicyInfo{
		Name:      policyName,
		Resource:  resource,
		Kind:      kind,
		Namespace: ns,
		success:   true, // fail to be set explicity
	}
}

//IsSuccessful checks if policy is succesful
// the policy is set to fail, if any of the rules have failed
func (pi *PolicyInfo) IsSuccessful() bool {
	return pi.success
}

//ErrorRules returns error msgs from all rule
func (pi *PolicyInfo) ErrorRules() string {
	errorMsgs := []string{}
	for _, r := range pi.Rules {
		if !r.IsSuccessful() {
			errorMsgs = append(errorMsgs, r.ToString())
		}
	}
	return strings.Join(errorMsgs, ";")
}

type RuleType int

const (
	Mutation RuleType = iota
	Validation
	Generation
)

//RuleInfo defines rule struct
type RuleInfo struct {
	Name     string
	Msgs     []string
	ruleType RuleType
	success  bool
}

//ToString reule information
func (ri *RuleInfo) ToString() string {
	str := "rulename: " + ri.Name
	msgs := strings.Join(ri.Msgs, ";")
	return strings.Join([]string{str, msgs}, ";")
}

//NewRuleInfo creates a new RuleInfo
func NewRuleInfo(ruleName string, ruleType RuleType) *RuleInfo {
	return &RuleInfo{
		Name:     ruleName,
		Msgs:     []string{},
		ruleType: ruleType,
		success:  true, // fail to be set explicity
	}
}

//Fail set the rule as failed
func (ri *RuleInfo) Fail() {
	ri.success = false
}

//IsSuccessful checks if rule is succesful
func (ri *RuleInfo) IsSuccessful() bool {
	return ri.success
}

//Add add msg
func (ri *RuleInfo) Add(msg string) {
	ri.Msgs = append(ri.Msgs, msg)
}

//Addf add msg with args
func (ri *RuleInfo) Addf(msg string, args ...interface{}) {
	ri.Msgs = append(ri.Msgs, fmt.Sprintf(msg, args...))
}

//RulesSuccesfuly check if the any rule has failed or not
func RulesSuccesfuly(rules []*RuleInfo) bool {
	for _, r := range rules {
		if !r.success {
			return false
		}
	}
	return true
}

//AddRuleInfos sets the rule information
func (pi *PolicyInfo) AddRuleInfos(rules []*RuleInfo) {
	if !RulesSuccesfuly(rules) {
		pi.success = false
	}
	pi.Rules = rules
}
