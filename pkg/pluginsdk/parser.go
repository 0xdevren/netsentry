// Package pluginsdk defines the public plugin interface for NetSentry parser plugins.
// Third-party developers implement these interfaces to add new device format support.
package pluginsdk

import (
	"context"
	"github.com/0xdevren/netsentry/internal/model"
)

// ParserPlugin is implemented by external parser plugins that add support for
// new device configuration formats. Register with the parser registry using
// parser.DefaultRegistry.Register(plugin).
type ParserPlugin interface {
	// DeviceType returns the platform identifier this parser handles.
	DeviceType() model.DeviceType
	// Parse converts raw configuration bytes into a structured ConfigModel.
	Parse(ctx context.Context, data []byte, device model.Device) (*model.ConfigModel, error)
}
