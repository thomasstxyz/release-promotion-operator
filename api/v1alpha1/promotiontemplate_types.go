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

// PromotionTemplateSpec defines the desired state of PromotionTemplate
type PromotionTemplateSpec struct {
	// CopySpec contains a list of source/destination pairs,
	// which represent file copy operations
	// between the source and destination environment.
	// +required
	CopySpec []CopyOperation `json:"copy"`
}
type CopyOperation struct {
	// Source is the path in the source environment.
	// Can be either a file or a directory.
	// +required
	Source string `json:"source"`

	// Destination is the path in the destination environment.
	// Can be either a file or a directory.
	// +required
	Destination string `json:"destination"`
}

// PromotionTemplateStatus defines the observed state of PromotionTemplate
type PromotionTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PromotionTemplate is the Schema for the promotiontemplates API
type PromotionTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromotionTemplateSpec   `json:"spec,omitempty"`
	Status PromotionTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PromotionTemplateList contains a list of PromotionTemplate
type PromotionTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PromotionTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PromotionTemplate{}, &PromotionTemplateList{})
}
