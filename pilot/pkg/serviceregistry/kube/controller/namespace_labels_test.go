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

package controller

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/test"
)

func TestKubeNamespaceLabelsGetter(t *testing.T) {
	client := kube.NewFakeClient()
	stop := test.NewStop(t)

	namespaces := kclient.New[*corev1.Namespace](client)

	// Create some test namespaces
	prodNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "prod",
			Labels: map[string]string{
				"env":  "production",
				"tier": "frontend",
			},
		},
	}
	testNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"env": "testing",
			},
		},
	}
	noLabelsNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "no-labels",
			Labels: map[string]string{},
		},
	}

	client.Kube().CoreV1().Namespaces().Create(nil, prodNs, metav1.CreateOptions{})
	client.Kube().CoreV1().Namespaces().Create(nil, testNs, metav1.CreateOptions{})
	client.Kube().CoreV1().Namespaces().Create(nil, noLabelsNs, metav1.CreateOptions{})

	// Wait for namespaces to be synced
	client.RunAndWait(stop)

	// Create the namespace labels getter
	getter := NewKubeNamespaceLabelsGetter(namespaces)

	tests := []struct {
		name           string
		namespace      string
		expectedLabels map[string]string
	}{
		{
			name:      "namespace with labels",
			namespace: "prod",
			expectedLabels: map[string]string{
				"env":  "production",
				"tier": "frontend",
			},
		},
		{
			name:      "namespace with one label",
			namespace: "test",
			expectedLabels: map[string]string{
				"env": "testing",
			},
		},
		{
			name:           "namespace with no labels",
			namespace:      "no-labels",
			expectedLabels: map[string]string{},
		},
		{
			name:           "non-existent namespace",
			namespace:      "non-existent",
			expectedLabels: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels := getter.GetNamespaceLabels(tt.namespace)
			if !mapsEqual(labels, tt.expectedLabels) {
				t.Errorf("GetNamespaceLabels(%s) = %v, expected %v", tt.namespace, labels, tt.expectedLabels)
			}
		})
	}
}

func TestKubeNamespaceLabelsGetter_EmptyNamespace(t *testing.T) {
	client := kube.NewFakeClient()
	stop := test.NewStop(t)

	namespaces := kclient.New[*corev1.Namespace](client)

	client.RunAndWait(stop)

	getter := NewKubeNamespaceLabelsGetter(namespaces)

	// Test with empty namespace name
	labels := getter.GetNamespaceLabels("")
	if labels != nil {
		t.Errorf("GetNamespaceLabels(\"\") = %v, expected nil", labels)
	}
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
