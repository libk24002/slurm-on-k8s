/*
Copyright 2025.

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

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ChartSpec struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Version    string `json:"version"`
	Namespace  string `json:"namespace,omitempty"`
}

type ServiceSpec struct {
	Type string `json:"type"`
}

type ValuesSpec struct {
	Service      ServiceSpec `json:"service"`
	ReplicaCount int32       `json:"replicaCount"`
}

// SlurmDeploymentSpec defines the desired state of SlurmDeployment.
type SlurmDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Chart  ChartSpec  `json:"chart"`
	Values ValuesSpec `json:"values,omitempty"`
}

// SlurmDeploymentStatus defines the observed state of SlurmDeployment.
type SlurmDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SlurmDeployment is the Schema for the slurmdeployments API.
type SlurmDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmDeploymentSpec   `json:"spec,omitempty"`
	Status SlurmDeploymentStatus `json:"status,omitempty"`
}

func (v *ValuesSpec) UnmarshalJSON(data []byte) error {
	type Alias ValuesSpec
	aux := &struct {
		Service      ServiceSpec `json:"service"`
		ReplicaCount interface{} `json:"replicaCount"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	v.Service = aux.Service
	switch val := aux.ReplicaCount.(type) {
	case int32:
		v.ReplicaCount = val
	case int:
		v.ReplicaCount = int32(val)
	case string:
		num, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		v.ReplicaCount = int32(num)
	case float64:
		// 处理 float64 类型
		v.ReplicaCount = int32(val)
	default:
		return fmt.Errorf("unexpected type for replicaCount: %T", val)
	}
	return nil
}

// +kubebuilder:object:root=true

// SlurmDeploymentList contains a list of SlurmDeployment.
type SlurmDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlurmDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlurmDeployment{}, &SlurmDeploymentList{})
}
