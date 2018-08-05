package zerobounceapi

import (
	"encoding/json"
	"errors"
)

// CreditResponse received when checking our balance
type CreditResponse struct {
	Credits int    `json:"credits"` // -1 in case of failure
	Err     string `json:"error"`   // In case of error
}

const (
	fieldCreditCredits = "credits"
	fieldCreditError   = "error"
)

// NewCreditResponseFromJSON creates a CreditResponse from a raw json
func NewCreditResponseFromJSON(jsonBytes []byte) (*CreditResponse, error) {

	var response map[string]interface{}
	err := json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}

	return parseCreditResponse(response), nil
}

func parseCreditResponse(response map[string]interface{}) *CreditResponse {

	creditResponse := &CreditResponse{
		Credits: getIntField(response[fieldCreditCredits]),
		Err:     getStringField(response[fieldCreditError]),
	}

	return creditResponse
}

// Error produced
func (r *CreditResponse) Error() error {
	return errors.New(r.Err)
}

// CreditsAvailable checks if there are credits available
func (r *CreditResponse) CreditsAvailable() (bool, error) {
	if r.Credits == -1 {
		return false, ErrCannotGetYourCreditBudget
	}

	if r.Err != "" {
		return r.Credits > 0, errors.New(r.Err)
	}
	return r.Credits > 0, nil
}
