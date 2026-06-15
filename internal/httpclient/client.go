package httpclient

import (
	"io"
	"net/http"
	"strings"
)

type Method string

const (
	GET    Method = http.MethodGet
	POST   Method = http.MethodPost
	PUT    Method = http.MethodPut
	PATCH  Method = http.MethodPatch
	DELETE Method = http.MethodDelete
)

type Client struct {
	clientHTTP *http.Client
}

func New() *Client {
	return &Client{clientHTTP: &http.Client{}}
}

func (c *Client) MakeRequest(method Method, url string, body string, headers map[string]string) (string, error) {
	req, err := http.NewRequest(string(method), url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.clientHTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer func(){
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
