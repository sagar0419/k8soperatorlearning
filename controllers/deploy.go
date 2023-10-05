package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
)

func (r *LearningReconciler) DeploymentForOperator(l *learningv1alpha1.Learning) *appsv1.Deployment {
	return &appsv1.Deployment{}
}
