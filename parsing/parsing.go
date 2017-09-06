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

type Parser interface {
	Parse(body *[]byte) bool
}

func (g *GithubStatus) Parse(body *[]byte) bool {
	parsed := GithubStatus{}
	json.Unmarshal(*body, &parsed)
	return parsed.Status == "good"
}

func (g *StatusPageStatus) Parse(body *[]byte) bool {
	parsed := StatusPageStatus{
		Page:   StatusPageInnerPage{},
		Status: StatusPageInnerStatus{},
	}
	json.Unmarshal(*body, &parsed)
	return parsed.Status.Indicator == "none"
}
