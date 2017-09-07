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
	SourceMessage *[]byte
}

type Store struct {
	LastUpdated time.Time
	Statuses    map[string]*Refined
	mutex       *sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		LastUpdated: time.Now(),
		Statuses:    make(map[string]*Refined),
		mutex:       new(sync.RWMutex),
	}
}

type ReadWrite interface {
	Read(url string) (*Refined, error)
	Write(newStore *Refined)
}

func (s *Store) Read(url string) (*Refined, error) {
	s.mutex.RLock()
	target, ok := s.Statuses[url]
	s.mutex.RUnlock()
	if ok == false {
		return nil, errors.New("URL not found")
	}
	return target, nil
}

func (s *Store) Write(newStore *Refined) {
	s.mutex.Lock()
	s.Statuses[newStore.Url] = newStore
	s.LastUpdated = time.Now()
	s.mutex.Unlock()
}

type Front struct {
	Alarm       bool      `json:"alarm"`
	LastUpdated time.Time `json:"last_updated"`
}

func (s *Store) CurrentValue() *Front {
	s.mutex.RLock()
	acc := true
	for _, r := range s.Statuses {
		acc = acc && r.Good
	}
	ret := &Front{Alarm: acc, LastUpdated: s.LastUpdated}
	s.mutex.RUnlock()
	return ret
}
