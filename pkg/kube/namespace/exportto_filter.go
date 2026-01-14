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

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/slices"
)

// ExportToFilter tracks namespace membership for exportTo selectors.
// It watches namespace changes and caches namespace labels to support
// efficient selector evaluation for exportTo visibility.
type ExportToFilter struct {
	mu sync.RWMutex

	// namespaceLabels caches labels for all namespaces
	namespaceLabels map[string]map[string]string

	// handlers to notify when namespace labels change
	handlers []func(namespaceName string)

	namespaces kclient.Client[*corev1.Namespace]
}

// NewExportToFilter creates a new ExportToFilter that watches namespace changes
// and caches namespace labels for efficient exportTo selector evaluation.
func NewExportToFilter(namespaces kclient.Client[*corev1.Namespace], stop <-chan struct{}) *ExportToFilter {
	f := &ExportToFilter{
		namespaceLabels: make(map[string]map[string]string),
		namespaces:      namespaces,
	}

	// Register event handlers for namespace changes
	namespaces.AddEventHandler(controllers.EventHandler[*corev1.Namespace]{
		AddFunc: func(ns *corev1.Namespace) {
			f.mu.Lock()
			f.namespaceCreatedLocked(ns)
			f.mu.Unlock()
			f.notifyHandlers(ns.Name)
		},
		UpdateFunc: func(oldObj, newObj *corev1.Namespace) {
			f.mu.Lock()
			labelsChanged := f.namespaceUpdatedLocked(oldObj, newObj)
			f.mu.Unlock()
			if labelsChanged {
				f.notifyHandlers(newObj.Name)
			}
		},
		DeleteFunc: func(ns *corev1.Namespace) {
			f.mu.Lock()
			f.namespaceDeletedLocked(ns)
			f.mu.Unlock()
			f.notifyHandlers(ns.Name)
		},
	})

	// Start namespaces and wait for it to be ready
	namespaces.Start(stop)
	kube.WaitForCacheSync("exportto filter", stop, namespaces.HasSynced)

	// Initialize the cache with existing namespaces
	f.initializeCache()

	return f
}

// initializeCache populates the namespace labels cache with all existing namespaces
func (f *ExportToFilter) initializeCache() {
	f.mu.Lock()
	defer f.mu.Unlock()

	namespaceList := f.namespaces.List("", labels.Everything())
	for _, ns := range namespaceList {
		f.namespaceLabels[ns.Name] = copyLabels(ns.Labels)
	}
}

// GetNamespaceLabels returns cached labels for a namespace.
// Returns nil if the namespace is not found in the cache.
func (f *ExportToFilter) GetNamespaceLabels(namespace string) map[string]string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	labels, exists := f.namespaceLabels[namespace]
	if !exists {
		return nil
	}
	// Return a copy to prevent external modifications
	return copyLabels(labels)
}

// MatchesSelector checks if a namespace matches a given selector.
// Returns false if the namespace is not found in the cache.
func (f *ExportToFilter) MatchesSelector(namespace string, selector labels.Selector) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	nsLabels, exists := f.namespaceLabels[namespace]
	if !exists {
		return false
	}

	return selector.Matches(labels.Set(nsLabels))
}

// AddHandler registers a callback for namespace label changes.
// The handler will be called whenever a namespace is created, updated (if labels changed), or deleted.
func (f *ExportToFilter) AddHandler(handler func(namespaceName string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers = append(f.handlers, handler)
}

// notifyHandlers calls all registered handlers with the namespace name that changed.
// Handlers are called without holding the lock to prevent deadlocks.
func (f *ExportToFilter) notifyHandlers(namespaceName string) {
	// Clone handlers; we handle dynamic handlers so they can change after the filter has started.
	// Important: handlers are not called under the lock to prevent deadlocks.
	f.mu.RLock()
	handlers := slices.Clone(f.handlers)
	f.mu.RUnlock()

	for _, h := range handlers {
		h(namespaceName)
	}
}

// namespaceCreatedLocked handles a newly created namespace by caching its labels.
// Must be called with lock held.
func (f *ExportToFilter) namespaceCreatedLocked(ns *corev1.Namespace) {
	f.namespaceLabels[ns.Name] = copyLabels(ns.Labels)
}

// namespaceUpdatedLocked handles a namespace update by updating cached labels.
// Returns true if the labels changed.
// Must be called with lock held.
func (f *ExportToFilter) namespaceUpdatedLocked(oldNs, newNs *corev1.Namespace) bool {
	// Check if labels actually changed
	if labelsEqual(oldNs.Labels, newNs.Labels) {
		return false
	}

	f.namespaceLabels[newNs.Name] = copyLabels(newNs.Labels)
	return true
}

// namespaceDeletedLocked handles a namespace deletion by removing it from the cache.
// Must be called with lock held.
func (f *ExportToFilter) namespaceDeletedLocked(ns *corev1.Namespace) {
	delete(f.namespaceLabels, ns.Name)
}

// copyLabels creates a deep copy of a label map
func copyLabels(labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}
	result := make(map[string]string, len(labels))
	for k, v := range labels {
		result[k] = v
	}
	return result
}

// labelsEqual compares two label maps for equality
func labelsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
