/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Module kubernetes object
type Module struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Properties ModuleProperties `json:"properties,omitempty"`

	Status ModuleStatus `json:"status,omitempty"`
}

type ModuleProperties struct {
	Weight int    `json:"weight"`
	State  string `json:"state"`
	Source string `json:"source"`
}

type ModuleStatus struct{}

type moduleKind struct{}

func (ms *ModuleStatus) GetObjectKind() schema.ObjectKind {
	return &moduleKind{}
}

func (mk *moduleKind) SetGroupVersionKind(_ schema.GroupVersionKind) {}
func (mk *moduleKind) GroupVersionKind() schema.GroupVersionKind {
	return ModuleGVK
}

var ModuleGVK = schema.GroupVersionKind{Group: "deckhouse.io", Version: "v1alpha1", Kind: "Module"}

func (m *Module) SetName(name string) {
	m.Name = name
	m.calculateLabels()
}

func (m *Module) SetWeight(weight int) {
	m.Properties.Weight = weight
}

func (m *Module) SetSource(source string) {
	if source == "" {
		source = "Embedded"
	}

	if source != "Embedded" {
		source = "External: " + source
	}

	m.Properties.Source = source
}
func (m *Module) SetEnabledState(enabled bool) {
	if enabled {
		m.Properties.State = "Enabled"
	} else {
		m.Properties.State = "Disabled"
	}
}

func (m *Module) calculateLabels() {
	// could be removed when we will ready properties from the module.yaml file

	if strings.HasPrefix(m.Name, "cni-") {
		m.Labels["module.deckhouse.io/cni"] = ""
	}

	if strings.HasPrefix(m.Name, "cloud-provider-") {
		m.Labels["module.deckhouse.io/cloud-provider"] = ""
	}

	if strings.HasSuffix(m.Name, "-crd") {
		m.Labels["module.deckhouse.io/crd"] = ""
	}
}
