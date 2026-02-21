package parser

import (
	"fmt"
	"sync"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser/arista"
	"github.com/0xdevren/netsentry/internal/parser/cisco"
	"github.com/0xdevren/netsentry/internal/parser/juniper"
)

// Registry maintains the mapping from DeviceType to DeviceParser.
type Registry struct {
	mu      sync.RWMutex
	parsers map[model.DeviceType]DeviceParser
}

// NewRegistry constructs an empty Registry.
func NewRegistry() *Registry {
	return &Registry{parsers: make(map[model.DeviceType]DeviceParser)}
}

// Register adds a parser to the registry.
func (r *Registry) Register(p DeviceParser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[p.DeviceType()] = p
}

// Get retrieves the parser for the given device type.
func (r *Registry) Get(dt model.DeviceType) (DeviceParser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.parsers[dt]
	return p, ok
}

// MustGet retrieves the parser or panics if not found.
func (r *Registry) MustGet(dt model.DeviceType) DeviceParser {
	p, ok := r.Get(dt)
	if !ok {
		panic(fmt.Sprintf("parser registry: no parser for %q", dt))
	}
	return p
}

// DefaultRegistry is the global registry pre-populated with all built-in parsers.
var DefaultRegistry = func() *Registry {
	r := NewRegistry()
	r.Register(cisco.NewIOSParser())
	r.Register(cisco.NewNXOSParser())
	r.Register(juniper.NewJunOSParser())
	r.Register(arista.NewEOSParser())
	return r
}()
