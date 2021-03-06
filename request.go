package goinsta

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"errors"
)

type reqOptions struct {
	Endpoint     string
	PostData     string
	IsLoggedIn   bool
	IgnoreStatus bool
	Query        map[string]string
}

type ErrResponse struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
}

func (insta *Instagram) OptionalRequest(endpoint string, a ...interface{}) (body []byte, err error) {
	return insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf(endpoint, a...),
	})
}

func (insta *Instagram) sendSimpleRequest(endpoint string, a ...interface{}) (body []byte, err error) {
	return insta.sendRequest(&reqOptions{
		Endpoint: fmt.Sprintf(endpoint, a...),
	})
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, err error) {

	if !insta.IsLoggedIn && !o.IsLoggedIn {
		return nil, fmt.Errorf("not logged in")
	}

	method := "GET"
	if len(o.PostData) > 0 {
		method = "POST"
	}

	u, err := url.Parse(GOINSTA_API_URL + o.Endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range o.Query {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	var req *http.Request
	req, err = http.NewRequest(method, u.String(), bytes.NewBuffer([]byte(o.PostData)))
	if err != nil {
		return
	}

	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", GOINSTA_USER_AGENT)

	client := &http.Client{
		Jar: insta.Cookiejar,
	}

	if insta.Proxy != "" {
		proxy, err := url.Parse(insta.Proxy)
		if err != nil {
			return body, err
		}
		insta.Transport.Proxy = http.ProxyURL(proxy)

		client.Transport = &insta.Transport
	} else {
		// Remove proxy if insta.Proxy was removed
		insta.Transport.Proxy = nil
		client.Transport = &insta.Transport
	}

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	u, _ = url.Parse(GOINSTA_API_URL)
	for _, value := range insta.Cookiejar.Cookies(u) {
		if strings.Contains(value.Name, "csrftoken") {
			insta.Informations.Token = value.Value
		}
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = insta.checkResponseError(resp.StatusCode, body)
	if err != nil {
		return nil, err
	}

	return body, err
}

//Checks if the response is successful, if not parses the error
func (insta *Instagram) checkResponseError(code int, body []byte) error {

	switch code { //Basic code flipping, anything not caught here is treated as a new errResp
	case 200:
		return nil
	case 429:
		return ErrRateLimit
	case 404:
		return ErrNotFound
	}

	var errResp ErrResponse
	err := json.Unmarshal(body, &errResp)
	if err != nil {
		return fmt.Errorf("invalid status code %s", string(body)) //Cant unmarshal so skip
	}

	switch errResp.ErrorType {
	case "rate_limit_error":
		return ErrRateLimit
	}

	return errors.New(errResp.Message)
}
