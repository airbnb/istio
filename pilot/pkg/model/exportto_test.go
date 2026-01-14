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

package model_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/labels"

	typev1beta1 "istio.io/api/type/v1beta1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config/visibility"
	"istio.io/istio/pkg/test/util/assert"
	"istio.io/istio/pkg/util/sets"
)

func TestExportToTarget_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		target   *model.ExportToTarget
		expected bool
	}{
		{
			name:     "nil target",
			target:   nil,
			expected: true,
		},
		{
			name: "empty target",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{},
			},
			expected: true,
		},
		{
			name: "target with static namespace",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](visibility.Public),
				Selectors:        []labels.Selector{},
			},
			expected: false,
		},
		{
			name: "target with selector",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{labels.Everything()},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.IsEmpty()
			if got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExportToTarget_Matches(t *testing.T) {
	// Create label selectors for testing
	envProdSelector, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "prod"},
	})
	envTestSelector, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "test"},
	})
	matchExpressionSelector, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
			{
				Key:      "tier",
				Operator: "In",
				Values:   []string{"frontend", "backend"},
			},
		},
	})

	tests := []struct {
		name            string
		target          *model.ExportToTarget
		namespace       string
		namespaceLabels map[string]string
		expected        bool
	}{
		{
			name:            "nil target",
			target:          nil,
			namespace:       "ns1",
			namespaceLabels: map[string]string{},
			expected:        false,
		},
		{
			name: "public visibility",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](visibility.Public),
				Selectors:        []labels.Selector{},
			},
			namespace:       "any-namespace",
			namespaceLabels: map[string]string{},
			expected:        true,
		},
		{
			name: "specific static namespace match",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			namespace:       "ns1",
			namespaceLabels: map[string]string{},
			expected:        true,
		},
		{
			name: "specific static namespace no match",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			namespace:       "ns3",
			namespaceLabels: map[string]string{},
			expected:        false,
		},
		{
			name: "label selector match",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace: "prod-ns",
			namespaceLabels: map[string]string{
				"env": "prod",
			},
			expected: true,
		},
		{
			name: "label selector no match",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace: "test-ns",
			namespaceLabels: map[string]string{
				"env": "test",
			},
			expected: false,
		},
		{
			name: "multiple selectors - one matches",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{envProdSelector, envTestSelector},
			},
			namespace: "test-ns",
			namespaceLabels: map[string]string{
				"env": "test",
			},
			expected: true,
		},
		{
			name: "combined static and selector - static matches",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("special-ns"),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace: "special-ns",
			namespaceLabels: map[string]string{
				"env": "dev",
			},
			expected: true,
		},
		{
			name: "combined static and selector - selector matches",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("special-ns"),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace: "prod-ns",
			namespaceLabels: map[string]string{
				"env": "prod",
			},
			expected: true,
		},
		{
			name: "combined static and selector - neither matches",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("special-ns"),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace: "other-ns",
			namespaceLabels: map[string]string{
				"env": "dev",
			},
			expected: false,
		},
		{
			name: "match expression selector - matches",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{matchExpressionSelector},
			},
			namespace: "app-ns",
			namespaceLabels: map[string]string{
				"tier": "frontend",
			},
			expected: true,
		},
		{
			name: "match expression selector - no match",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{matchExpressionSelector},
			},
			namespace: "db-ns",
			namespaceLabels: map[string]string{
				"tier": "database",
			},
			expected: false,
		},
		{
			name: "selector with nil labels",
			target: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{envProdSelector},
			},
			namespace:       "ns1",
			namespaceLabels: nil,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.Matches(tt.namespace, tt.namespaceLabels)
			if got != tt.expected {
				t.Errorf("Matches() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseExportTo(t *testing.T) {
	tests := []struct {
		name               string
		exportToList       []string
		exportToSelectors  []*typev1beta1.LabelSelector
		expectError        bool
		validateResult     func(*testing.T, *model.ExportToTarget)
	}{
		{
			name:               "empty inputs",
			exportToList:       []string{},
			exportToSelectors:  []*typev1beta1.LabelSelector{},
			expectError:        false,
			validateResult: func(t *testing.T, target *model.ExportToTarget) {
				assert.Equal(t, true, target.IsEmpty())
			},
		},
		{
			name:               "static namespaces only",
			exportToList:       []string{"ns1", "ns2", "*"},
			exportToSelectors:  []*typev1beta1.LabelSelector{},
			expectError:        false,
			validateResult: func(t *testing.T, target *model.ExportToTarget) {
				assert.Equal(t, 3, target.Len())
				assert.Equal(t, true, target.Contains(visibility.Instance("ns1")))
				assert.Equal(t, true, target.Contains(visibility.Instance("ns2")))
				assert.Equal(t, true, target.Contains(visibility.Public))
			},
		},
		{
			name:         "selectors only",
			exportToList: []string{},
			exportToSelectors: []*typev1beta1.LabelSelector{
				{
					MatchLabels: map[string]string{"env": "prod"},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, target *model.ExportToTarget) {
				assert.Equal(t, true, target.HasSelectors())
				assert.Equal(t, 1, len(target.Selectors))
			},
		},
		{
			name:         "combined static and selectors",
			exportToList: []string{"ns1", "."},
			exportToSelectors: []*typev1beta1.LabelSelector{
				{
					MatchLabels: map[string]string{"env": "prod"},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, target *model.ExportToTarget) {
				assert.Equal(t, 2, target.Len())
				assert.Equal(t, true, target.HasSelectors())
				assert.Equal(t, true, target.Contains(visibility.Private))
				assert.Equal(t, true, target.Contains(visibility.Instance("ns1")))
			},
		},
		{
			name:         "match expressions",
			exportToList: []string{},
			exportToSelectors: []*typev1beta1.LabelSelector{
				{
					MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
						{
							Key:      "tier",
							Operator: "In",
							Values:   []string{"frontend", "backend"},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, target *model.ExportToTarget) {
				assert.Equal(t, true, target.HasSelectors())
				assert.Equal(t, 1, len(target.Selectors))
			},
		},
		{
			name:         "invalid selector operator",
			exportToList: []string{},
			exportToSelectors: []*typev1beta1.LabelSelector{
				{
					MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
						{
							Key:      "tier",
							Operator: "InvalidOp",
							Values:   []string{"value"},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := model.ParseExportTo(tt.exportToList, tt.exportToSelectors)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseExportTo() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseExportTo() unexpected error: %v", err)
				return
			}
			if tt.validateResult != nil {
				tt.validateResult(t, result)
			}
		})
	}
}

func TestExportToTarget_Copy(t *testing.T) {
	selector, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "prod"},
	})

	original := &model.ExportToTarget{
		StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
		Selectors:        []labels.Selector{selector},
	}

	// Test copy
	copied := original.Copy()

	// Verify deep copy of static namespaces
	if !original.StaticNamespaces.Equals(copied.StaticNamespaces) {
		t.Errorf("Copy() static namespaces not equal")
	}

	// Modify original and verify copied is not affected
	original.StaticNamespaces.Insert("ns3")
	if copied.StaticNamespaces.Contains("ns3") {
		t.Errorf("Copy() is not a deep copy - modifications affect the copy")
	}

	// Test nil copy
	var nilTarget *model.ExportToTarget
	if nilTarget.Copy() != nil {
		t.Errorf("Copy() of nil should return nil")
	}
}

func TestExportToTarget_Equals(t *testing.T) {
	selector1, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "prod"},
	})
	selector2, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "test"},
	})

	tests := []struct {
		name     string
		target1  *model.ExportToTarget
		target2  *model.ExportToTarget
		expected bool
	}{
		{
			name:     "both nil",
			target1:  nil,
			target2:  nil,
			expected: true,
		},
		{
			name:     "one nil",
			target1:  nil,
			target2:  &model.ExportToTarget{StaticNamespaces: sets.New[visibility.Instance]()},
			expected: false,
		},
		{
			name: "same static namespaces, no selectors",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			expected: true,
		},
		{
			name: "different static namespaces",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1"),
				Selectors:        []labels.Selector{},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns2"),
				Selectors:        []labels.Selector{},
			},
			expected: false,
		},
		{
			name: "same selectors",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector1},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector1},
			},
			expected: true,
		},
		{
			name: "different selectors",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector1},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector2},
			},
			expected: false,
		},
		{
			name: "different selector count",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector1},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance](),
				Selectors:        []labels.Selector{selector1, selector2},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target1.Equals(tt.target2)
			if got != tt.expected {
				t.Errorf("Equals() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExportToTarget_IsSuperset(t *testing.T) {
	selector, _ := model.LabelSelectorAsSelector(&typev1beta1.LabelSelector{
		MatchLabels: map[string]string{"env": "prod"},
	})

	tests := []struct {
		name     string
		target1  *model.ExportToTarget
		target2  *model.ExportToTarget
		expected bool
	}{
		{
			name: "superset static namespaces",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2", "ns3"),
				Selectors:        []labels.Selector{},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			expected: true,
		},
		{
			name: "not a superset",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1"),
				Selectors:        []labels.Selector{},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{},
			},
			expected: false,
		},
		{
			name: "with selectors - returns false",
			target1: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1", "ns2"),
				Selectors:        []labels.Selector{selector},
			},
			target2: &model.ExportToTarget{
				StaticNamespaces: sets.New[visibility.Instance]("ns1"),
				Selectors:        []labels.Selector{},
			},
			expected: false,
		},
		{
			name:     "nil targets",
			target1:  nil,
			target2:  nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target1.IsSuperset(tt.target2)
			if got != tt.expected {
				t.Errorf("IsSuperset() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLabelSelectorAsSelector(t *testing.T) {
	tests := []struct {
		name         string
		selector     *typev1beta1.LabelSelector
		expectError  bool
		testMatch    func(*testing.T, labels.Selector)
	}{
		{
			name:        "nil selector",
			selector:    nil,
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				// Nothing() selector matches nothing
				if sel.Matches(labels.Set{"any": "label"}) {
					t.Error("nil selector should match nothing")
				}
			},
		},
		{
			name:        "empty selector",
			selector:    &typev1beta1.LabelSelector{},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				// Everything() selector matches everything
				if !sel.Matches(labels.Set{"any": "label"}) {
					t.Error("empty selector should match everything")
				}
			},
		},
		{
			name: "match labels",
			selector: &typev1beta1.LabelSelector{
				MatchLabels: map[string]string{
					"env":  "prod",
					"tier": "frontend",
				},
			},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				// Should match when all labels present
				if !sel.Matches(labels.Set{"env": "prod", "tier": "frontend", "extra": "label"}) {
					t.Error("should match when all required labels present")
				}
				// Should not match when missing a label
				if sel.Matches(labels.Set{"env": "prod"}) {
					t.Error("should not match when missing required labels")
				}
			},
		},
		{
			name: "match expressions - In",
			selector: &typev1beta1.LabelSelector{
				MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
					{
						Key:      "env",
						Operator: "In",
						Values:   []string{"prod", "staging"},
					},
				},
			},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				if !sel.Matches(labels.Set{"env": "prod"}) {
					t.Error("should match when value is in list")
				}
				if sel.Matches(labels.Set{"env": "dev"}) {
					t.Error("should not match when value not in list")
				}
			},
		},
		{
			name: "match expressions - NotIn",
			selector: &typev1beta1.LabelSelector{
				MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
					{
						Key:      "env",
						Operator: "NotIn",
						Values:   []string{"dev"},
					},
				},
			},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				if !sel.Matches(labels.Set{"env": "prod"}) {
					t.Error("should match when value not in excluded list")
				}
				if sel.Matches(labels.Set{"env": "dev"}) {
					t.Error("should not match when value in excluded list")
				}
			},
		},
		{
			name: "match expressions - Exists",
			selector: &typev1beta1.LabelSelector{
				MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
					{
						Key:      "env",
						Operator: "Exists",
					},
				},
			},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				if !sel.Matches(labels.Set{"env": "any-value"}) {
					t.Error("should match when label exists")
				}
				if sel.Matches(labels.Set{"other": "label"}) {
					t.Error("should not match when label doesn't exist")
				}
			},
		},
		{
			name: "match expressions - DoesNotExist",
			selector: &typev1beta1.LabelSelector{
				MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
					{
						Key:      "env",
						Operator: "DoesNotExist",
					},
				},
			},
			expectError: false,
			testMatch: func(t *testing.T, sel labels.Selector) {
				if !sel.Matches(labels.Set{"other": "label"}) {
					t.Error("should match when label doesn't exist")
				}
				if sel.Matches(labels.Set{"env": "prod"}) {
					t.Error("should not match when label exists")
				}
			},
		},
		{
			name: "invalid operator",
			selector: &typev1beta1.LabelSelector{
				MatchExpressions: []*typev1beta1.LabelSelectorRequirement{
					{
						Key:      "env",
						Operator: "InvalidOp",
						Values:   []string{"value"},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, err := model.LabelSelectorAsSelector(tt.selector)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tt.testMatch != nil {
				tt.testMatch(t, selector)
			}
		})
	}
}
