package state

import (
	"github.com/derwolfe/ticktock/parsing"
	"sync"
)

type State struct {
	statuses map[string]*parsing.Unified
	mutex    sync.Mutex
}

type ReadWrite interface {
	Read(url string) (*parsing.Unified, error)
	Write(url string, newState *parsing.Unified)
}

func (s *State) Read(url string) (*parsing.Unified, error) {
	s.RWLock()
	state, err := s[url]
	if err != nil {
		return nil, err
	}
	s.RWUnlock()
	return state, nil
}

func (s *State) Write(url string, newState *parsing.Unified) {
	s.Lock()
	s[url] = newState
	s.Unlock()
}
