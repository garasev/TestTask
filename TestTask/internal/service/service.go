package ser

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const PINGURL = "/api/ping"

type Service struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func NewService(url *url.URL, proxy *httputil.ReverseProxy) *Service {
	return &Service{
		URL:          url,
		Alive:        true,
		ReverseProxy: proxy,
	}
}

func (s *Service) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

func (s *Service) IsAlive() (alive bool) {
	s.mux.Lock()
	alive = s.Alive
	s.mux.Unlock()
	return
}

type ServicePool struct {
	services []*Service
	cur      int
}

func NewServicePool() *ServicePool {
	return &ServicePool{
		services: make([]*Service, 0),
		cur:      0,
	}
}

func (sp *ServicePool) Add(service *Service) {
	sp.services = append(sp.services, service)
}

func (sp *ServicePool) HealthCheck() {
	for _, service := range sp.services {
		alive := isBackendAlive(service.URL)
		service.SetAlive(alive)
	}
}

func (sp *ServicePool) String() string {
	res := ""
	for _, service := range sp.services {
		res += service.URL.Host + fmt.Sprintf(": %t;", service.Alive)
	}

	return res
}

func (sp *ServicePool) NextIndex() int {
	sp.cur += 1
	return (sp.cur) % (len(sp.services))
}

func (s *ServicePool) GetNextPeer() *Service {
	next := s.NextIndex()
	l := len(s.services) + next
	for i := next; i < l; i++ {
		idx := i % len(s.services)
		if s.services[idx].IsAlive() {
			if i != next {
				s.cur = idx
			}
			return s.services[idx]
		}
	}
	return nil
}

func isBackendAlive(u *url.URL) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(u.String() + PINGURL)
	if err != nil {
		//fmt.Println(err.Error())
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func (sp *ServicePool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, s := range sp.services {
		if s.URL.String() == backendUrl.String() {
			s.SetAlive(alive)
			break
		}
	}
}
