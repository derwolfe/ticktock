package main

import (
	"encoding/json"
	"github.com/derwolfe/ticktock/parsing"
	"github.com/derwolfe/ticktock/state"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
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

// GLOBAL :barf:
var (
	DataStore      = state.NewStore()
	inflightStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ticktock",
			Subsystem: "status_checks",
			Name:      "in_flight",
			Help:      "Number of in flight status checks.",
		})
)

func metricsInit() {
	prometheus.MustRegister(inflightStatus)
}

func updaterInit() {
	ticker := time.NewTicker(1 * time.Minute)
	updateState(DataStore)
	go func() {
		for {
			select {
			case <-ticker.C:
				updateState(DataStore)
			}
		}
	}()
}

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

func updateState(store *state.Store) {
	var wg sync.WaitGroup
	sources := []string{GITHUB, TRAVIS, QUAY, CODECOV}
	// match the length of sources
	wg.Add(4)

	for _, url := range sources {
		go func(url string) {
			log.Printf("Started Fetching: %s", url)
			inflightStatus.Inc()
			defer inflightStatus.Dec()
			defer wg.Done()
			body, err := statusFetch(url)
			// bail out after a few attempts if we've encountered a few errors
			if err != nil {
				log.Printf("Error fetching: %s, %s", url, err)
				return
			}
			log.Printf("Succeeded fetching: %s", url)

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
}

func status(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(DataStore.CurrentValue())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	metricsInit()
	updaterInit()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", status) // set router
	http.ListenAndServe(":9090", nil)
}
