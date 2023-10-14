package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *LearningReconciler) HeadLessSvc(m *learningv1alpha1.Learning) *corev1.Service {

	lbl := map[string]string{
		"app":    m.Spec.DbName,
		"labels": m.Spec.DbName,
	}

	targetPort := intstr.FromInt(m.Spec.Service.TargetPort)

	hdlSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.DbName,
			Namespace: m.Spec.Namespace,
			Labels:    lbl,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       m.Spec.DbPort,
					TargetPort: targetPort,
				},
			},
			ClusterIP: "None",
			Selector:  lbl,
		},
	}
	ctrl.SetControllerReference(m, hdlSvc, r.Scheme)
	return hdlSvc
}
