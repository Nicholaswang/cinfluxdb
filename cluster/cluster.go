package cluster

import (
	"fmt"
	"log"
	"sync"
)

type Service struct {
	relays map[string]Relay
	etcd   EtcdConfig
}

func New(config Config) (*Service, error) {
	s := new(Service)
	s.relays = make(map[string]Relay)

	for _, cfg := range config.UDPBackends {
		u, err := NewUDP(cfg)
		if err != nil {
			return nil, err
		}
		if s.relays[u.Name()] != nil {
			return nil, fmt.Errorf("duplicate relay: %q", u.Name())
		}
		s.relays[u.Name()] = u
	}

	return s, nil
}

func (s *Service) Run() {
	var wg sync.WaitGroup
	wg.Add(len(s.relays))

	for k := range s.relays {
		relay := s.relays[k]
		go func() {
			defer wg.Done()

			if err := relay.Run(); err != nil {
				log.Printf("Error running relay %q: %v", relay.Name(), err)
			}
		}()
	}

	wg.Wait()
}

func (s *Service) Stop() {
	for _, v := range s.relays {
		v.Stop()
	}
}

type Relay interface {
	Name() string
	Run() error
	Stop() error
	Status() string //get status from etcd
}
