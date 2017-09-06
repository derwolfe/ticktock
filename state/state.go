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
	statuses map[string]*Refined
	mutex    *sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		statuses: make(map[string]*Refined),
		mutex:    new(sync.RWMutex),
	}
}

type ReadWrite interface {
	Read(url string) (*Refined, error)
	Write(newStore *Refined)
}

func (s *Store) Read(url string) (*Refined, error) {
	s.mutex.RLock()
	target, ok := s.statuses[url]
	s.mutex.RUnlock()
	if ok == false {
		return nil, errors.New("URL not found")
	}
	return target, nil
}

func (s *Store) Write(newStore *Refined) {
	s.mutex.Lock()
	s.statuses[newStore.Url] = newStore
	s.mutex.Unlock()
}
