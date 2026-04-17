// Package portclassify assigns a category to a port based on its number and label.
package portclassify

// Category represents a broad classification of a port's typical use.
type Category string

const (
	CategorySystem    Category = "system"    // 0–1023
	CategoryUser      Category = "user"      // 1024–49151
	CategoryDynamic   Category = "dynamic"   // 49152–65535
	CategoryDatabase  Category = "database"
	CategoryWeb       Category = "web"
	CategoryUnknown   Category = "unknown"
)

var webPorts = map[int]bool{80: true, 443: true, 8080: true, 8443: true, 3000: true, 5000: true}
var dbPorts = map[int]bool{3306: true, 5432: true, 6379: true, 27017: true, 1433: true, 5984: true}

// Classifier categorises ports.
type Classifier struct{}

// New returns a new Classifier.
func New() *Classifier { return &Classifier{} }

// Classify returns the Category for the given port number.
func (c *Classifier) Classify(port int) Category {
	if port < 0 || port > 65535 {
		return CategoryUnknown
	}
	if webPorts[port] {
		return CategoryWeb
	}
	if dbPorts[port] {
		return CategoryDatabase
	}
	switch {
	case port <= 1023:
		return CategorySystem
	case port <= 49151:
		return CategoryUser
	default:
		return CategoryDynamic
	}
}

// ClassifyAll returns a map of port → Category for each port in the slice.
func (c *Classifier) ClassifyAll(ports []int) map[int]Category {
	out := make(map[int]Category, len(ports))
	for _, p := range ports {
		out[p] = c.Classify(p)
	}
	return out
}
