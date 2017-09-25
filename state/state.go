package state

import (
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
	Bodies      map[string]string
	LastUpdated time.Time
	Alarm       bool
}

func NewStore() *Store {
	return &Store{
		LastUpdated: time.Now(),
		Statuses:    make(map[string]*Refined),
		Bodies:      make(map[string]string),
		Alarm:       false,
	}
}

type Write interface {
	Write(status *Refined)
	Read() *Front
}

func (s *Store) Write(status *Refined) {
	s.Lock()
	defer s.Unlock()
	// reevaluate whether or not alarm state is bad and add a new status
	// message from travis. This should be parsed a bit better
	s.Statuses[status.Url] = status

	acc := true
	for _, r := range s.Statuses {
		acc = acc && r.Good
		s.Bodies[r.Url] = string(r.SourceMessage)
	}
	s.Alarm = !acc
	s.LastUpdated = time.Now()
}

type Front struct {
	Alarm       bool              `json:"alarm"`
	LastUpdated time.Time         `json:"last_updated"`
	Bodies      map[string]string `json:"status_messages"`
}

func (s Store) Read() *Front {
	s.RLock()
	ret := &Front{LastUpdated: s.LastUpdated, Alarm: s.Alarm, Bodies: s.Bodies}
	s.RUnlock()
	return ret
}
