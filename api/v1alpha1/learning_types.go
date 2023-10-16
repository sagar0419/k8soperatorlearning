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
	Namespace             string              `json:"namespace"`
	AppContainerName      string              `json:"appContainerName"`
	AppImage              string              `json:"appImage"`
	AppName               string              `json:"appName"`
	AppPort               int32               `json:"appPort"`
	AppReplica            int32               `json:"appReplica"`
	DbContainerName       string              `json:"dbContainerName"`
	DbImage               string              `json:"dbImage"`
	DbPort                int32               `json:"dbPort"`
	DbReplica             int32               `json:"dbReplica"`
	DbName                string              `json:"dbName"`
	DbVolumeName          string              `json:"dbVolumeName"`
	DbVolumePath          string              `json:"dbVolumePath"`
	DbVolumePvcName       string              `json:"dbVolumePvcName"`
	DbVolumeSize          string              `json:"dbVolumeSize"`
	StorageClassNameMysql string              `json:"storageClassNameMysql"`
	Service               LearningServiceSpec `json:"service"`
	Env                   LearningDbEnvVar    `json:"env"`
}

type LearningServiceSpec struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
	NodePort   string `json:"nodePort,omitempty"`
	Type       string `json:"type,omitempty"`
}

type LearningDbEnvVar struct {
	MysqlDb            string `json:"mysqlDb"` // Old DB
	MysqlUser          string `json:"mysqlUser"`
	MysqlPassword      string `json:"mysqlPassword"`
	MysqlRootPassword  string `json:"mysqlRootPassword"`
	VMysqlDb           string `json:"vMysqlDb"`
	VMysqlUser         string `json:"vMysqlUser"`
	VMysqlPassword     string `json:"vMysqlPassword"`
	VMysqlRootPassword string `json:"vMysqlRootPassword"`
	AppDb              string `json:"appDb"` //DB that application will use.
	VAppDb             string `json:"vAppDb"`
	DbHostName         string `json:"dbHostName"`
	VDbHostName        string `json:"vDbHostName"`
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
