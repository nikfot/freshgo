package gvhttp

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//HTTPClient implements a struct for http client
type HTTPClient struct {
	Description string
	Timeout     time.Duration
	Headers     map[string]string
	ServerURL   *url.URL
	httpClient  *http.Client
}

//HTTPRequest is a struct fot http requests
type HTTPRequest struct {
	request *http.Request
}

//RequestsClient is a generic HTTP client with time out
var RequestsClient = &http.Client{
	Timeout: time.Second * 10,
}

//SetHTTPTimeout sets timeout to Alethea http client
func (client *HTTPClient) SetHTTPTimeout(timeout time.Duration) {
	client.httpClient.Timeout = timeout
}

//Request implements a struct of http request method "GET"
func (client *HTTPClient) Request(method string, urlStr string, requestBody []byte, username string, password string, headers map[string]string) ([]byte, error) {
	//headers := client.Headers
	method = strings.ToUpper(method)
	var r HTTPRequest
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r, err = NewHTTPRequest(method, url.String(), requestBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if username != "" {
		r.SetRequestAuth(username, password)
	}
	if headers != nil {
		r.SetHeaders(headers)
	}

	res, err := client.httpClient.Do(r.request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//SetRequestHeader sets headers to Alethea http requests
func (r *HTTPRequest) SetRequestHeader(key string, value string) {
	r.request.Header.Set(key, value)
}

//SetRequestAuth sets headers to Alethea http requests
func (r *HTTPRequest) SetRequestAuth(usr string, psw string) {
	r.request.SetBasicAuth(usr, psw)
}

//NewHTTPRequest wraps the HTTPRequest of the package around classic http request
func NewHTTPRequest(method string, url string, requestBody []byte) (request HTTPRequest, err error) {
	var req *http.Request
	if strings.ToUpper(method) == "GET" {
		req, err = http.NewRequest(strings.ToUpper(method), url, nil)
	} else if strings.ToUpper(method) == "POST" {
		req, err = http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(requestBody))
	}
	if err != nil {
		return HTTPRequest{}, err
	}
	request.request = req
	return request, nil
}

//SetHeaders adds multiple request headers
func (r *HTTPRequest) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		r.SetRequestHeader(k, v)
	}
}

//NewHTTPClient creates a new HTTPClient
func NewHTTPClient(description string, url string, timeout time.Duration, headers map[string]string, ignoreInsecure bool) *HTTPClient {
	//init vars
	var client HTTPClient
	client.Headers = make(map[string]string)
	//set values
	client.Description = description
	client.Timeout = timeout
	client.Headers = headers
	if ignoreInsecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.httpClient = &http.Client{Transport: tr}
	} else {
		client.httpClient = &http.Client{}
	}
	//Add timeout
	client.SetHTTPTimeout(timeout)
	return &client
}
