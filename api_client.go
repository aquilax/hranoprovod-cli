package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/Hranoprovod/shared"
)

type APINode *shared.APINode
type APINodeList *shared.APINodeList

type APIClientOptions struct {
	BaseURL string
}

func GetDefaultAPIClientOptions() *APIClientOptions {
	return &APIClientOptions{
		"http://hranoprovod.appspot.com/api/v1/",
	}
}

type APIClient struct {
	aco *APIClientOptions
}

func NewAPIClient(aco *APIClientOptions) *APIClient {
	return &APIClient{
		aco,
	}
}

func httpRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func (ac *APIClient) Search(q string) (*APINodeList, error) {
	url := ac.aco.BaseURL + "search?q=" + q
	resp, err := httpRequest("GET", url, nil)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var nl APINodeList;
	err = json.Unmarshal(contents, &nl);
	if err != nil {
		return nil, err
	}
	return &nl, nil
}