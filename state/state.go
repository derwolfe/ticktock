package state

import (
	"errors"
	"sync"
	"time"
)

type Refined struct {
	Url           string
	LastUpdated   time.Time
	Good          bool
	SourceMessage []byte
}

type Store struct {
	sync.RWMutex
	Statuses    map[string]*Refined
	LastUpdated time.Time
	Alarm       bool
}

func NewStore() *Store {
	return &Store{
		LastUpdated: time.Now(),
		Statuses:    make(map[string]*Refined),
	}
}

type ReadWrite interface {
	Read(url string) (*Refined, error)
	Write(newStore *Refined)
	CurrentValue() *Front
}

func (s Store) Read(url string) (*Refined, error) {
	s.RLock()
	target, ok := s.Statuses[url]
	s.RUnlock()
	if ok == false {
		return nil, errors.New("URL not found")
	}
	return target, nil
}

func (s *Store) Write(status *Refined) {
	s.Lock()
	s.Statuses[status.Url] = status
	acc := true
	for _, r := range s.Statuses {
		acc = acc && r.Good
	}

	s.Alarm = !acc
	s.LastUpdated = time.Now()
	s.Unlock()
}

type Front struct {
	Alarm       bool      `json:"alarm"`
	LastUpdated time.Time `json:"last_updated"`
}

func (s Store) CurrentValue() *Front {
	s.RLock()
	ret := &Front{LastUpdated: s.LastUpdated, Alarm: s.Alarm}
	s.RUnlock()
	return ret
}
