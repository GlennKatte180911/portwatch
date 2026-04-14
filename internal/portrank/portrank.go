// Package portrank assigns a severity rank to ports based on how
// commonly they are targeted or how sensitive their services are.
// Higher ranks indicate ports that warrant closer attention when
// they appear in a diff.
package portrank

import "sync"

// Rank represents the sensitivity level of a port.
type Rank int

const (
	RankLow    Rank = 1
	RankMedium Rank = 2
	RankHigh   Rank = 3
)

// String returns a human-readable label for the rank.
func (r Rank) String() string {
	switch r {
	case RankHigh:
		return "high"
	case RankMedium:
		return "medium"
	default:
		return "low"
	}
}

// Ranker maps ports to sensitivity ranks.
type Ranker struct {
	mu    sync.RWMutex
	ranks map[int]Rank
}

// defaultRanks contains well-known ports and their default ranks.
var defaultRanks = map[int]Rank{
	21:   RankHigh,   // FTP
	22:   RankHigh,   // SSH
	23:   RankHigh,   // Telnet
	25:   RankMedium, // SMTP
	80:   RankLow,    // HTTP
	443:  RankLow,    // HTTPS
	3306: RankHigh,   // MySQL
	5432: RankHigh,   // PostgreSQL
	6379: RankHigh,   // Redis
	8080: RankLow,    // HTTP alt
	27017: RankHigh,  // MongoDB
}

// New returns a Ranker pre-loaded with default port ranks.
func New() *Ranker {
	r := &Ranker{ranks: make(map[int]Rank, len(defaultRanks))}
	for port, rank := range defaultRanks {
		r.ranks[port] = rank
	}
	return r
}

// Set assigns a rank to a port, overwriting any existing value.
func (r *Ranker) Set(port int, rank Rank) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ranks[port] = rank
}

// Get returns the rank for a port. If the port is not registered,
// RankLow is returned as the default.
func (r *Ranker) Get(port int) Rank {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rank, ok := r.ranks[port]; ok {
		return rank
	}
	return RankLow
}

// Annotate returns a map of port -> Rank for every port in the slice.
func (r *Ranker) Annotate(ports []int) map[int]Rank {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int]Rank, len(ports))
	for _, p := range ports {
		if rank, ok := r.ranks[p]; ok {
			out[p] = rank
		} else {
			out[p] = RankLow
		}
	}
	return out
}

// MaxRank returns the highest rank found among the given ports.
// Returns RankLow when the slice is empty.
func (r *Ranker) MaxRank(ports []int) Rank {
	max := RankLow
	for port, rank := range r.Annotate(ports) {
		_ = port
		if rank > max {
			max = rank
		}
	}
	return max
}
