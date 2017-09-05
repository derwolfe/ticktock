package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const TIMEOUT = 3

const GITHUB = "https://status.github.com/api/status.json"
const CODECOV = "https://wdzsn5dlywj9.statuspage.io/api/v2/status.json"
const TRAVIS = "https://pnpcptp8xh9k.statuspage.io/api/v2/status.json"
const QUAY = "https://8szqd6w4s277.statuspage.io/api/v2/status.json"

type status struct {
	body []byte
	url  string
}

// make a parser for the different kinds of status messages. statuspage statuses tend to all contain the same elements.
type serviceState interface {
	// given a message, parse the relevant elements into a single booleanvalue
	parse(body []byte) (bool, error)

	// get the current state
	state()

	// mutate self to update its state
	update()
}

type statuspageStatus struct {
	state bool
}

type githubStatus struct {
	state bool
}

// func (s *statuspageStatus) parse(body []byte) (bool, error) {
//
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
