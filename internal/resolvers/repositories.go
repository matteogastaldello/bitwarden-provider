package resolvers

import (
	"context"
	"fmt"

	rtv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	repositories "github.com/matteogastaldello/bitwarden-provider/api/secret/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/types"
)

func ResolveGitRepository(ctx context.Context, kube client.Client, ref *rtv1.Reference) (repositories.BitwardenSecret, error) {
	res := repositories.BitwardenSecret{}
	if ref == nil {
		return res, fmt.Errorf("no %s referenced", res.Kind)
	}

	err := kube.Get(ctx, types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}, &res)
	return res, err
}
