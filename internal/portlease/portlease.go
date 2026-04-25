// Package portlease tracks temporary port leases with expiry times,
// allowing short-lived port openings to be distinguished from persistent ones.
package portlease

import (
	"fmt"
	"sync"
	"time"
)

// Lease represents a time-bounded claim on a port.
type Lease struct {
	Port      int
	GrantedAt time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the lease has passed its expiry time.
func (l Lease) IsExpired(now time.Time) bool {
	return now.After(l.ExpiresAt)
}

// Registry manages active port leases.
type Registry struct {
	mu     sync.Mutex
	leases map[int]Lease
	now    func() time.Time
}

// New returns an initialised lease Registry.
func New() *Registry {
	return &Registry{
		leases: make(map[int]Lease),
		now:    time.Now,
	}
}

// Grant creates or renews a lease for port with the given duration.
// Returns an error if port is out of range or duration is non-positive.
func (r *Registry) Grant(port int, duration time.Duration) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portlease: port %d out of range", port)
	}
	if duration <= 0 {
		return fmt.Errorf("portlease: duration must be positive")
	}
	now := r.now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.leases[port] = Lease{
		Port:      port,
		GrantedAt: now,
		ExpiresAt: now.Add(duration),
	}
	return nil
}

// Revoke removes the lease for port if it exists.
func (r *Registry) Revoke(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.leases, port)
}

// Active returns true if port has a current, non-expired lease.
func (r *Registry) Active(port int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	l, ok := r.leases[port]
	if !ok {
		return false
	}
	return !l.IsExpired(r.now())
}

// Expired returns all leases that have passed their expiry time.
func (r *Registry) Expired() []Lease {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	var out []Lease
	for _, l := range r.leases {
		if l.IsExpired(now) {
			out = append(out, l)
		}
	}
	return out
}

// Purge removes all expired leases and returns the number removed.
func (r *Registry) Purge() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	count := 0
	for port, l := range r.leases {
		if l.IsExpired(now) {
			delete(r.leases, port)
			count++
		}
	}
	return count
}
