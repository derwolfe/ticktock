package main

import (
	"encoding/json"
	"fmt"
	"github.com/derwolfe/ticktock/parsing"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TIMEOUT = 3 //seconds
	GITHUB  = "https://status.github.com/api/status.json"
	CODECOV = "https://wdzsn5dlywj9.statuspage.io/api/v2/status.json"
	TRAVIS  = "https://pnpcptp8xh9k.statuspage.io/api/v2/status.json"
	QUAY    = "https://8szqd6w4s277.statuspage.io/api/v2/status.json"
)

type status struct {
	body []byte
	url  string
}

// var parsers make(map[string]Parser) {
// 	"CODECOV":
// 	"QUAY":
// 	"TRAVIS":
// 	"GITHUB":
// }

func statusFetch(url string) (*status, error) {
	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	// how to know if there was a timeout?
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &status{body: body, url: url}, nil
}

func main() {
	sources := []string{GITHUB, TRAVIS, QUAY, CODECOV}
	for _, url := range sources {
		status, err := statusFetch(url)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(status.body))
	}
}
