package parsing

import (
	"encoding/json"
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

type Unified struct {
	Good      bool
	UpdatedAt time.Time
}

type Parser interface {
	Parse(body []byte) *Unified
}

func (g *GithubStatus) Parse(body []byte) *Unified {
	parsed := GithubStatus{}
	json.Unmarshal(body, &parsed)

	retval := Unified{
		UpdatedAt: time.Now(),
	}
	if parsed.Status != "good" {
		retval.Good = false
	} else {
		retval.Good = true
	}
	return &retval
}

func (g *StatusPageStatus) Parse(body []byte) *Unified {
	parsed := StatusPageStatus{
		Page:   StatusPageInnerPage{},
		Status: StatusPageInnerStatus{},
	}
	json.Unmarshal(body, &parsed)

	retval := Unified{
		UpdatedAt: time.Now(),
	}
	if parsed.Status.Indicator != "none" {
		retval.Good = false
	} else {
		retval.Good = true
	}
	return &retval
}
