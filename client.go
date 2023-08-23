package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

var DefaultClient = Client{}

func NewClient(baseUrl string, timeout ...int) *Client {
	secondsTimeOut := 60
	if len(timeout) > 0 {
		secondsTimeOut = timeout[0]
	}
	client := &Client{BaseURL: baseUrl, httpClient: &http.Client{
		Timeout: time.Duration(secondsTimeOut) * time.Second,
	}}
	return client
}

func (c *Client) getBaseUrl() string {
	if c.BaseURL == "" {
		return BASE_URL
	}
	return c.BaseURL
}

func (c *Client) getHttpClient() *http.Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: time.Minute,
		}
	}
	return c.httpClient
}

func (c *Client) get(data requestData) (string, error) {
	response, err := c.makeRequest(data, http.MethodGet)
	return response, err
}

func (c *Client) post(data requestData) (string, error) {
	response, err := c.makeRequest(data, http.MethodPost)
	return response, err
}

func (c *Client) delete(data requestData) error {
	_, err := c.makeRequest(data, http.MethodDelete)
	return err
}

func (c *Client) getNewRequest(data requestData, method string) (*http.Request, error) {
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0{
		return nil, errors.New("Missing API Key ")
	}
	url := c.getBaseUrl() + "/" + data.EntityEndpoint
	if data.Id != "" {
		url += "/" + data.Id
	}
	if data.QueryParams != "" {
		url += "?" + data.QueryParams
	}
	req, err := http.NewRequest(method, url, c.bodyToJson(data.Body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Add("X-Authorization", apiKey)
	return req, nil
}

func (c *Client) makeRequest(data requestData, method string) (string, error) {
	req, err := c.getNewRequest(data, method)
	if err != nil {
		return "", err
	}

	res, err := c.getHttpClient().Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if !c.validStatus(req.Method, res.StatusCode) {
		return "", errors.New(string(body))
	}
	return string(body), nil
}


func (c *Client) bodyToJson(body interface{}) io.Reader {
	var reader io.Reader
	if body != nil {
		inrec, err := json.Marshal(&body)
		if err != nil {
			return nil
		}
		reader = bytes.NewReader(inrec)
		return reader
	}
	return nil
}

func (c *Client) validStatus(method string, statusCode int) bool {
	switch method {
	case http.MethodPost:
		return statusCode == http.StatusCreated || statusCode == http.StatusOK
	case http.MethodGet:
		return statusCode == http.StatusOK
	case http.MethodDelete:
		return statusCode == http.StatusNoContent
	}
	return statusCode == http.StatusOK
}
