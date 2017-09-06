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

func main() {
	store := state.NewStore()
	var wg sync.WaitGroup

	sources := []string{GITHUB, TRAVIS, QUAY, CODECOV}
	// match the length of sources
	wg.Add(4)
	for _, url := range sources {
		go func(url string) {
			defer wg.Done()
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

			r := state.Refined{
				Good:          good,
				SourceMessage: body,
				LastUpdated:   time.Now(),
				Url:           url,
			}
			store.Write(&r)
		}(url)
	}
	wg.Wait()
	fmt.Println("%#v", *store)
}
