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

	// TemplateRef specifies the reference to the PromotionTemplate.
	// +required
	TemplateRef TemplateRef `json:"templateRef"`

	// Strategy specifies how to promote.
	// +required
	Strategy Strategy `json:"strategy"`

	// A list of resources to be included in the readiness check.
	// +optional
	ReadinessChecks ReadinessChecks `json:"readinessChecks,omitempty"`
}

// TypedLocalObjectReference defines the readiness checks to be done before doing the promotion.
type ReadinessChecks struct {
	// A list of objects (in the same namespace) to be included in the readiness check.
	LocalObjectsRef []LocalObjectsRef `json:"localObjectsRef"`
}

type LocalObjectsRef struct {
	GroupVersionResource metav1.GroupVersionResource `json:"groupVersionResource"`

	Name string `json:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty"`
}

func (in *Promotion) GetLocalObjectsRefsForReadinessChecks() []LocalObjectsRef {
	return in.Spec.ReadinessChecks.LocalObjectsRef
}

// FromSpec defines the source of the promotion.
type FromSpec struct {
	EnvironmentRef EnvironmentReference `json:"environmentRef"`
}

// ToSpec defines the destination of the promotion.
type ToSpec struct {
	EnvironmentRef EnvironmentReference `json:"environmentRef"`
}

// TemplateRef defines the reference to the PromotionTemplate.
type TemplateRef struct {
	Name string `json:"name"`
}

// Strategy defines the strategy for the promotion.
type Strategy struct {
	PullRequest bool `json:"pull-request"`
}

// PromotionStatus defines the observed state of Promotion
type PromotionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions holds the conditions for the Promotion.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// DependentObjectsReady ...
	// +optional
	DependentObjectsReady bool `json:"dependentObjectsReady"`
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
