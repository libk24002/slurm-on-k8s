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

type MariaDBSpec struct {
	Enabled bool               `json:"enabled"`
	Port    int32              `json:"port"`
	Auth    MariaDBAuthSpec    `json:"auth"`
	Primary MariaDBPrimarySpec `json:"primary"`
}

type MariaDBAuthSpec struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database"`
}

type MariaDBPrimarySpec struct {
	Persistence MariaDBPrimaryPersistenceSpec `json:"persistence"`
}

type MariaDBPrimaryPersistenceSpec struct {
	Enabled      bool   `json:"enabled"`
	StorageClass string `json:"storageClass"`
	Size         string `json:"size"`
}

type AuthSpec struct {
	SSH AuthSSHSpec `json:"ssh"`
}

type AuthSSHSpec struct {
	Secret    AuthSSHSecretSpec    `json:"secret"`
	ConfigMap AuthSSHConfigmapSpec `json:"configmap"`
}

type AuthSSHSecretSpec struct {
	Name string                `json:"name"`
	Keys AuthSSHSecretKeysSpec `json:"keys"`
}

type AuthSSHSecretKeysSpec struct {
	Public         string `json:"public"`
	Private        string `json:"private"`
	AuthorizedKeys string `json:"authorizedKeys"`
}

type AuthSSHConfigmapSpec struct {
	Name          string   `json:"name"`
	PrefabPubKeys []string `json:"prefabPubKeys"`
}

type PersistenceSpec struct {
	Shared PersistenceSharedSpec `json:"shared"`
}

type PersistenceSharedSpec struct {
	Enabled       bool     `json:"enabled"`
	Name          string   `json:"name"`
	ExistingClaim string   `json:"existingClaim"`
	AccessModes   []string `json:"accessModes"`
	StorageClass  string   `json:"storageClass"`
	Size          string   `json:"size"`
}

type ImageSpec struct {
	Registry    string   `json:"registry"`
	Repository  string   `json:"repository"`
	Tag         string   `json:"tag"`
	PullPolicy  string   `json:"pullPolicy"`
	PullSecrets []string `json:"pullSecrets"`
}

type DiagnosticModeSpec struct {
	Enabled bool     `json:"enabled"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

type ExtraVolumeMountsSpec struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type MungedSpec struct {
	Name              string                  `json:"name"`
	CommonLables      []string                `json:"commonLables"`
	Image             ImageSpec               `json:"image"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode"`
	ExtraVolumes      []map[string]string     `json:"extraVolumes"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts"`
}

type SlurmctldSpec struct {
	Name              string                  `json:"name"`
	CommonLables      []string                `json:"commonLables"`
	Image             ImageSpec               `json:"image"`
	CheckDNS          SlurmctldCheckDNS       `json:"checkDns"`
	ReplicaCount      int32                   `json:"replicaCount"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode"`
	ExtraVolumes      []map[string]string     `json:"extraVolumes"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts"`
}

type SlurmctldCheckDNS struct {
	Image ImageSpec `json:"image"`
}

type SlurmdSpec struct {
	Name              string                  `json:"name"`
	CommonLables      []string                `json:"commonLables"`
	Image             ImageSpec               `json:"image"`
	ReplicaCount      int32                   `json:"replicaCount"`
	Resources         ResourceSpec            `json:"resources"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode"`
	ExtraVolumes      []map[string]string     `json:"extraVolumes"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts"`
}

type ResourceSpec struct {
	Requests ResourceRequestSpec `json:"requests"`
	Limits   ResourceLimitSpec   `json:"limits"`
}

type ResourceRequestSpec struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	EphemeralStorage string `json:"ephemeral-storage"`
}

type ResourceLimitSpec struct {
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	EphemeralStorage string `json:"ephemeral-storage"`
}

type SlurmdbdSpec struct {
	Name              string                  `json:"name"`
	CommonLables      []string                `json:"commonLables"`
	Image             ImageSpec               `json:"image"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode"`
	ExtraVolumes      []map[string]string     `json:"extraVolumes"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts"`
}

type SlurmLogindSpec struct {
	Name              string                  `json:"name"`
	CommonLables      []string                `json:"commonLables"`
	Image             ImageSpec               `json:"image"`
	Resources         ResourceSpec            `json:"resources"`
	DiagnosticMode    DiagnosticModeSpec      `json:"diagnosticMode"`
	ExtraVolumes      []map[string]string     `json:"extraVolumes"`
	ExtraVolumeMounts []ExtraVolumeMountsSpec `json:"extraVolumeMounts"`
}

// type ServiceAccountSpec struct {
// 	Create bool `json:"create"`
// }

// type SlurmConfigSpec struct {
// 	SlurmConf string `json:"slurm.conf"`
// }

type ValuesSpec struct {
	Mariadb     MariaDBSpec     `json:"mariadb,omitempty"`
	Auth        AuthSpec        `json:"auth,omitempty"`
	Persistence PersistenceSpec `json:"persistence,omitempty"`
	Munged      MungedSpec      `json:"munged,omitempty"`
	Slurmctld   SlurmctldSpec   `json:"slurmctld,omitempty"`
	Slurmd      SlurmdSpec      `json:"slurmd,omitempty"`
	Slurmdbd    SlurmdbdSpec    `json:"slurmdbd,omitempty"`
	SlurmLogin  SlurmLogindSpec `json:"login,omitempty"`
	// ServiceAccount    ServiceAccountSpec `json:"serviceAccount,omitempty"`
	// SlurmConfig       SlurmConfigSpec    `json:"configuration,omitempty"`
	NameOverride      string   `json:"nameOverride"`
	FullnameOverride  string   `json:"fullnameOverride"`
	CommonAnnotations []string `json:"commonAnnotations"`
	CommonLables      []string `json:"commonLables"`
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
	// aux := &struct {
	// 	Service      ServiceSpec `json:"service"`
	// 	ReplicaCount interface{} `json:"replicaCount"`
	// }{}
	// if err := json.Unmarshal(data, &aux); err != nil {
	// 	return err
	// }
	// v.Service = aux.Service
	// switch val := aux.ReplicaCount.(type) {
	// case int32:
	// 	v.ReplicaCount = val
	// case int:
	// 	v.ReplicaCount = int32(val)
	// case string:
	// 	num, err := strconv.ParseInt(val, 10, 32)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	v.ReplicaCount = int32(num)
	// case float64:
	// 	// 处理 float64 类型
	// 	v.ReplicaCount = int32(val)
	// default:
	// 	return fmt.Errorf("unexpected type for replicaCount: %T", val)
	// }
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
