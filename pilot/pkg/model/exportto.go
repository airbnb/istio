// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	typev1beta1 "istio.io/api/type/v1beta1"
	"istio.io/istio/pkg/config/visibility"
	"istio.io/istio/pkg/util/sets"
)

// ExportToTarget represents the visibility configuration for a service or config object.
// It can contain both static namespace names (including visibility constants like ".", "*", "~")
// and label selectors for dynamic namespace matching.
type ExportToTarget struct {
	// StaticNamespaces contains explicit namespace names and visibility constants (., *, ~)
	StaticNamespaces sets.Set[visibility.Instance]

	// Selectors contains label selectors for dynamic namespace matching
	Selectors []labels.Selector
}

// IsEmpty returns true if the ExportToTarget has no static namespaces and no selectors
func (e *ExportToTarget) IsEmpty() bool {
	if e == nil {
		return true
	}
	return len(e.StaticNamespaces) == 0 && len(e.Selectors) == 0
}

// IsPublic returns true if the ExportToTarget contains the public visibility constant "*"
func (e *ExportToTarget) IsPublic() bool {
	if e == nil {
		return false
	}
	return e.StaticNamespaces.Contains(visibility.Public)
}

// IsPrivate returns true if the ExportToTarget contains the private visibility constant "."
func (e *ExportToTarget) IsPrivate() bool {
	if e == nil {
		return false
	}
	return e.StaticNamespaces.Contains(visibility.Private)
}

// IsNone returns true if the ExportToTarget contains the none visibility constant "~"
func (e *ExportToTarget) IsNone() bool {
	if e == nil {
		return false
	}
	return e.StaticNamespaces.Contains(visibility.None)
}

// Matches evaluates if a namespace matches the export target based on:
// 1. Static namespace names (including visibility constants)
// 2. Label selectors against the namespace labels
//
// Parameters:
//   - namespace: the namespace name to check
//   - namespaceLabels: the labels of the namespace
//
// Returns true if the namespace matches any of the export targets
func (e *ExportToTarget) Matches(namespace string, namespaceLabels map[string]string) bool {
	if e == nil {
		return false
	}

	// Check if public visibility is enabled
	if e.StaticNamespaces.Contains(visibility.Public) {
		return true
	}

	// Check if the specific namespace is in the static list
	if e.StaticNamespaces.Contains(visibility.Instance(namespace)) {
		return true
	}

	// Check label selectors
	if len(e.Selectors) > 0 && namespaceLabels != nil {
		labelSet := labels.Set(namespaceLabels)
		for _, selector := range e.Selectors {
			if selector.Matches(labelSet) {
				return true
			}
		}
	}

	return false
}

// ParseExportTo converts a list of exportTo strings and label selectors to an ExportToTarget.
// This function is used to convert API types to internal types.
//
// Parameters:
//   - exportToList: list of static namespace names or visibility constants
//   - exportToSelectors: list of label selectors from the API
//
// Returns an ExportToTarget with both static namespaces and selectors populated
func ParseExportTo(exportToList []string, exportToSelectors []*typev1beta1.LabelSelector) (*ExportToTarget, error) {
	target := &ExportToTarget{
		StaticNamespaces: sets.New[visibility.Instance](),
		Selectors:        []labels.Selector{},
	}

	// Parse static namespace names and visibility constants
	for _, ns := range exportToList {
		target.StaticNamespaces.Insert(visibility.Instance(ns))
	}

	// Parse label selectors
	for _, ls := range exportToSelectors {
		selector, err := LabelSelectorAsSelector(ls)
		if err != nil {
			return nil, fmt.Errorf("invalid label selector: %v", err)
		}
		target.Selectors = append(target.Selectors, selector)
	}

	return target, nil
}

// Contains checks if the given visibility instance is in the static namespaces set.
// This provides backward compatibility with code that used sets.Set[visibility.Instance].Contains().
func (e *ExportToTarget) Contains(v visibility.Instance) bool {
	if e == nil {
		return false
	}
	return e.StaticNamespaces.Contains(v)
}

