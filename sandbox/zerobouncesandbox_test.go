package zerobounceapisandbox_test

import (
	"os"
	"testing"

	zerobounceapi "github.com/wakumaku/go-zerobounceapi"
)

var sandboxEmails = map[string]bool{
	"disposable@example.com":                  true,
	"invalid@example.com":                     false,
	"valid@example.com":                       true,
	"toxic@example.com":                       true,
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
	"leading_period_removed@example.com":      true,
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

func TestSandbox(t *testing.T) {
	client := zerobounceapi.New(os.Getenv("ZB_APIKEY"), nil)

	for email, valid := range sandboxEmails {
		r, err := client.EmailValidation(email)
		if err == nil {
			if isValid := r.IsValid(); isValid != valid {
				t.Fatalf("Expected %s to be %v", email, valid)
			}
		} else {
			t.Fatal(err)
		}
	}
}
