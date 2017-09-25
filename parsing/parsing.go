package parsing

import (
	"encoding/json"
	"github.com/derwolfe/ticktock/state"
	"time"
)

type StatusPageInnerPage struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StatusPageInnerStatus struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description`
}

type StatusPageStatus struct {
	Page   StatusPageInnerPage   `json:"page"`
	Status StatusPageInnerStatus `json:"status"`
}

type GithubStatus struct {
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
}

type Parser func([]byte) state.Refined

func GithubParser(body []byte) state.Refined {
	parsed := GithubStatus{}
	json.Unmarshal(body, &parsed)
	good := parsed.Status == "good"
	return state.Refined{
		Url:           "https://status.github.com/api/status.json",
		LastUpdated:   parsed.LastUpdated,
		ServiceName:   "Github",
		SourceMessage: parsed.Status,
		Good:          good,
	}
}

func StatusPageParser(body []byte) state.Refined {
	parsed := StatusPageStatus{
		Page:   StatusPageInnerPage{},
		Status: StatusPageInnerStatus{},
	}
	json.Unmarshal(body, &parsed)
	good := parsed.Status.Indicator == "none"
	return state.Refined{
		Url:           parsed.Page.Url,
		LastUpdated:   parsed.Page.UpdatedAt,
		ServiceName:   parsed.Page.Name,
		SourceMessage: parsed.Status.Description,
		Good:          good,
	}
}

func DefaultParser(body []byte) state.Refined {
	return state.Refined{}
}
