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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LearningSpec defines the desired state of Learning
type LearningSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Learning. Edit learning_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	ApplicationDescription string `json:"applicationDescription"`
	AppContainerName       string `json:"appContainerName"`
	AppImage               string `json:"appImage"`
	AppPort                int32  `json:"appPort"`
	AppSize                int32  `json:"appSize"`
	DatabaseDescription    string `json:"databaseDescription"`
	DbContainerName        string `json:"dbContainerName"`
	DbImage                string `json:"dbImage"`
	DbPort                 int32  `json:"dbPort"`
	DbStoragePath          string `json:"dbStoragePath"`
	DataStorageSize        string `json:"dataStorageSize"`
	DbSize                 int32  `json:"dbSize"`
}

type LearningServiceSpec struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	NodePort   string `json:"nodePort,omitempty"`
	Type       string `json:"type"`
}

// LearningStatus defines the observed state of Learning
type LearningStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodList []string `json:"podList"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Learning is the Schema for the learnings API
type Learning struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LearningSpec   `json:"spec,omitempty"`
	Status LearningStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LearningList contains a list of Learning
type LearningList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Learning `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Learning{}, &LearningList{})
}
