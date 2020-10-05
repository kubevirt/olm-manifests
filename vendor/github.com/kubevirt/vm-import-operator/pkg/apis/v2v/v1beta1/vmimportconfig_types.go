/*
Copyright 2020 The vm import Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sdkapi "kubevirt.io/controller-lifecycle-operator-sdk/pkg/sdk/api"
)

// VMImportConfigSpec defines the desired state of VMImportConfig
type VMImportConfigSpec struct {
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Rules on which nodes controller pod(s) will be scheduled
	// +optional
	Infra sdkapi.NodePlacement `json:"infra,omitempty"`
}

// VMImportConfigStatus defines the observed state of VMImportConfig
type VMImportConfigStatus struct {
	sdkapi.Status `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VMImportConfig is the Schema for the vmimportconfigs API
// +kubebuilder:resource:path=vmimportconfigs,scope=Namespaced
type VMImportConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMImportConfigSpec   `json:"spec,omitempty"`
	Status VMImportConfigStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VMImportConfigList contains a list of VMImportConfig
type VMImportConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMImportConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMImportConfig{}, &VMImportConfigList{})
}
