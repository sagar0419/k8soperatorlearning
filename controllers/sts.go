package controllers

import (
	learningv1alpha1 "k8soperatorlearning/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *LearningReconciler) StsForOperator(m *learningv1alpha1.Learning) *appsv1.StatefulSet {
	lbl := map[string]string{
		"app":    m.Spec.DbName,
		"labels": m.Spec.DbName,
	}
	// Converting termination grace period to int64
	gracePeriodSeconds := int64(30)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.DbName,
			Namespace: m.Spec.Namespace,
			Labels:    lbl,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &m.Spec.DbReplica,
			Selector: &metav1.LabelSelector{
				MatchLabels: lbl,
			},
			ServiceName: m.Spec.Service.Name,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lbl,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &gracePeriodSeconds,
					Containers: []v1.Container{
						{
							Name:  m.Spec.DbContainerName,
							Image: m.Spec.DbImage,
							Env: []v1.EnvVar{
								{
									Name:  m.Spec.Env.MysqlDb,
									Value: m.Spec.Env.VMysqlDb,
								},
								{
									Name:  m.Spec.Env.MysqlUser,
									Value: m.Spec.Env.VMysqlUser,
								},
								{
									Name:  m.Spec.Env.MysqlPassword,
									Value: m.Spec.Env.VMysqlPassword,
								},
								{
									Name:  m.Spec.Env.MysqlRootPassword,
									Value: m.Spec.Env.VMysqlRootPassword,
								},
							},
							Ports: []v1.ContainerPort{
								{
									Name:          "Http",
									ContainerPort: m.Spec.DbPort,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      m.Spec.DbVolumeName,
									MountPath: m.Spec.DbVolumePath,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: m.Spec.DbVolumePvcName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes: []v1.PersistentVolumeAccessMode{
							v1.ReadWriteOnce,
						},
						StorageClassName: &m.Spec.StorageClassNameMysql,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceStorage: resource.MustParse(m.Spec.DbVolumeSize),
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(m, sts, r.Scheme)
	return sts
}