// Copy creates a deep copy of ExportToTarget.
// Note: labels.Selector is typically immutable once created, so we can share the selector references.
func (e *ExportToTarget) Copy() *ExportToTarget {
	if e == nil {
		return nil
	}
	out := &ExportToTarget{}

	// Preserve nil for StaticNamespaces
	if e.StaticNamespaces != nil {
		out.StaticNamespaces = e.StaticNamespaces.Copy()
	}

	// Preserve nil for Selectors
	if e.Selectors != nil {
		out.Selectors = make([]labels.Selector, len(e.Selectors))
		copy(out.Selectors, e.Selectors)
	}

	return out
}

// Equals checks if two ExportToTarget objects are equal.
// Note: Label selectors are compared by their string representation.
func (e *ExportToTarget) Equals(other *ExportToTarget) bool {
	if e == nil && other == nil {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	// Compare static namespaces
	if !e.StaticNamespaces.Equals(other.StaticNamespaces) {
		return false
	}
	// Compare selectors by count and string representation
	if len(e.Selectors) != len(other.Selectors) {
		return false
	}
	// Compare each selector's string representation
	// Note: Order matters for this comparison
	for i, sel := range e.Selectors {
		if sel.String() != other.Selectors[i].String() {
			return false
		}
	}
	return true
}

// StaticNamespacesList returns the static namespaces as a list for iteration.
// This is useful for code that needs to iterate over the namespaces.
func (e *ExportToTarget) StaticNamespacesList() []visibility.Instance {
	if e == nil {
		return nil
	}
	return e.StaticNamespaces.UnsortedList()
}

// HasSelectors returns true if there are any label selectors defined.
func (e *ExportToTarget) HasSelectors() bool {
	return e != nil && len(e.Selectors) > 0
}

// IsSuperset checks if this ExportToTarget is a superset of another.
// For static namespaces, checks set containment.
// If either has selectors, this is conservatively false (cannot determine superset relationship with dynamic selectors).
func (e *ExportToTarget) IsSuperset(other *ExportToTarget) bool {
	if e == nil || other == nil {
		return false
	}
	// If either has selectors, we cannot determine superset relationship
	if e.HasSelectors() || other.HasSelectors() {
		return false
	}
	// Check if all namespaces in 'other' are in 'e'
	return e.StaticNamespaces.SupersetOf(other.StaticNamespaces)
}

// Len returns the number of static namespaces (for backward compatibility).
func (e *ExportToTarget) Len() int {
	if e == nil {
		return 0
	}
	return len(e.StaticNamespaces)
}

// LabelSelectorAsSelector converts a type v1beta1 LabelSelector to a labels.Selector.
// This follows the same pattern used in the ambient controller.
func LabelSelectorAsSelector(ps *typev1beta1.LabelSelector) (labels.Selector, error) {
	if ps == nil {
		return labels.Nothing(), nil
	}
	if len(ps.MatchLabels)+len(ps.MatchExpressions) == 0 {
		return labels.Everything(), nil
	}
	requirements := make([]labels.Requirement, 0, len(ps.MatchLabels)+len(ps.MatchExpressions))
	for k, v := range ps.MatchLabels {
		r, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	for _, expr := range ps.MatchExpressions {
		var op selection.Operator
		switch metav1.LabelSelectorOperator(expr.Operator) {
		case metav1.LabelSelectorOpIn:
			op = selection.In
		case metav1.LabelSelectorOpNotIn:
			op = selection.NotIn
		case metav1.LabelSelectorOpExists:
			op = selection.Exists
		case metav1.LabelSelectorOpDoesNotExist:
			op = selection.DoesNotExist
		default:
			return nil, fmt.Errorf("%q is not a valid label selector operator", expr.Operator)
		}
		r, err := labels.NewRequirement(expr.Key, op, append([]string(nil), expr.Values...))
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *r)
	}
	selector := labels.NewSelector()
	selector = selector.Add(requirements...)
	return selector, nil
}
