package deluge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// func (c *Client) url(path string) string {
// 	if path == "" || path[0:1] != "/" {
// 		path = fmt.Sprintf("/%s", path)
// 	}

// 	if c.token != "" {
// 		path = fmt.Sprintf("%s&token=%s", path, c.token)
// 	}
// 	return fmt.Sprintf("%s%s", c.API, path)
// }

func (c *Client) request(method, path string, payload []byte, headers *http.Header) (*http.Response, error) {
	if c == nil {
		return nil, fmt.Errorf("Cannot make a request with a nil client")
	}
	in := bytes.NewBuffer(payload)
	req, err := http.NewRequest(method, c.API, in)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for header, values := range *headers {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}

	res, err := c.user_agent.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) post(path string, payload []byte, headers *http.Header) (*http.Response, error) {
	return c.request("POST", path, payload, headers)
}

func (c *Client) put(path string, payload []byte, headers *http.Header) (*http.Response, error) {
	return c.request("PUT", path, payload, headers)
}

func (c *Client) get(path string, headers *http.Header) (*http.Response, error) {
	return c.request("GET", path, nil, headers)
}

func (c *Client) delete(path string, headers *http.Header) (*http.Response, error) {
	return c.request("DELETE", path, nil, headers)
}

// func (c *Client) action(action string, hash string, headers *http.Header) error {
// 	res, err := c.get(fmt.Sprintf("/?action=%s&hash=%s", action, hash), headers)
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode != 200 {
// 		return fmt.Errorf("status code: %d", res.StatusCode)
// 	}
// 	return nil
// }

func (c *Client) action(method string, params string, decoder interface{}) error {
	var payload = fmt.Sprintf(`{"id":%d, "method":"%s", "params":[%s]}`, c.index, method, params)
	c.index++

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	res, err := c.post("", []byte(payload), &header)

	//Deluge hangs if the action is invalid or hash doesnt match a torrent
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return errors.New("request timed out. Check to make sure the action is valid for the speficied torrent")
	} else if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("error status: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("unable to read response body")
	}

	err = json.Unmarshal(body, &decoder)
	if err != nil {
		return errors.New("unable to parse response body: " + err.Error())
	}

	return nil
}
