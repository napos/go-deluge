package deluge

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	API      string
	Username string
	Password string

	token      string
	user_agent *http.Client
}

func (c *Client) setToken() error {
	var rr RpcResponse
	err := c.action("auth.login", fmt.Sprintf("\"%s\"", c.Password), &rr)

	if err != nil {
		return err
	}
	if !rr.Result {
		return fmt.Errorf("error code %d! %s", rr.Error.Code, rr.Error.Message)
	}

	return nil
}

func NewClient(c *Client) (*Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	cookieJar, _ := cookiejar.New(&options)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	c.user_agent = &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	if c.API == "" {
		c.API = "http://localhost:8112/json"
	}

	err := c.setToken()
	if err != nil {
		return nil, err
	}

	return c, nil
}
