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

package controllers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/thomasstxyz/release-promotion-operator/api/v1alpha1"
)

// PromotionReconciler reconciles a Promotion object
type PromotionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.release-promotion-operator.io,resources=promotions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.release-promotion-operator.io,resources=promotions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.release-promotion-operator.io,resources=promotions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Promotion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PromotionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Begin Reconciliation", "ctx", ctx, "req", req)

	// Get Promotion object
	promotion := &apiv1alpha1.Promotion{}
	if err := r.Get(ctx, req.NamespacedName, promotion); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Do readiness checks
	var unreadyResources []apiv1alpha1.LocalObjectsRef
	ReadinessChecksSucceeded, unreadyResources, err := r.readinessChecks(ctx, promotion)
	if err != nil {
		fmt.Println(err)
	}

	// Update status of Promotion
	if err := r.updateStatus(ctx, promotion, unreadyResources, ReadinessChecksSucceeded); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// readinessChecks checks the status of all dependent objects for readiness
// using [kstatus](https://github.com/kubernetes-sigs/cli-utils/tree/master/pkg/kstatus)
func (r *PromotionReconciler) readinessChecks(ctx context.Context, promotion *apiv1alpha1.Promotion) (bool, []apiv1alpha1.LocalObjectsRef, error) {
	// Create dynamic client to be able to get unstructured resources
	dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())

	var dependentResources []apiv1alpha1.LocalObjectsRef = promotion.GetLocalObjectsRefsForReadinessChecks()

	// Check ready status of each specified object
	var unreadyResources []apiv1alpha1.LocalObjectsRef
	for _, dr := range dependentResources {
		gvr := schema.GroupVersionResource{
			Group:    dr.GroupVersionResource.Group,
			Version:  dr.GroupVersionResource.Version,
			Resource: dr.GroupVersionResource.Resource,
		}

		var ns string
		if dr.Namespace != "" {
			ns = dr.Namespace
		} else {
			ns = promotion.Namespace
		}

		// Fetch dependent object
		var obj *unstructured.Unstructured
		obj, err = dynamicClient.Resource(gvr).Namespace(ns).Get(ctx, dr.Name, v1.GetOptions{}, "")
		if err != nil {
			fmt.Printf("error getting obj: %v\n", err)
			return false, unreadyResources, err
		}

		// Check status of dependent object
		result, err := status.Compute(obj)
		if err != nil {
			fmt.Printf("error retrieving status of obj: %v\n", err)
			return false, unreadyResources, err
		}
		fmt.Printf("The status of %s is: %q. Message: %s\n", dr.Name, result.Status, result.Message)

		// Mark object unready if status is any other than status.CurrentStatus (Ready status)
		if result.Status != status.CurrentStatus {
			unreadyResources = append(unreadyResources, dr)
		}
	}

	return true, unreadyResources, err
}

func (r *PromotionReconciler) updateStatus(ctx context.Context, promotion *apiv1alpha1.Promotion, unreadyResources []apiv1alpha1.LocalObjectsRef, ReadinessChecksSucceeded bool) error {
	// get beforeConditions
	var beforeConditions []metav1.Condition = promotion.Status.Conditions

	var afterConditions []metav1.Condition = beforeConditions

	// Get index of existing condition of type "Ready"
	var readyTypeConditionIndex int
	var readyTypeConditionExists bool
	for i, bc := range beforeConditions {
		switch bc.Type {
		case "Ready":
			readyTypeConditionIndex = i
			readyTypeConditionExists = true
		}
	}

	if !ReadinessChecksSucceeded {
		promotion.Status.DependentObjectsReady = false

		condition := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Reason:             "ReadinessChecksFailed",
			// TODO: output the object, which failed
			Message: "Error fetching dependent objects!",
		}

		if readyTypeConditionExists {
			afterConditions[readyTypeConditionIndex] = condition
		} else {
			afterConditions = append(afterConditions, condition)
		}
	}

	if ReadinessChecksSucceeded {
		promotion.Status.DependentObjectsReady = true

		condition := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Reason:             "ReadinessChecksSucceeded",
			Message:            "All dependent objects are ready!",
		}

		if readyTypeConditionExists {
			afterConditions[readyTypeConditionIndex] = condition
		} else {
			afterConditions = append(afterConditions, condition)
		}
	}

	if len(unreadyResources) != 0 && ReadinessChecksSucceeded {
		promotion.Status.DependentObjectsReady = false

		condition := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Reason:             "ReadinessChecksFailed",
			Message:            fmt.Sprintf("Dependent objects are not ready:\n%v", unreadyResources),
		}

		if readyTypeConditionExists {
			afterConditions[readyTypeConditionIndex] = condition
		} else {
			afterConditions = append(afterConditions, condition)
		}
	}

	promotion.Status.Conditions = afterConditions

	err := r.Status().Update(ctx, promotion, &client.SubResourceUpdateOptions{})
	if err != nil {
		fmt.Printf("error updating status of resource: %v\n", err)
		return err
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *PromotionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Promotion{}).
		Complete(r)
}
