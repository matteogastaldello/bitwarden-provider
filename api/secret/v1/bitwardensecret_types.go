/*
Copyright 2023.

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
	rtv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  int    `json:"type"`
}

type URI struct {
	Match int    `json:"match"`
	URI   string `json:"uri"`
}

type Login struct {
	Uris     []URI  `json:"uris"`
	Username string `json:"username"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}

// BitwardenSecretSpec defines the desired state of BitwardenSecret
type BitwardenSecretSpec struct {
	rtv1.ManagedSpec `json:",inline"`

	Credentials    *rtv1.CredentialSelectors `json:"credentials"`
	OrganizationID string                    `json:"organizationid"`
	CollectionIDs  string                    `json:"collectionids"`
	FolderID       *string                   `json:"folderid"`
	Type           int                       `json:"type"`
	Name           string                    `json:"name"`
	Notes          *string                   `json:"notes"`
	Favorite       bool                      `json:"favorite"`
	Fields         []Field                   `json:"fields"`
	Login          Login                     `json:"login"`
	Reprompt       int                       `json:"reprompt"`
}

// BitwardenSecretStatus defines the observed state of BitwardenSecret
type BitwardenSecretStatus struct {
	rtv1.ManagedStatus `json:",inline"`

	SecretId string `json:"bitwardenId,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BitwardenSecret is the Schema for the bitwardensecrets API
type BitwardenSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BitwardenSecretSpec   `json:"spec,omitempty"`
	Status BitwardenSecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BitwardenSecretList contains a list of BitwardenSecret
type BitwardenSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BitwardenSecret `json:"items"`
}
