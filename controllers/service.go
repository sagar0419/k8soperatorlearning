package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

func (r *LearningReconciler) ServiceForOperator(m *learningv1alpha1.Learning) *corev1.Service {
	return &corev1.Service{}
}
