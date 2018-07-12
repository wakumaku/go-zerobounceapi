package zerobounceapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	zerobounceapi "github.com/wakumaku/go-zerobounceapi"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *zerobounceapi.Client
)

var sandboxEmails = map[string]bool{
	"disposable@example.com":                  true,
	"invalid@example.com":                     false,
	"valid@example.com":                       true,
	"toxic@example.com":                       false,
	"donotmail@example.com":                   false,
	"spamtrap@example.com":                    false,
	"abuse@example.com":                       false,
	"unknown@example.com":                     false,
	"catch_all@example.com":                   false,
	"antispam_system@example.com":             false,
	"does_not_accept_mail@example.com":        false,
	"exception_occurred@example.com":          false,
	"failed_smtp_connection@example.com":      false,
	"failed_syntax_check@example.com":         false,
	"forcible_disconnect@example.com":         false,
	"global_suppression@example.com":          false,
	"greylisted@example.com":                  false,
	"leading_period_removed@example.com":      false,
	"mail_server_did_not_respond@example.com": false,
	"mail_server_temporary_error@example.com": false,
	"mailbox_quota_exceeded@example.com":      false,
	"mailbox_not_found@example.com":           false,
	"no_dns_entries@example.com":              false,
	"possible_trap@example.com":               false,
	"possible_typo@example.com":               false,
	"role_based@example.com":                  false,
	"timeout_exceeded@example.com":            false,
	"unroutable_ip_address@example.com":       false,
}

var (
	validResponse = `{
	"address":"flowerjill@aol.com",
	"status":"Valid",
	"sub_status":"", 
	"account":"flowerjill",
	"domain":"aol.com",
	"disposable":false,
	"toxic":false,
	"firstname":"Jill",
	"lastname":"Stein",
	"gender":"female",
	"location":null,
	"country":"United States",
	"region":"Florida",
	"city":"West Palm Beach",
	"zipcode":"33401",
	"creationdate":null,
	"processedat":"2017-04-01 02:48:02.592"
}`
	inValidResponse = `{
	"status":"Invalid"
}`
	validCreditResponse   = `{"credits":2375323}`
	inValidCreditResponse = `{"credits":-1}`
	errorResponse         = `{"error":"Invalid API Key or your account ran out of credits"}`
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = zerobounceapi.NewWith(server.URL, "a_valid_api_key", nil)

	return func() {
		server.Close()
	}
}

func TestValidate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Validate], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validResponse))
	})

	r, err := client.Validate("a_valid@email.com")
	if err != nil {
		t.Fatal(err)
	}

	if valid := r.IsValid(); !valid {
		t.Errorf("Should be valid: %v, %s", valid, err.Error())
	}
}

func TestValidateWithIP(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.ValidateWithIP], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validResponse))
	})

	r, err := client.ValidateWithIP("a_valid@email.com", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	if valid := r.IsValid(); !valid {
		t.Errorf("Should be valid: %v", valid)
	}
}

func TestValidateAnInvalidEmail(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Validate], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(inValidResponse))
	})

	r, err := client.Validate("invalid@email.com")
	if err != nil {
		t.Fatal(err)
	}

	if valid := r.IsValid(); valid {
		t.Errorf("Should be invalid: %v", valid)
	}
}

func TestInvalidAPIkey(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Validate], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(errorResponse))
	})

	r, err := client.Validate("test@email.com")
	if err != nil {
		t.Fatal(err)
	}

	if valid := r.IsValid(); valid {
		t.Errorf("Should be invalid: %v", valid)
	}
}

func TestCredit(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Credits], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validCreditResponse))
	})

	r, err := client.Credits()
	if err != nil {
		t.Fatal(err)
	}

	if available, err := r.CreditsAvailable(); !available {
		t.Errorf("Should be invalid: %v, %s", available, err.Error())
	}
}

func TestCreditBalance(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Credits], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(validCreditResponse))
	})

	r, err := client.Credits()
	if err != nil {
		t.Fatal(err)
	}

	if balance, err := r.CreditsBalance(); balance <= 0 {
		t.Errorf("Should be invalid: %v, %s", balance, err.Error())
	}
}

func TestCreditError(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(zerobounceapi.Endpoints[zerobounceapi.Credits], func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(inValidCreditResponse))
	})

	r, err := client.Credits()
	if err != nil {
		t.Fatal(err)
	}

	available, err := r.CreditsAvailable()
	if available {
		t.Errorf("Should not be available: %v, %s", available, err.Error())
	}

	if err != zerobounceapi.ErrCannotGetYourCreditBudget {
		t.Errorf("Error expected: %s\nReceived: %s", zerobounceapi.ErrCannotGetYourCreditBudget.Error(), err.Error())
	}
}
