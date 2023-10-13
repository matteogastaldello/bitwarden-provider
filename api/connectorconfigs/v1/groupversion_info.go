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
	ConnectorConfigKind             = reflect.TypeOf(ConnectorConfig{}).Name()
	ConnectorConfigGroupKind        = schema.GroupKind{Group: Group, Kind: ConnectorConfigKind}.String()
	ConnectorConfigKindAPIVersion   = ConnectorConfigKind + "." + SchemeGroupVersion.String()
	ConnectorConfigGroupVersionKind = SchemeGroupVersion.WithKind(ConnectorConfigKind)
)

func init() {
	SchemeBuilder.Register(&ConnectorConfig{}, &ConnectorConfigList{})
}
