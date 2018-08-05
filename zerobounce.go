// Package zerobounceapi client based on
// https://docs.zerobounce.net/docs
package zerobounceapi

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// BaseURL of zerobounce services endpoint
const BaseURL = "https://api.zerobounce.net/v1"

// ServiceType .
type ServiceType int

// Service types supported
const (
	_ ServiceType = iota
	Validate
	ValidateWithIP
	Credits
)

// Endpoints endpoint definitions
var Endpoints = map[ServiceType]string{
	Validate:       "/validate",
	ValidateWithIP: "/validatewithip",
	Credits:        "/getcredits",
}

// Errors
var (
	ErrEmptyResponse             = errors.New("Empty body response received from service")
	ErrResponseHadEmptyStructure = errors.New("Could not unmarshall correctly the response")
	ErrCannotGetYourCreditBudget = errors.New("Cannot get your credits at this moment (-1)")
)

// Client holding credentials and connection
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Error returned
type Error struct {
	section string
	action  string
	err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.section, e.action, e.err)
}

// New Client
func New(apiKey string, httpClient *http.Client) *Client {
	return NewWithBaseURL(BaseURL, apiKey, httpClient)
}

// NewWithBaseURL a baseURL for test stuff
func NewWithBaseURL(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 2 * time.Second,
		}
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// EmailValidation calls the validation email service
func (c *Client) EmailValidation(email string) (*EmailResponse, error) {
	params := map[string]string{"email": email}
	return c.callValidateService(Validate, params)
}

// EmailValidationWithIP calls the validation email service with an IP
func (c *Client) EmailValidationWithIP(email, ip string) (*EmailResponse, error) {
	params := map[string]string{"email": email, "ipaddress": ip}
	return c.callValidateService(ValidateWithIP, params)
}

// CreditBalance calls the credit service
func (c *Client) CreditBalance() (*CreditResponse, error) {
	params := map[string]string{}
	return c.callCreditService(params)
}

func (c *Client) callValidateService(service ServiceType, params map[string]string) (*EmailResponse, error) {
	request, err := c.buildRequest("GET", Endpoints[service], params, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.callService(request)
	if err != nil {
		return nil, err
	}

	return NewEmailResponseFromJSON(body)
}

func (c *Client) callCreditService(params map[string]string) (*CreditResponse, error) {
	request, err := c.buildRequest("GET", Endpoints[Credits], params, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.callService(request)

	return NewCreditResponseFromJSON(body)
}

func (c *Client) callService(request *http.Request) ([]byte, error) {
	_, body, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	return body, nil
}

func (c *Client) buildRequest(method, path string, params map[string]string, body io.Reader) (*http.Request, error) {
	// Default required values
	params["apikey"] = c.apiKey

	URL, err := c.buildURL(path, params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, &Error{"buildRequest", "NewRequest", err}
	}

	return req, nil
}

func (c *Client) buildURL(path string, params map[string]string) (string, error) {

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", &Error{"buildURL", "parseURL", err}
	}
	u.Path += path

	queryString := u.Query()
	for k, v := range params {
		queryString.Set(k, v)
	}

	u.RawQuery = queryString.Encode()

	return u.String(), nil
}

func (c *Client) doRequest(request *http.Request) (int, []byte, error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return 0, nil, &Error{"doRequest", "Do", err}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, &Error{"doRequest", "ReadAll", err}
	}

	return resp.StatusCode, body, nil
}
