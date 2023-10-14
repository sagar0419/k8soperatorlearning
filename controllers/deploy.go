package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *LearningReconciler) DeploymentForOperator(l *learningv1alpha1.Learning) *appsv1.Deployment {
	labl := map[string]string{
		"app":    l.Name,
		"labels": l.Name,
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      l.Name,
			Labels:    labl,
			Namespace: l.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &l.Spec.AppSize,
			Selector: &metav1.LabelSelector{
				MatchLabels: labl,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labl,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  l.Spec.AppContainerName,
							Image: l.Spec.AppImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: l.Spec.AppPort,
									Name:          "http",
								},
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(l, deploy, r.Scheme)
	return deploy
}
