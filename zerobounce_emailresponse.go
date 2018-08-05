package zerobounceapi

import (
	"encoding/json"
	"errors"
	"strings"
)

// EmailResponse when validating emails
type EmailResponse struct {
	Address      string `json:"address"`      //:"flowerjill@aol.com",
	Status       string `json:"status"`       //:"Valid",
	SubStatus    string `json:"sub_status"`   //:"",
	Account      string `json:"account"`      //:"flowerjill",
	Domain       string `json:"domain"`       //:"aol.com",
	Disposable   bool   `json:"disposable"`   //:false,
	Toxic        bool   `json:"toxic"`        //:false,
	Firstname    string `json:"firstname"`    //:"Jill",
	Lastname     string `json:"lastname"`     //:"Stein",
	Gender       string `json:"gender"`       //:"female",
	Location     string `json:"location"`     //:null,
	Country      string `json:"country"`      //:"United States",
	Region       string `json:"region"`       //:"Florida",
	City         string `json:"city"`         //:"West Palm Beach",
	Zipcode      string `json:"zipcode"`      //:"33401",
	Creationdate string `json:"creationdate"` //:null,
	Processedate string `json:"processedate"` //:"2017-04-01 02:48:02.592"
	Err          string `json:"error"`        // In case of error
}

const (
	fieldEmailAddress      = "address"
	fieldEmailStatus       = "status"
	fieldEmailSubStatus    = "sub_status"
	fieldEmailAccount      = "account"
	fieldEmailDomain       = "domain"
	fieldEmailDisposable   = "disposable"
	fieldEmailToxic        = "toxic"
	fieldEmailFirstname    = "firstname"
	fieldEmailLastname     = "lastname"
	fieldEmailGender       = "gender"
	fieldEmailLocation     = "location"
	fieldEmailCountry      = "country"
	fieldEmailRegion       = "region"
	fieldEmailCity         = "city"
	fieldEmailZipcode      = "zipcode"
	fieldEmailCreationdate = "creationdate"
	fieldEmailProcessedate = "processedate"
	fieldEmailError        = "error"
)

// NewEmailResponseFromJSON creates a CreditResponse from a raw json
func NewEmailResponseFromJSON(jsonBytes []byte) (*EmailResponse, error) {

	var response map[string]interface{}
	err := json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}

	return parseEmailResponse(response), nil
}

func parseEmailResponse(response map[string]interface{}) *EmailResponse {

	emailResponse := &EmailResponse{
		Address:      getStringField(response[fieldEmailAddress]),
		Status:       getStringField(response[fieldEmailStatus]),
		SubStatus:    getStringField(response[fieldEmailSubStatus]),
		Account:      getStringField(response[fieldEmailAccount]),
		Domain:       getStringField(response[fieldEmailDomain]),
		Disposable:   getBoolField(response[fieldEmailDisposable]),
		Toxic:        getBoolField(response[fieldEmailToxic]),
		Firstname:    getStringField(response[fieldEmailFirstname]),
		Lastname:     getStringField(response[fieldEmailLastname]),
		Gender:       getStringField(response[fieldEmailGender]),
		Location:     getStringField(response[fieldEmailLocation]),
		Country:      getStringField(response[fieldEmailCountry]),
		Region:       getStringField(response[fieldEmailRegion]),
		City:         getStringField(response[fieldEmailCity]),
		Zipcode:      getStringField(response[fieldEmailZipcode]),
		Creationdate: getStringField(response[fieldEmailCreationdate]),
		Processedate: getStringField(response[fieldEmailProcessedate]),
		Err:          getStringField(response[fieldEmailError]),
	}

	return emailResponse
}

// Error produced
func (r *EmailResponse) Error() error {
	return errors.New(r.Err)
}

// IsValid checks in a response if an email have been validated
func (r *EmailResponse) IsValid() bool {
	return strings.ToLower(r.Status) == "valid"
}
