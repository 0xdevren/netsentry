package cisco

import (
	"context"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
)

// NXOSParser parses Cisco NX-OS device configurations.
// NX-OS uses a similar line-based format to IOS with NX-OS-specific stanzas.
type NXOSParser struct {
	lexer *Lexer
}

// NewNXOSParser constructs an NXOSParser.
func NewNXOSParser() *NXOSParser {
	return &NXOSParser{lexer: NewLexer()}
}

// DeviceType returns the platform this parser handles.
func (p *NXOSParser) DeviceType() model.DeviceType {
	return model.DeviceTypeCiscoNXOS
}

// Parse converts raw NX-OS configuration into a ConfigModel.
// NX-OS shares the IOS line grammar, so the IOS parser handles the bulk of
// the work. NX-OS-specific directives are captured in GlobalSettings.
func (p *NXOSParser) Parse(ctx context.Context, data []byte, device model.Device) (*model.ConfigModel, error) {
	// Delegate to IOS parser for shared grammar.
	iosParser := NewIOSParser()
	cfg, err := iosParser.Parse(ctx, data, device)
	if err != nil {
		return nil, err
	}
	// Override device type.
	cfg.Device.Type = model.DeviceTypeCiscoNXOS

	// Capture NX-OS-specific top-level features.
	for _, line := range cfg.Lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "feature "):
			feature := strings.TrimPrefix(trimmed, "feature ")
			cfg.GlobalSettings["feature_"+feature] = "enabled"
		case strings.HasPrefix(trimmed, "vpc domain "):
			cfg.GlobalSettings["vpc_domain"] = strings.TrimPrefix(trimmed, "vpc domain ")
		case strings.HasPrefix(trimmed, "fabric forwarding anycast-gateway-mac "):
			cfg.GlobalSettings["anycast_gw_mac"] = strings.TrimPrefix(trimmed, "fabric forwarding anycast-gateway-mac ")
		case strings.HasPrefix(trimmed, "nv overlay evpn"):
			cfg.GlobalSettings["evpn"] = "enabled"
		}
	}

	return cfg, nil
}
