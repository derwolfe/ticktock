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
			InflightStatus.Inc()
			defer InflightStatus.Dec()
			defer wg.Done()
			body, err := statusFetch(url)
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

func requestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	metricsInit()
	updaterInit()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", status)
	http.ListenAndServe(":9090", requestLog(http.DefaultServeMux))
}
