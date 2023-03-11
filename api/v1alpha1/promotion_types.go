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

// PromotionSpec defines the desired state of Promotion
type PromotionSpec struct {
	// FromSpec specifies where to promote from.
	// +required
	FromSpec FromSpec `json:"from"`

	// ToSpec specifies where to promote to.
	// +required
	ToSpec ToSpec `json:"to"`

	// Strategy specifies how to promote.
	// +required
	Strategy Strategy `json:"strategy"`
}

// FromSpec defines the source of the promotion.
type FromSpec struct {
	EnvironmentRef EnvironmentReference `json:"environmentRef"`
}

// ToSpec defines the destination of the promotion.
type ToSpec struct {
	EnvironmentRef EnvironmentReference `json:"environmentRef"`
}

// Strategy defines the strategy for the promotion.
type Strategy struct {
	PullRequest bool `json:"pull-request"`
}

// PromotionStatus defines the observed state of Promotion
type PromotionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Promotion is the Schema for the promotions API
type Promotion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromotionSpec   `json:"spec,omitempty"`
	Status PromotionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PromotionList contains a list of Promotion
type PromotionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Promotion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Promotion{}, &PromotionList{})
}
