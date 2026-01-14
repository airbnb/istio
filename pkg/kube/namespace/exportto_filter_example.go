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

package namespace

// Example usage of ExportToFilter for exportTo label selector support.
//
// This example demonstrates how to use ExportToFilter to track namespace
// label changes and evaluate exportTo selectors for resource visibility.
//
// Basic usage:
//
//	// Create the filter
//	filter := NewExportToFilter(namespaceClient, stop)
//
//	// Register a handler to be notified of namespace changes
//	filter.AddHandler(func(namespaceName string) {
//	    log.Infof("Namespace %s labels changed", namespaceName)
//	    // Trigger re-evaluation of exportTo visibility for resources
//	    // that reference this namespace
//	})
//
//	// Check if a namespace matches an exportTo selector
//	selector, _ := labels.Parse("env=prod,team=platform")
//	if filter.MatchesSelector("my-namespace", selector) {
//	    // Resource is visible in my-namespace
//	}
//
//	// Get namespace labels for custom evaluation
//	labels := filter.GetNamespaceLabels("my-namespace")
//	if labels["env"] == "prod" {
//	    // Do something based on label value
//	}
//
// Integration with exportTo:
//
// When a VirtualService, DestinationRule, or ServiceEntry has an exportToSelector,
// the ExportToFilter can be used to determine which namespaces match the selector:
//
//	// Example: VirtualService with exportToSelectors
//	vs := &v1alpha3.VirtualService{
//	    Spec: v1alpha3.VirtualServiceSpec{
//	        ExportToSelectors: []*v1alpha1.LabelSelector{
//	            {
//	                MatchLabels: map[string]string{
//	                    "env": "prod",
//	                },
//	            },
//	        },
//	    },
//	}
//
//	// Convert to k8s label selector
//	selector, _ := LabelSelectorAsSelector(vs.Spec.ExportToSelectors[0])
//
//	// Check which namespaces match
//	for _, ns := range allNamespaces {
//	    if filter.MatchesSelector(ns, selector) {
//	        // VirtualService should be visible in namespace ns
//	    }
//	}
//
// Handling namespace changes:
//
// The filter automatically watches for namespace create/update/delete events
// and triggers registered handlers when namespace labels change. This allows
// the control plane to reactively update resource visibility:
//
//	filter.AddHandler(func(namespaceName string) {
//	    // Re-evaluate all exportTo selectors that might match this namespace
//	    for _, resource := range resourcesWithExportToSelectors {
//	        for _, selector := range resource.ExportToSelectors {
//	            sel, _ := LabelSelectorAsSelector(selector)
//	            if filter.MatchesSelector(namespaceName, sel) {
//	                // Push config update to proxies in namespaceName
//	                pushConfigUpdate(namespaceName, resource)
//	            }
//	        }
//	    }
//	})
