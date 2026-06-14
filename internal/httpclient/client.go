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
	baseHeader string
	clientHTTP *http.Client
	baseAuth   string
}

func New() *Client {
	return &Client{
		baseHeader: "application/x-www-form-urlencoded",
		clientHTTP: &http.Client{},
		baseAuth:   "",
	}
}

func (c *Client) MakeRequest(url string, body string, method Method) (string, error) {
	req, err := http.NewRequest(string(method), url, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", c.baseAuth)
	req.Header.Set("Content-Type", c.baseHeader)

	resp, err := c.clientHTTP.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Client) ChangeHeader(header string) {
	c.baseHeader = header
}

func (c *Client) ChangeAuth(auth string) {
	c.baseAuth = auth
}
