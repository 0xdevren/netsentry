package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
)

// Detector analyses raw configuration text to identify the vendor/platform DeviceType.
type Detector struct{}

// NewDetector constructs a new Detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect examines the raw configuration byte slice and returns the most likely DeviceType.
func (d *Detector) Detect(data []byte) model.DeviceType {
	content := strings.ToLower(string(data))
	lines := bytes.Split(data, []byte("\n"))

	// Cisco NX-OS fingerprinting - check for NX-OS specific directives first.
	if strings.Contains(content, "nxos") ||
		strings.Contains(content, "feature nxapi") ||
		strings.Contains(content, "vpc domain") ||
		strings.Contains(content, "fabric forwarding") {
		return model.DeviceTypeCiscoNXOS
	}

	// Cisco IOS / IOS-XE fingerprinting.
	if d.hasAnyPrefix(lines,
		"version ",
		"service timestamps",
		"ip cef",
		"ip routing",
		"no ip domain-lookup",
	) || strings.Contains(content, "cisco ios") {
		return model.DeviceTypeCiscoIOS
	}

	// Juniper JunOS fingerprinting.
	if strings.Contains(content, "system {") ||
		strings.Contains(content, "interfaces {") ||
		strings.Contains(content, "protocols {") ||
		d.hasAnyPrefix(lines, "set system", "set interfaces", "set protocols") {
		return model.DeviceTypeJuniperOS
	}

	// Arista EOS fingerprinting.
	if strings.Contains(content, "arista") ||
		strings.Contains(content, "eos") ||
		strings.Contains(content, "management api http-commands") ||
		strings.Contains(content, "daemon terminattr") {
		return model.DeviceTypeAristaEOS
	}

	return model.DeviceTypeUnknown
}

// hasAnyPrefix reports whether any line in lines starts with any of the given prefixes.
func (d *Detector) hasAnyPrefix(lines [][]byte, prefixes ...string) bool {
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.ToLower(string(line)))
		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				return true
			}
		}
	}
	return false
}

// String returns a descriptive label for a DeviceType.
func DeviceTypeLabel(dt model.DeviceType) string {
	switch dt {
	case model.DeviceTypeCiscoIOS:
		return "Cisco IOS / IOS-XE"
	case model.DeviceTypeCiscoNXOS:
		return "Cisco NX-OS"
	case model.DeviceTypeJuniperOS:
		return "Juniper JunOS"
	case model.DeviceTypeAristaEOS:
		return "Arista EOS"
	default:
		return fmt.Sprintf("Unknown (%s)", string(dt))
	}
}
