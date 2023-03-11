/*
Copyright 2023 Thomas Stadler <thomas@thomasst.xyz>.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnvironmentSpec defines the desired state of Environment
type EnvironmentSpec struct {
	// Source specifies the source Git Repository.
	// +required
	Source *SourceSpec `json:"source"`

	// Path to the directory which represents the environment.
	// Defaults to 'None', which translates to the root path of the Source.
	// +optional
	Path string `json:"path,omitempty"`
}

// SourceSpec includes the Git reference of the source Git Repository.
type SourceSpec struct {
	// URL specifies the Git repository URL, it can be an HTTP/S or SSH address.
	// +kubebuilder:validation:Pattern="^(http|https|ssh)://.*$"
	// +required
	URL string `json:"url"`

	// Reference specifies the Git reference to resolve and monitor for
	// changes, defaults to the 'master' branch.
	// +optional
	Reference *GitRepositoryRef `json:"ref,omitempty"`
}

// GitRepositoryRef specifies the Git reference to resolve and checkout.
type GitRepositoryRef struct {
	// Branch to check out, defaults to 'master' if no other field is defined.
	// +optional
	Branch string `json:"branch,omitempty"`
}

// EnvironmentStatus defines the observed state of Environment
type EnvironmentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Environment is the Schema for the environments API
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EnvironmentList contains a list of Environment
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
