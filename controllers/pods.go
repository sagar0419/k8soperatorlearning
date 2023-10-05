package controllers

import (
	corev1 "k8s.io/api/core/v1"
)

func GetPodNames(pods []corev1.Pod) []string {
	var podNames []string

	for _, pods := range pods {
		podNames = append(podNames, pods.Name)
	}
	return podNames
}
