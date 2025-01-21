package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota // select randomly
	RoundRobinSelect                   // select using Robbin algorithm
)

type Discovery interface {
	Refresh() error // refresh from remote registry
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}

// MultiServerDiscovery is a discovery for multi servers without a registry center
// it will try each server address in a round-robin manner
type MultiServerDiscovery struct {
	r       *rand.Rand   // random number generator
	mu      sync.RWMutex // lock
	servers []string     // server addresses
	index   int          // record the selected position for robin algorithm
}

// NewMultiServerDiscovery creates a MultiServerDiscovery instance
// @param servers: server addresses
// @return: *MultiServerDiscovery
// @note: servers is a list of server addresses
func NewMultiServerDiscovery(servers []string) *MultiServerDiscovery {
	d := &MultiServerDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

var _ Discovery = (*MultiServerDiscovery)(nil)

// Refresh doesn't make sense for MultiServerDiscovery
// @return: error
// @note: MultiServerDiscovery doesn't need to refresh
func (m *MultiServerDiscovery) Refresh() error {
	return nil
}

// Update updates the servers of discovery
func (m *MultiServerDiscovery) Update(servers []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.servers = servers
	return nil
}

// Get a server according to mode
// @param mode: select mode
// @return: server address, error
func (m *MultiServerDiscovery) Get(mode SelectMode) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := len(m.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case RandomSelect:
		return m.servers[m.r.Intn(n)], nil
	case RoundRobinSelect:
		s := m.servers[m.index%n]
		m.index = (m.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: invalid select mode")
	}
}

// GetAll returns all server addresses
// @return: server addresses, error
func (m *MultiServerDiscovery) GetAll() ([]string, error) {
	// read lock
	m.mu.RLock()
	defer m.mu.RUnlock()
	// return a copy of servers
	servers := make([]string, len(m.servers), len(m.servers))
	copy(servers, m.servers)
	return servers, nil
}


