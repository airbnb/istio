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
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/test"
	"istio.io/istio/pkg/test/util/retry"
)

func TestExportToFilter(t *testing.T) {
	stop := test.NewStop(t)
	client := kube.NewFakeClient()
	client.RunAndWait(stop)

	// Create the filter
	filter := NewExportToFilter(client.Kube().CoreV1().Namespaces(), stop)

	// Create test namespaces
	ns1 := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns1",
			Labels: map[string]string{
				"env": "prod",
				"app": "web",
			},
		},
	}

	ns2 := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns2",
			Labels: map[string]string{
				"env": "dev",
				"app": "api",
			},
		},
	}

	// Add namespace 1
	_, err := client.Kube().CoreV1().Namespaces().Create(nil, ns1, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond) // Give time for event to propagate

	// Test GetNamespaceLabels
	t.Run("GetNamespaceLabels", func(t *testing.T) {
		labels := filter.GetNamespaceLabels("ns1")
		if labels == nil {
			t.Fatal("expected labels for ns1, got nil")
		}
		if labels["env"] != "prod" {
			t.Errorf("expected env=prod, got %s", labels["env"])
		}
		if labels["app"] != "web" {
			t.Errorf("expected app=web, got %s", labels["app"])
		}

		// Non-existent namespace should return nil
		labels = filter.GetNamespaceLabels("nonexistent")
		if labels != nil {
			t.Errorf("expected nil for nonexistent namespace, got %v", labels)
		}
	})

	// Test MatchesSelector
	t.Run("MatchesSelector", func(t *testing.T) {
		// Create a selector that matches env=prod
		selector, err := labels.Parse("env=prod")
		if err != nil {
			t.Fatal(err)
		}

		matches := filter.MatchesSelector("ns1", selector)
		if !matches {
			t.Error("expected ns1 to match selector env=prod")
		}

		// Selector that doesn't match
		selector, err = labels.Parse("env=dev")
		if err != nil {
			t.Fatal(err)
		}

		matches = filter.MatchesSelector("ns1", selector)
		if matches {
			t.Error("expected ns1 to not match selector env=dev")
		}

		// Non-existent namespace should not match
		matches = filter.MatchesSelector("nonexistent", selector)
		if matches {
			t.Error("expected nonexistent namespace to not match")
		}
	})

	// Test namespace updates
	t.Run("NamespaceUpdate", func(t *testing.T) {
		handlerCalled := make(chan string, 10)

		filter.AddHandler(func(ns string) {
			handlerCalled <- ns
		})

		// Update namespace labels
		ns1Updated := ns1.DeepCopy()
		ns1Updated.Labels["env"] = "staging"

		_, err := client.Kube().CoreV1().Namespaces().Update(nil, ns1Updated, metav1.UpdateOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// Handler should have been called
		select {
		case ns := <-handlerCalled:
			if ns != "ns1" {
				t.Errorf("expected handler to be called with ns1, got %s", ns)
			}
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for handler to be called")
		}

		// Verify labels were updated
		retry.UntilOrFail(t, func() bool {
			labels := filter.GetNamespaceLabels("ns1")
			return labels != nil && labels["env"] == "staging"
		}, retry.Timeout(2*time.Second))
	})

	// Test namespace creation
	t.Run("NamespaceCreation", func(t *testing.T) {
		handlerCalled := make(chan string, 10)

		filter.AddHandler(func(ns string) {
			handlerCalled <- ns
		})

		// Create namespace 2
		_, err := client.Kube().CoreV1().Namespaces().Create(nil, ns2, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// Handler should have been called
		select {
		case ns := <-handlerCalled:
			if ns != "ns2" {
				t.Errorf("expected handler to be called with ns2, got %s", ns)
			}
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for handler to be called")
		}

		// Verify labels were cached
		retry.UntilOrFail(t, func() bool {
			labels := filter.GetNamespaceLabels("ns2")
			return labels != nil && labels["env"] == "dev" && labels["app"] == "api"
		}, retry.Timeout(2*time.Second))
	})

	// Test namespace deletion
	t.Run("NamespaceDeletion", func(t *testing.T) {
		handlerCalled := make(chan string, 10)

		filter.AddHandler(func(ns string) {
			handlerCalled <- ns
		})

		// Delete namespace 2
		err := client.Kube().CoreV1().Namespaces().Delete(nil, "ns2", metav1.DeleteOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// Handler should have been called
		select {
		case ns := <-handlerCalled:
			if ns != "ns2" {
				t.Errorf("expected handler to be called with ns2, got %s", ns)
			}
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for handler to be called")
		}

		// Verify labels were removed from cache
		retry.UntilOrFail(t, func() bool {
			labels := filter.GetNamespaceLabels("ns2")
			return labels == nil
		}, retry.Timeout(2*time.Second))
	})
}

func TestExportToFilterConcurrency(t *testing.T) {
	stop := test.NewStop(t)
	client := kube.NewFakeClient()
	client.RunAndWait(stop)

	filter := NewExportToFilter(client.Kube().CoreV1().Namespaces(), stop)

	// Create a namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ns",
			Labels: map[string]string{
				"env": "test",
			},
		},
	}
	_, err := client.Kube().CoreV1().Namespaces().Create(nil, ns, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test concurrent reads and writes
	t.Run("ConcurrentAccess", func(t *testing.T) {
		done := make(chan struct{})

		// Concurrent reads
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- struct{}{} }()
				for j := 0; j < 100; j++ {
					filter.GetNamespaceLabels("test-ns")
					selector, _ := labels.Parse("env=test")
					filter.MatchesSelector("test-ns", selector)
				}
			}()
		}

		// Concurrent handler registration
		for i := 0; i < 5; i++ {
			go func() {
				defer func() { done <- struct{}{} }()
				filter.AddHandler(func(ns string) {
					// Empty handler
				})
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 15; i++ {
			<-done
		}
	})
}

func TestLabelsCopyAndEquality(t *testing.T) {
	t.Run("copyLabels", func(t *testing.T) {
		original := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		copy := copyLabels(original)
		if copy["key1"] != "value1" {
			t.Errorf("expected key1=value1, got %s", copy["key1"])
		}
		if copy["key2"] != "value2" {
			t.Errorf("expected key2=value2, got %s", copy["key2"])
		}

		// Modify copy should not affect original
		copy["key3"] = "value3"
		_, exists := original["key3"]
		if exists {
			t.Error("modifying copy should not affect original")
		}

		// Nil input should return nil
		nilCopy := copyLabels(nil)
		if nilCopy != nil {
			t.Errorf("expected nil, got %v", nilCopy)
		}
	})

	t.Run("labelsEqual", func(t *testing.T) {
		a := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		b := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		if !labelsEqual(a, b) {
			t.Error("expected labels to be equal")
		}

		// Different values
		b["key2"] = "different"
		if labelsEqual(a, b) {
			t.Error("expected labels to be different")
		}

		// Different keys
		c := map[string]string{
			"key1": "value1",
		}
		if labelsEqual(a, c) {
			t.Error("expected labels to be different (different keys)")
		}

		// Empty maps
		d := map[string]string{}
		e := map[string]string{}
		if !labelsEqual(d, e) {
			t.Error("expected empty labels to be equal")
		}
	})
}
