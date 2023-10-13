package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	connector "github.com/matteogastaldello/bitwarden-provider/api/connectorconfigs/v1"
	bwv1 "github.com/matteogastaldello/bitwarden-provider/api/secret/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		bwv1.SchemeBuilder.AddToScheme,
		connector.SchemeBuilder.AddToScheme,
	)
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
