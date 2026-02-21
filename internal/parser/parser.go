// Package parser defines the DeviceParser interface and the central parser
// registry used by the validation pipeline.
package parser

import (
	"context"
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
)

// DeviceParser is the interface implemented by all vendor-specific parsers.
// Each parser converts raw configuration bytes into a structured ConfigModel.
type DeviceParser interface {
	// DeviceType returns the platform this parser handles.
	DeviceType() model.DeviceType
	// Parse converts raw configuration bytes into a ConfigModel.
	Parse(ctx context.Context, data []byte, device model.Device) (*model.ConfigModel, error)
}

// Parse is a convenience function that locates the appropriate parser from
// the default registry and parses the configuration.
func Parse(ctx context.Context, deviceType model.DeviceType, data []byte, device model.Device) (*model.ConfigModel, error) {
	p, ok := DefaultRegistry.Get(deviceType)
	if !ok {
		return nil, fmt.Errorf("parser: no parser registered for device type %q", deviceType)
	}
	return p.Parse(ctx, data, device)
}
