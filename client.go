package deluge

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

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
	var res BoolResponse
	err := c.action("auth.login", fmt.Sprintf("\"%s\"", c.Password), &res)

	if err != nil {
		return err
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("Error pausing torrent: %s", res.Error.Message)
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
		Timeout:   time.Second * 10,
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
