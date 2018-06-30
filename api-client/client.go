package client

import (
	"encoding/json"
	"github.com/aquilax/hranoprovod-cli/shared"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Options contains the option for the API Client
type Options struct {
	BaseURL string
}

// URLParams is used to add reuqest params to the URL
type URLParams map[string]string

// NewDefaultOptions returns the default options for the API Client
func NewDefaultOptions() *Options {
	return &Options{
		"http://hranoprovod.appspot.com/api/v1/",
	}
}

// APIClient the base client structure
type APIClient struct {
	options *Options
}

// NewAPIClient returns new API Client
func NewAPIClient(options *Options) *APIClient {
	return &APIClient{
		options,
	}
}

func httpRequest(method string, url string, body io.Reader, v interface{}) error {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var ae shared.APIError
		err = json.Unmarshal(contents, &ae)
		if err != nil {
			return err
		}
		return ae
	}
	return json.Unmarshal(contents, v)
}

func (ac *APIClient) buildURL(path string, params URLParams) (string, error) {
	u, err := url.Parse(ac.options.BaseURL)
	if err != nil {
		return "", err
	}
	p := strings.Split(strings.Trim(u.Path, "/"), "/")
	p2 := strings.Split(strings.Trim(path, "/"), "/")
	p = append(p, p2...)
	u.Path = strings.Join(p, "/")

	qp := url.Values{}
	for name, value := range params {
		qp.Set(name, value)
	}
	u.RawQuery = qp.Encode()
	return u.String(), nil
}

// Search searches the service for the provided query
func (ac *APIClient) Search(q string) (*shared.APINodeList, error) {
	url, err := ac.buildURL("search", URLParams{"q": q})
	if err != nil {
		return nil, err
	}
	var nl shared.APINodeList
	err = httpRequest("GET", url, nil, &nl)
	return &nl, err
}
