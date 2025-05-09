/*
Copyright 2022.

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

package v1beta1

import (
	kmmv1beta1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/kmm/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagedClusterModuleSpec defines the desired state of ManagedClusterModule
type ManagedClusterModuleSpec struct {
	// ModuleSpec describes how the KMM operator should deploy a Module on those nodes that need it.
	ModuleSpec kmmv1beta1.ModuleSpec `json:"moduleSpec,omitempty"`

	// SpokeNamespace describes the Spoke namespace, in which the ModuleSpec should be applied.
	SpokeNamespace string `json:"spokeNamespace"`

	// Selector describes on which managed clusters the ModuleSpec should be applied.
	Selector map[string]string `json:"selector"`
}

// ManagedClusterModuleStatus defines the observed state of ManagedClusterModule.
type ManagedClusterModuleStatus struct {
	// Number of ManifestWorks to be applied.
	NumberDesired int32 `json:"numberDesired,omitempty"`

	// Number of ManifestWorks that have been successfully applied.
	NumberApplied int32 `json:"numberApplied,omitempty"`

	// Number of ManifestWorks that could not be successfully applied.
	NumberDegraded int32 `json:"numberDegraded,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=managedclustermodules,scope=Cluster
//+kubebuilder:subresource:status

// ManagedClusterModule describes how to load a kernel module on managed clusters
// +operator-sdk:csv:customresourcedefinitions:displayName="Managed Cluster Module"
type ManagedClusterModule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterModuleSpec   `json:"spec,omitempty"`
	Status ManagedClusterModuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ManagedClusterModuleList contains a list of ManagedClusterModule
type ManagedClusterModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedClusterModule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedClusterModule{}, &ManagedClusterModuleList{})
}
