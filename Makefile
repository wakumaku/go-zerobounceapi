.PHONY: test-sandbox check-env

test:
	@go test -v --race .

test-sandbox: check-env
	@go test -v --race ./sandbox/.

check-env:
ifndef ZB_APIKEY
$(error ZB_APIKEY is not set)
endif