# zerobounce API Client [![Build Status](https://travis-ci.org/wakumaku/go-zerobounceapi.svg?branch=master)](https://travis-ci.org/wakumaku/go-zerobounceapi) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/9b66f7d42dcb413bbf96f8f4d1471020)](https://www.codacy.com/app/wakumaku/go-zerobounceapi?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=wakumaku/go-zerobounceapi&amp;utm_campaign=Badge_Grade) [![Code Coverage](https://scrutinizer-ci.com/g/wakumaku/go-zerobounceapi/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/wakumaku/go-zerobounceapi/?branch=master) [![GoDoc](https://godoc.org/github.com/wakumaku/go-zerobounceapi?status.svg)](https://godoc.org/github.com/wakumaku/go-zerobounceapi)
### Source: https://docs.zerobounce.net/docs

```
go get github.com/wakumaku/go-zerobounceapi
```

Email validation:
```
client = zerobounceapi.New(apiKey, nil)
response, err := client.EmailValidation("email@domain.tld")
if err != nil {
    panic(err)
}

fmt.Println(response.IsValid())
```

Get credits:
```
client = zerobounceapi.New(apiKey, nil)
response, err := client.CreditBalance()
if err != nil {
    panic(err)
}

fmt.Println(response.Credits)
```

Makefile:
* `make test` Runs tests
* `make test-sandbox` Runs tests against zerobounce services, you need a valid API key (will not consume credits: https://docs.zerobounce.net/docs/sandbox-mode)