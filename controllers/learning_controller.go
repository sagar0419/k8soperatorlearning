/*
Copyright 2023 Sagar.

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
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"
)

// LearningReconciler reconciles a Learning object
type LearningReconciler struct {
	client.Client
	// logger interface that allows you to log messages with different severity levels and structured data.
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=learning.sagar,resources=learnings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=learning.sagar,resources=learnings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=learning.sagar,resources=learnings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Learning object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LearningReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Learning", req.NamespacedName)

	lrnOperator := &learningv1alpha1.Learning{}
	// check if the operator is deployed or not. If there is a issue with operator it will redeploy it, if the operator is deleted it will ignore.
	err := r.Client.Get(ctx, req.NamespacedName, lrnOperator)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Learning Operator is not found, Ignoring since the object is deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get opertaor")
		return ctrl.Result{}, err
	}

	// checking deployment
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: lrnOperator.Name, Namespace: lrnOperator.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		deploy := r.DeploymentForOperator(lrnOperator)
		log.Info("Creating new deployment")
		err = r.Create(ctx, deploy)
		if err != nil {
			log.Error(err, "Failed to create the deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, err
	} else if err != nil {
		log.Error(err, "Failed to get the deployment manifest")
		return ctrl.Result{}, err
	}

	// Updating deployment pod spec template if deployment already exist
	deploy := r.DeploymentForOperator(lrnOperator)
	if !equality.Semantic.DeepDerivative(deploy.Spec.Template, found.Spec.Template) {
		found = deploy
		log.Info("updating deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		err := r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update the deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// updating the replicas of the pods
	size := lrnOperator.Spec.AppSize

	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size

		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update the deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Checking Service
	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: foundService.Name, Namespace: foundService.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		dep := r.ServiceForOperator(lrnOperator)
		log.Info("Creating Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to Get the service", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to update the pod list status")
		return ctrl.Result{}, err
	}

	// PodList
	podlist := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(found.Namespace),
		client.MatchingLabels(map[string]string{"app": found.Name, "labels": found.Name}),
	}
	if r.List(ctx, podlist, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
		return ctrl.Result{}, err
	}

	podNames := GetPodNames(podlist.Items)
	if !reflect.DeepEqual(podNames, lrnOperator.Status.PodList) {
		lrnOperator.Status.PodList = podNames
		err = r.Status().Update(ctx, lrnOperator)
		if err != nil {
			log.Error(err, "Failed to update the pod list status")
			return ctrl.Result{}, err
		}
	}

	// Checking Statefulset
	sts := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Namespace: lrnOperator.Namespace, Name: lrnOperator.Name}, sts)
	if err != nil && errors.IsNotFound(err) {
		stsDeploy := r.StsForOperator(lrnOperator)
		log.Info("Creating New StatefulSet")
		err = r.Create(ctx, stsDeploy)
		if err != nil {
			log.Error(err, "Failed to Create a StatefulSet", "StatefulSet.Namespace", sts.Namespace, "StatefulSet.Name", sts.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, err
	} else if err != nil {
		log.Error(err, "Failed to get the StatefulSet manifest")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LearningReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learningv1alpha1.Learning{}).
		Complete(r)
}
