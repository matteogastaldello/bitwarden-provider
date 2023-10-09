package secrets

import (
	"context"
	"net/http"
	"path"

	"github.com/lucasepe/httplib"
	client "github.com/matteogastaldello/bitwarden-provider/internal/clients"
)

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

type Secret struct {
	OrganizationID string  `json:"organizationid"`
	CollectionIDs  string  `json:"collectionids"`
	FolderID       *string `json:"folderid"`
	Type           int     `json:"type"`
	Name           string  `json:"name"`
	Notes          *string `json:"notes"`
	Favorite       bool    `json:"favorite"`
	Fields         []Field `json:"fields"`
	Login          Login   `json:"login"`
	Reprompt       int     `json:"reprompt"`
}

type SecretResponse struct {
	Success bool   `json:"success"`
	Data    Secret `json:"data"`
}

// func Get(ctx context.Context, cli *azuredevops.Client, opts GetOptions) (bool, map[string]interface{}) {
// 	resp, err := http.Get(fmt.Sprintf("%s/object/item/%s", bw_server, id))

// 	if err != nil {
// 		log.Error(err, "Errore", "GET")
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return false, nil
// 	}

// 	var res map[string]interface{}

// 	json.NewDecoder(resp.Body).Decode(&res)
// 	// str, _ := json.Marshal(res)
// 	// fmt.Println(fmt.Printf("\n\n%s\n\n", str))
// 	return res["success"].(bool), res
// }

// Get retrieve a git repository.
// GET https://dev.azure.com/{organization}/{project}/_apis/git/repositories/{repositoryId}?api-version=7.0
func (s *SecretResponse) Get(ctx context.Context, cli *client.Client, secretId string) (*Secret, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(client.Default),
		Path:    path.Join("object/item", secretId),
	}).Build()
	if err != nil {
		return nil, err
	}
	req, err := httplib.Get(uri.String())
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	apiErr := &client.APIError{}
	val := &Secret{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})

	return val, err
}
func Exists(ctx context.Context, cli *client.Client, secretId string) (bool, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(client.Default),
		Path:    path.Join("object/item", secretId),
	}).Build()
	if err != nil {
		return false, err
	}
	req, err := httplib.Get(uri.String())
	if err != nil {
		return false, err
	}
	req = req.WithContext(ctx)

	apiErr := &client.APIError{}
	val := &SecretResponse{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})

	return val.Success, err
}
