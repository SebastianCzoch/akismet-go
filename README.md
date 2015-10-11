# akismet-go [![Build Status](https://travis-ci.org/SebastianCzoch/akismet-go.svg?branch=master)](https://travis-ci.org/SebastianCzoch/akismet-go) [![Code Climate](https://codeclimate.com/github/SebastianCzoch/akismet-go/badges/gpa.svg)](https://codeclimate.com/github/SebastianCzoch/akismet-go) [![GoDoc](https://godoc.org/github.com/SebastianCzoch/akismet-go?status.svg)](https://godoc.org/github.com/SebastianCzoch/akismet-go)  [![License](https://img.shields.io/badge/licence-GNU%20v2-green.svg)](./LICENSE)



Go utils for working with [Akismet](http://www.akismet.com/) spam detection service.

## Examples
### Example 1 - Verification Client
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/akismet-go"
)

func main() {
	client := akismet.NewClient("api_key", "site")
	if client.VeryfiClient() != nil {
		fmt.Println("Client can't be verified, ", err)
	}
}
```

### Example 2 - Check spam
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/akismet-go"
)

func main() {
	client := akismet.NewClient("api_key", "site")
	if client.VeryfiClient() != nil {
		fmt.Println("Client can't be verified, ", err)
	}

	options := &akismet.Options{
		UserIP: "127.0.0.1",
		UserAgent: "Test-Agent",
		Content: "It's a spam?",
	}

	isSpam, err := client.IsSpam(options)
	if err != nil {
		fmt.Println(err)
	}

	if isSpam {
		fmt.Println("This is spam")
	} else {
		fmt.Println("This is not spam")
	}
}
```
## Install

```
$ go get github.com/SebastianCzoch/akismet-go
````

or via [Godep](https://github.com/tools/godep)
```
$ godep get github.com/SebastianCzoch/akismet-go
```


## API
### NewClient(apiKey, site string) *Client
Create new client and return pointer to it

### (c *Client) VeryfiClient() (error)
Check if passed key and blog values are correct, if not return error

### (c *Client) IsSpam(o *Options) (bool, error)
Check if passed Options struct is a spam or not

### (c *Client) SubmitSpam(o *Options) error
This call is for submitting comments that weren't marked as spam but should have been.

### (c *Client) SubmitHam(o *Options) error
This call is intended for the submission of false positives - items that were incorrectly classified as spam by Akismet.

### Options struct
```
	UserIP      string (required) IP address of the comment submitter
	UserAgent   string (required) User agent string of the web browser submitting the comment
	Referrer    string The content of the HTTP_REFERER header should be sent here
	Permalink   string The permanent location of the entry the comment was submitted to
	Author      string Name submitted with the comment
	AuthorEmail string Email address submitted with the comment
	AuthorURL   string URL submitted with comment
	Content     string The content that was submitted
	Created     string Datetime when content was created (RFC 3339 format)
	Modified    string Datetime when content was modified (RFC 3339 format)
	Lang        string Indicates the language(s) in use on the blog or site, in ISO 639-1 format, comma-separated. A site with articles in English and French might use "en, fr_ca"
	Charset     string The character encoding for the form values, such as "UTF-8" or "ISO-8859-1"
	UserRole    string The user role of the user who submitted the comment. This is an optional parameter. If you set it to "administrator", Akismet will always return false.
	IsTest      string This is an optional parameter. You can use it when submitting test queries to Akismet.
```
## Tests
Required go in version >=1.4

```
$ go test ./...
````

## License

[GNU v2](./LICENSE)
