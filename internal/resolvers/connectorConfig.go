package resolvers

import (
	"context"
	"fmt"

	rtv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	"github.com/krateoplatformops/provider-runtime/pkg/resource"
	connectorconfigs "github.com/matteogastaldello/bitwarden-provider/api/connectorconfigs/v1"
	bwclient "github.com/matteogastaldello/bitwarden-provider/internal/clients"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ResolveConnectorConfig(ctx context.Context, kube client.Client, ref *rtv1.Reference) (bwclient.ClientOptions, error) {
	opts := bwclient.ClientOptions{}

	cfg := connectorconfigs.ConnectorConfig{}
	if ref == nil {
		return opts, fmt.Errorf("no %s referenced", cfg.Kind)
	}

	err := kube.Get(ctx, types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}, &cfg)
	if err != nil {
		return opts, errors.Wrapf(err, "cannot get %s connector config", ref.Name)
	}

	csr := cfg.Spec.Credentials.SecretRef
	if csr == nil {
		return opts, fmt.Errorf("no credentials secret referenced")
	}

	sec := corev1.Secret{}
	err = kube.Get(ctx, types.NamespacedName{Namespace: csr.Namespace, Name: csr.Name}, &sec)
	if err != nil {
		return opts, errors.Wrapf(err, "cannot get %s secret", ref.Name)
	}

	password, err := resource.GetSecret(ctx, kube, csr.DeepCopy())
	if err != nil {
		return opts, err
	}

	opts.ApiUrl = cfg.Spec.ApiUrl
	opts.Password = password

	return opts, nil
}
