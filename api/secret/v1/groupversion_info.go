// Package v1 contains API Schema definitions for the github v1 API group
// +kubebuilder:object:generate=true
// +groupName=bitwarden.provider.matteogastaldello.provider
// +versionName=v1
package v1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	Group   = "bitwarden.provider.matteogastaldello.provider"
	Version = "v1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

var (
	BitwardenSecretKind             = reflect.TypeOf(BitwardenSecret{}).Name()
	BitwardenSecretGroupKind        = schema.GroupKind{Group: Group, Kind: BitwardenSecretKind}.String()
	BitwardenSecretAPIVersion       = BitwardenSecretKind + "." + SchemeGroupVersion.String()
	BitwardenSecretGroupVersionKind = SchemeGroupVersion.WithKind(BitwardenSecretKind)
)

func init() {
	SchemeBuilder.Register(&BitwardenSecret{}, &BitwardenSecretList{})
}
