// Package arista provides parsers for Arista EOS device configurations.
package arista

import (
	"context"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/parser/cisco"
)

// EOSParser parses Arista EOS device configurations.
// EOS uses a Cisco IOS-compatible CLI grammar with EOS-specific extensions.
type EOSParser struct {
	ios *cisco.IOSParser
}

// NewEOSParser constructs an EOSParser.
func NewEOSParser() *EOSParser {
	return &EOSParser{ios: cisco.NewIOSParser()}
}

// DeviceType returns the platform this parser handles.
func (p *EOSParser) DeviceType() model.DeviceType { return model.DeviceTypeAristaEOS }

// Parse converts Arista EOS configuration into a ConfigModel.
// EOS shares the IOS grammar, so the IOS parser handles the structural
// parsing. EOS-specific directives (MLAG, EVPN, management API) are
// added to GlobalSettings.
func (p *EOSParser) Parse(ctx context.Context, data []byte, device model.Device) (*model.ConfigModel, error) {
	cfg, err := p.ios.Parse(ctx, data, device)
	if err != nil {
		return nil, err
	}
	cfg.Device.Type = model.DeviceTypeAristaEOS

	for _, line := range cfg.Lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "management api http-commands"):
			cfg.GlobalSettings["management_api"] = "http-commands"
		case strings.HasPrefix(trimmed, "daemon terminattr"):
			cfg.GlobalSettings["terminattr"] = "enabled"
		case strings.HasPrefix(trimmed, "mlag configuration"):
			cfg.GlobalSettings["mlag"] = "configured"
		case strings.HasPrefix(trimmed, "router bgp") && cfg.BGPConfig != nil:
			// Already handled by IOS parser.
		case strings.HasPrefix(trimmed, "ip virtual-router mac-address"):
			cfg.GlobalSettings["virtual_router_mac"] = strings.TrimPrefix(trimmed, "ip virtual-router mac-address ")
		case strings.HasPrefix(trimmed, "vxlan vni"):
			cfg.GlobalSettings["vxlan"] = "configured"
		}
	}

	return cfg, nil
}
