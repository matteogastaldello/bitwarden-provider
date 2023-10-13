package unlocker

import (
	"context"
	"net/http"
	"path"

	rtv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	"github.com/lucasepe/httplib"
	bwclient "github.com/matteogastaldello/bitwarden-provider/internal/clients"
	"github.com/matteogastaldello/bitwarden-provider/internal/resolvers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
Sblocca il Vault

Le API di bitwarden non richiedono autenticazione ad ogni chiamata,
ma il Vault viene sbloccato "una tantum" prima di avviare le chiamate
*/
func Unlock(ctx context.Context, cli *bwclient.Client, kubeCli *client.Client, ref *rtv1.Reference) (*bwclient.Response, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(bwclient.Default),
		Path:    path.Join("unlock"),
	}).Build()
	if err != nil {
		return nil, err
	}

	opts, err := resolvers.ResolveConnectorConfig(ctx, *kubeCli, ref)

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	values := struct {
		Password string `json:"password"`
	}{
		Password: opts.Password,
	}
	req, err := httplib.Post(uri.String(), httplib.ToJSON(values))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)

	apiErr := &bwclient.APIError{}

	val := &bwclient.Response{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})
	return val, err
}
