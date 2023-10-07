package controllers

import (
	"context"

	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *LearningReconciler) StatefulSet(ctx context.Context, req ctrl.Request) (ctrl.Request, error) {
	log := r.Log.WithValues("learning", req.NamespacedName)

	lrnOperator := &learningv1alpha1.Learning{}

	found := &appsv1.StatefulSet{}
	err := r.Get(ctx, types.NamespacedName{Name: lrnOperator.Name, Namespace: lrnOperator.Namespace}, found)
	if err != nil && errors.IsNotFound() {
		stateSet := r.StateFulSetForOperator(lrnOperator)
		log.Info("Deploying Statefulset")
		err = r.Create(ctx, stateSet)
		if err != nil {
			log.Error(err, "Failed to create the statefulset", "StatefulSet.Namespace", StatefulSet.Namespace, "StatefulSet.Name", StatefulSet.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
}
