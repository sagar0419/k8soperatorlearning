package controllers

import (
	"context"
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *LearningReconciler) Deploy(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Learning", req.NamespacedName)
	lrnOperator := &learningv1alpha1.Learning{}
	// checking deployment
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: lrnOperator.Name, Namespace: lrnOperator.Namespace}, found)
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
	return ctrl.Result{}, nil
}
