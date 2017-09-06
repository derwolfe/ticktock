package main

import (
	"fmt"
	"github.com/derwolfe/ticktock/parsing"
	"github.com/derwolfe/ticktock/state"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	TIMEOUT = 3 //seconds
	GITHUB  = "https://status.github.com/api/status.json"
	CODECOV = "https://wdzsn5dlywj9.statuspage.io/api/v2/status.json"
	TRAVIS  = "https://pnpcptp8xh9k.statuspage.io/api/v2/status.json"
	QUAY    = "https://8szqd6w4s277.statuspage.io/api/v2/status.json"
)

func statusFetch(url string) (*[]byte, error) {
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
	return &body, nil
}

type Refined struct {
	Url           string
	LastUpdated   time.Time
	Good          bool
	SourceMessage *[]byte
}

func main() {
	sources := []string{GITHUB, TRAVIS, QUAY, CODECOV}
	for _, url := range sources {
		body, err := statusFetch(url)
		if err != nil {
			panic(err)
		}

		var good bool
		switch url {
		case GITHUB:
			source := parsing.GithubStatus{}
			good = source.Parse(body)
		case CODECOV, TRAVIS, QUAY:
			source := parsing.StatusPageStatus{}
			good = source.Parse(body)
		}

		r := Refined{
			Good:          good,
			SourceMessage: body,
			LastUpdated:   time.Now(),
			Url:           url,
		}

		fmt.Println(string(*r.SourceMessage))
	}
}
