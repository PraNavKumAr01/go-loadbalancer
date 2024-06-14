package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	Mutex        sync.RWMutex
}

type LoadBalancer struct {
	backends []*Backend
	current  uint32
}

func (lb *LoadBalancer) NextIndex() int {
	return int(atomic.AddUint32(&lb.current, uint32(1)) % uint32(len(lb.backends)))
}

func (lb *LoadBalancer) GetNextPeer() *Backend {
	next := lb.NextIndex()
	backend := lb.backends[next]
	if backend.Alive {
		return backend
	}
	return nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := lb.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}

func (lb *LoadBalancer) HealthCheck() {
	for {
		for _, backend := range lb.backends {
			go func(b *Backend) {
				resp, err := http.Get(b.URL.String())
				b.Mutex.Lock()
				if err != nil || resp.StatusCode != http.StatusOK {
					b.Alive = false
				} else {
					b.Alive = true
				}
				b.Mutex.Unlock()
			}(backend)
		}
		time.Sleep(time.Second * 10)
	}
}

func NewLoadBalancer(backendURLs []string) *LoadBalancer {
	var backends []*Backend
	for _, backendURL := range backendURLs {
		url, err := url.Parse(backendURL)
		if err != nil {
			log.Fatalf("Could not parse URL: %v", err)
		}
		backend := &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}
		backends = append(backends, backend)
	}
	lb := &LoadBalancer{
		backends: backends,
	}
	go lb.HealthCheck()
	return lb
}
