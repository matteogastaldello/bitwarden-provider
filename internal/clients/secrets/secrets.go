package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"

	"github.com/lucasepe/httplib"
	bwv1 "github.com/matteogastaldello/bitwarden-provider/api/secret/v1"
	client "github.com/matteogastaldello/bitwarden-provider/internal/clients"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Data interface{}

type ResponseGeneric struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
	Data    Data    `json:"data"`
}

func (m *ResponseGeneric) ValidData() error {

	if reflect.TypeOf(m.Data) != reflect.TypeOf(bwv1.Secret{}) || m.Message == nil {
		fmt.Println(reflect.TypeOf(m.Data), reflect.TypeOf(bwv1.Secret{}), "message:", m.Message)
		return errors.New("Not a valid data field")
	}
	return nil
	// switch m.Data.(type) {
	// case string:
	// 	// handle string
	// 	return nil
	// case bwv1.Secret:
	// 	// handle bwv1.Secret
	// 	return nil
	// default:
	// 	// handle unexpected type
	// 	return errors.New("Not a valid data field")
	// }
}

// type SecretResponse struct {
// 	Success bool        `json:"success"`
// 	Data    bwv1.Secret `json:"data"`
// }

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
func (s *ResponseGeneric) Get(ctx context.Context, cli *client.Client, secretId string) (*bwv1.Secret, error) {
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
	val := &bwv1.Secret{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})

	return val, err
}
func Exists(ctx context.Context, cli *client.Client, secretId string) (*bool, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(client.Default),
		Path:    path.Join("object/item", secretId),
	}).Build()
	retVal := false
	if err != nil {
		return &retVal, nil
	}
	req, err := httplib.Get(uri.String())
	if err != nil {
		return &retVal, nil
	}
	req = req.WithContext(ctx)

	//apiErr := &client.APIError{}
	val := &ResponseGeneric{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		// Validators: []httplib.HandleResponseFunc{
		// 	httplib.ErrorJSON(apiErr, http.StatusOK), httplib.ErrorJSON(apiErr, http.StatusBadRequest),
		// },
	})
	if err != nil {
		return &retVal, err
	}

	// err = val.ValidData()
	// if err != nil {
	// 	return false, err
	// }

	return &val.Success, nil
}

func Create(ctx context.Context, cli *client.Client, secret bwv1.Secret) (*bwv1.Secret, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(client.Default),
		Path:    path.Join("object/item"),
	}).Build()
	if err != nil {
		return nil, errors.New("error occurred building uri")
	}
	b, err := json.Marshal(secret)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	req, err := httplib.Post(uri.String(), httplib.ToJSON(secret))
	if err != nil {
		return nil, errors.New("error occurred making POST req")
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)

	apiErr := &client.APIError{}
	val := &ResponseGeneric{}
	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		Validators: []httplib.HandleResponseFunc{
			httplib.ErrorJSON(apiErr, http.StatusOK),
		},
	})
	if err != nil {
		return nil, err
	}
	sec := bwv1.Secret{}
	err = mapstructure.Decode(val.Data, &sec)
	if err != nil {
		return nil, err
	}
	return &sec, nil
}

func Delete(ctx context.Context, cli *client.Client, secretId string) (bool, error) {
	uri, err := httplib.NewURLBuilder(httplib.URLBuilderOptions{
		BaseURL: cli.BaseURL(client.Default),
		Path:    path.Join("object/item", secretId),
	}).Build()
	if err != nil {
		return false, nil
	}
	req, err := httplib.Delete(uri.String())
	if err != nil {
		return false, nil
	}
	req = req.WithContext(ctx)

	//apiErr := &client.APIError{}
	val := &ResponseGeneric{}

	err = httplib.Fire(cli.HTTPClient(), req, httplib.FireOptions{
		ResponseHandler: httplib.FromJSON(val),
		// Validators: []httplib.HandleResponseFunc{
		// 	httplib.ErrorJSON(apiErr, http.StatusOK), httplib.ErrorJSON(apiErr, http.StatusBadRequest),
		// },
	})
	if err != nil {
		return false, err
	}
	b, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	if err != nil {
		return false, err
	}

	// err = val.ValidData()
	// if err != nil {
	// 	return false, err
	// }

	return val.Success, nil
}
