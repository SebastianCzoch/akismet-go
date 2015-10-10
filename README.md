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

## Tests

```
$ go test ./...
````

## License

[GNU v2](./LICENSE)
