package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *LearningReconciler) ServiceForOperator(m *learningv1alpha1.Learning) *corev1.Service {

	lbl := map[string]string{
		"app":    m.Name,
		"labels": m.Name,
	}

	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    lbl,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": m.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.Protocol(m.Spec.Service.Protocol),
					Port:       int32(m.Spec.Service.Port),
					TargetPort: intstr.FromInt(m.Spec.Service.TargetPort),
				},
			},
			Type: corev1.ServiceType(m.Spec.Service.Type),
		},
	}
	ctrl.SetControllerReference(m, svc, r.Scheme)
	return svc
}
