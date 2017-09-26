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
	TIMEOUT = 3 * time.Second
	GITHUB  = "https://status.github.com/api/status.json"
	CODECOV = "https://wdzsn5dlywj9.statuspage.io/api/v2/status.json"
	TRAVIS  = "https://pnpcptp8xh9k.statuspage.io/api/v2/status.json"
	QUAY    = "https://8szqd6w4s277.statuspage.io/api/v2/status.json"
)

var (
	DataStore      = state.NewStore()
	InflightStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "ticktock",
			Subsystem: "status_checks",
			Name:      "in_flight",
			Help:      "Number of in flight status checks.",
		})
)

func metricsInit() {
	prometheus.MustRegister(InflightStatus)
}

func updaterInit() {
	ticker := time.NewTicker(1 * time.Minute)
	updateState(DataStore)
	go func() {
		for _ = range ticker.C {
			updateState(DataStore)
		}
	}()
}

func statusFetch(url string, client http.Client) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func updateState(store *state.Store) {
	wg := sync.WaitGroup{}

	sources := []string{GITHUB, TRAVIS, QUAY, CODECOV}
	// match the length of sources
	wg.Add(len(sources))

	timeout := time.Duration(TIMEOUT)
	client := http.Client{
		Timeout: timeout,
	}

	for _, url := range sources {
		go func(url string) {
			log.Printf("Started Fetching: %s", url)
			InflightStatus.Inc()
			defer InflightStatus.Dec()
			defer wg.Done()
			body, err := statusFetch(url, client)
			if err != nil {
				log.Printf("Error fetching: %s, %s", url, err)
				return
			}
			log.Printf("Succeeded fetching: %s", url)

			var parser parsing.Parser
			switch url {
			case GITHUB:
				parser = parsing.GithubParser
			case CODECOV, TRAVIS, QUAY:
				parser = parsing.StatusPageParser
			default:
				parser = parsing.DefaultParser
			}
			r := parser(body)
			store.Write(&r)
		}(url)
	}
	wg.Wait()
}

func status(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(DataStore.Read())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func requestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	metricsInit()
	updaterInit()

	http.Handle("/", http.FileServer(http.Dir("/src/github.com/derwolfe/ticktock/static/")))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api", status)
	http.ListenAndServe(":9090", requestLog(http.DefaultServeMux))
}
