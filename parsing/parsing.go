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

// Should parser copy memory instead of allocating constantly?
type Parser interface {
	Parse(body []byte) (Unified, error)
}

func (g *GithubStatus) Parse(body []byte) (*Unified, error) {
	res := GithubStatus{}
	json.Unmarshal(body, &res)

	retval := Unified{}
	if res.Status != "good" {
		retval.Good = false
	} else {
		retval.Good = true
	}
	retval.UpdatedAt = time.Now()
	return &retval, nil
}

func (g *StatusPageStatus) Parse(body []byte) (*Unified, error) {
	res := StatusPageStatus{
		Page:   StatusPageInnerPage{},
		Status: StatusPageInnerStatus{},
	}
	json.Unmarshal(body, &res)

	retval := Unified{}
	if res.Status.Inidicator != "None" {
		retval.Good = false
	} else {
		retval.Good = true
	}
	retval.UpdatedAt = time.Now()
	return &retval, nil
}
