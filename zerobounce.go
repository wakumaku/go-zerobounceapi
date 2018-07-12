// Package zerobounceapi client based on
// https://docs.zerobounce.net/docs
package zerobounceapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// EmailResponse when validating emails
type EmailResponse struct {
	Address      string `json:"address,omitempty"`      //:"flowerjill@aol.com",
	Status       string `json:"status,omitempty"`       //:"Valid",
	SubStatus    string `json:"sub_status,omitempty"`   //:"",
	Account      string `json:"account,omitempty"`      //:"flowerjill",
	Domain       string `json:"domain,omitempty"`       //:"aol.com",
	Disposable   bool   `json:"disposable,omitempty"`   //:false,
	Toxic        bool   `json:"toxic,omitempty"`        //:false,
	Firstname    string `json:"firstname,omitempty"`    //:"Jill",
	Lastname     string `json:"lastname,omitempty"`     //:"Stein",
	Gender       string `json:"gender,omitempty"`       //:"female",
	Location     string `json:"location,omitempty"`     //:null,
	Country      string `json:"country,omitempty"`      //:"United States",
	Region       string `json:"region,omitempty"`       //:"Florida",
	City         string `json:"city,omitempty"`         //:"West Palm Beach",
	Zipcode      string `json:"zipcode,omitempty"`      //:"33401",
	Creationdate string `json:"creationdate,omitempty"` //:null,
	Processedat  string `json:"processedat,omitempty"`  //:"2017-04-01 02:48:02.592"
	// In case of error
	Error string `json:"error,omitempty"`
}

// CreditResponse received when checking our balance
type CreditResponse struct {
	Credits int `json:"credits,omitempty"` // -1 in case of failure
	// In case of error
	Error string `json:"error,omitempty"`
}

// IsValid checks in a response if an email have been validated
func (r *EmailResponse) IsValid() bool {
	return strings.ToLower(r.Status) == "valid"
}

// CreditsAvailable checks if there are credits available
func (r *CreditResponse) CreditsAvailable() (bool, error) {
	if r.Credits == -1 {
		return false, ErrCannotGetYourCreditBudget
	}

	if r.Error != "" {
		return r.Credits > 0, errors.New(r.Error)
	}
	return r.Credits > 0, nil
}

// CreditsBalance returns the current balance
func (r *CreditResponse) CreditsBalance() (int, error) {
	if r.Credits == -1 {
		return r.Credits, ErrCannotGetYourCreditBudget
	}

	if r.Error != "" {
		return r.Credits, errors.New(r.Error)
	}
	return r.Credits, nil
}

// Client holding credentials and connection
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New Client
func New(apiKey string, httpClient *http.Client) *Client {
	return NewWith(BaseURL, apiKey, httpClient)
}

// NewWith a baseURL for test stuff
func NewWith(baseURL, apiKey string, httpClient *http.Client) *Client {
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

// Validate .
func (c *Client) Validate(email string) (*EmailResponse, error) {
	params := map[string]string{"email": email}
	return c.callValidateService(Validate, params)
}

// ValidateWithIP .
func (c *Client) ValidateWithIP(email, ip string) (*EmailResponse, error) {
	params := map[string]string{"email": email, "ipaddress": ip}
	return c.callValidateService(ValidateWithIP, params)
}

// Credits gets the current credits available
func (c *Client) Credits() (*CreditResponse, error) {
	params := map[string]string{}
	return c.callCreditService(params)
}

func (c *Client) callValidateService(service ServiceType, params map[string]string) (*EmailResponse, error) {
	request, err := c.buildRequest("GET", Endpoints[service], params, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.callService(request)

	var response EmailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, err
}

func (c *Client) callCreditService(params map[string]string) (*CreditResponse, error) {
	request, err := c.buildRequest("GET", Endpoints[Credits], params, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.callService(request)

	var response CreditResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, err
}

func (c *Client) callService(request *http.Request) ([]byte, error) {
	_, body, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("Error doing request: %s", err.Error())
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

	return http.NewRequest(method, URL, body)
}

func (c *Client) buildURL(path string, params map[string]string) (string, error) {

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
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
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}
