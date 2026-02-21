package cisco

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
)

// IOSParser parses Cisco IOS and IOS-XE device configurations.
type IOSParser struct {
	lexer *Lexer
}

// NewIOSParser constructs an IOSParser.
func NewIOSParser() *IOSParser {
	return &IOSParser{lexer: NewLexer()}
}

// DeviceType returns the platform this parser handles.
func (p *IOSParser) DeviceType() model.DeviceType {
	return model.DeviceTypeCiscoIOS
}

// Parse converts raw Cisco IOS configuration bytes into a ConfigModel.
func (p *IOSParser) Parse(_ context.Context, data []byte, device model.Device) (*model.ConfigModel, error) {
	lines := splitLines(data)
	cfg := &model.ConfigModel{
		Device:         device,
		RawText:        string(data),
		Lines:          lines,
		GlobalSettings: make(map[string]string),
	}

	tokens := p.lexer.Tokenise(data)

	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		text := tok.Text

		switch {
		case strings.HasPrefix(text, "hostname "):
			cfg.Device.Hostname = strings.TrimPrefix(text, "hostname ")
			cfg.GlobalSettings["hostname"] = cfg.Device.Hostname

		case strings.HasPrefix(text, "interface "):
			iface, consumed := p.parseInterface(tokens, i)
			cfg.Interfaces = append(cfg.Interfaces, iface)
			i += consumed
			continue

		case strings.HasPrefix(text, "ip access-list "):
			acl, consumed := p.parseACL(tokens, i)
			cfg.ACLs = append(cfg.ACLs, acl)
			i += consumed
			continue

		case strings.HasPrefix(text, "router bgp "):
			bgpCfg, consumed := p.parseBGP(tokens, i)
			cfg.BGPConfig = bgpCfg
			i += consumed
			continue

		case strings.HasPrefix(text, "router ospf "):
			ospfCfg, consumed := p.parseOSPF(tokens, i)
			cfg.OSPFConfig = ospfCfg
			i += consumed
			continue

		case strings.HasPrefix(text, "ip route "):
			route := p.parseStaticRoute(text)
			cfg.StaticRoutes = append(cfg.StaticRoutes, route)

		case strings.HasPrefix(text, "vlan "):
			vlan := p.parseVLAN(text)
			cfg.VLANs = append(cfg.VLANs, vlan)

		case strings.HasPrefix(text, "logging "):
			cfg.GlobalSettings["logging"] = strings.TrimPrefix(text, "logging ")

		case strings.HasPrefix(text, "ntp server "):
			cfg.GlobalSettings["ntp_server"] = strings.TrimPrefix(text, "ntp server ")

		case text == "no ip domain-lookup":
			cfg.GlobalSettings["no_domain_lookup"] = "true"

		case strings.HasPrefix(text, "enable secret ") || strings.HasPrefix(text, "enable password "):
			cfg.GlobalSettings["enable_secret"] = "configured"

		case strings.HasPrefix(text, "snmp-server "):
			cfg.GlobalSettings["snmp_server_"+fmt.Sprintf("%d", i)] = text
		}

		i++
	}

	return cfg, nil
}

// parseInterface extracts an interface block beginning at index i.
// Returns the populated Interface and the number of tokens consumed.
func (p *IOSParser) parseInterface(tokens []Token, start int) (model.Interface, int) {
	iface := model.Interface{
		Name:       strings.TrimPrefix(tokens[start].Text, "interface "),
		Attributes: make(map[string]string),
	}
	consumed := 1
	baseDepth := tokens[start].Depth

	for i := start + 1; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.Depth <= baseDepth && tok.Type != TokenBlockStart {
			break
		}
		text := tok.Text
		switch {
		case strings.HasPrefix(text, "description "):
			iface.Description = strings.TrimPrefix(text, "description ")
		case strings.HasPrefix(text, "ip address "):
			parts := strings.Fields(strings.TrimPrefix(text, "ip address "))
			if len(parts) >= 2 {
				iface.IPAddress = parts[0]
				iface.SubnetMask = parts[1]
			}
		case strings.HasPrefix(text, "ipv6 address "):
			iface.IPv6Address = strings.TrimPrefix(text, "ipv6 address ")
		case text == "shutdown":
			iface.Shutdown = true
		case strings.HasPrefix(text, "switchport mode "):
			iface.VLANMode = strings.TrimPrefix(text, "switchport mode ")
		case strings.HasPrefix(text, "switchport access vlan "):
			if id, err := strconv.Atoi(strings.TrimPrefix(text, "switchport access vlan ")); err == nil {
				iface.AccessVLAN = id
			}
		case strings.HasPrefix(text, "mtu "):
			if m, err := strconv.Atoi(strings.TrimPrefix(text, "mtu ")); err == nil {
				iface.MTU = m
			}
		case strings.HasPrefix(text, "ip access-group "):
			parts := strings.Fields(strings.TrimPrefix(text, "ip access-group "))
			if len(parts) == 2 {
				if parts[1] == "in" {
					iface.InboundACL = parts[0]
				} else if parts[1] == "out" {
					iface.OutboundACL = parts[0]
				}
			}
		case text == "spanning-tree portfast":
			iface.SpanningTreePortFast = true
		default:
			iface.Attributes[fmt.Sprintf("line_%d", consumed)] = text
		}
		consumed++
	}

	return iface, consumed
}

// parseACL extracts an IP access-list block.
func (p *IOSParser) parseACL(tokens []Token, start int) (model.ACL, int) {
	header := strings.TrimPrefix(tokens[start].Text, "ip access-list ")
	parts := strings.Fields(header)
	acl := model.ACL{}
	if len(parts) >= 2 {
		acl.Type = parts[0]
		acl.Name = parts[1]
	}
	consumed := 1
	baseDepth := tokens[start].Depth

	for i := start + 1; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.Depth <= baseDepth && tok.Type != TokenBlockStart {
			break
		}
		entry := parseACLEntry(tok.Text)
		acl.Entries = append(acl.Entries, entry)
		consumed++
	}

	return acl, consumed
}

// parseACLEntry converts a single ACL entry line into an ACLEntry.
func parseACLEntry(text string) model.ACLEntry {
	entry := model.ACLEntry{}
	parts := strings.Fields(text)
	idx := 0

	// Optional sequence number.
	if idx < len(parts) {
		if n, err := strconv.Atoi(parts[idx]); err == nil {
			entry.Sequence = n
			idx++
		}
	}
	// Action.
	if idx < len(parts) {
		entry.Action = model.ACLAction(parts[idx])
		idx++
	}
	// Protocol.
	if idx < len(parts) {
		entry.Protocol = parts[idx]
		idx++
	}
	// Source.
	if idx < len(parts) {
		entry.Source = parts[idx]
		idx++
	}
	// Destination.
	if idx < len(parts) {
		entry.Destination = parts[idx]
		idx++
	}
	// Log keyword.
	for _, p := range parts[idx:] {
		if p == "log" {
			entry.Log = true
		}
	}
	return entry
}

// parseBGP extracts the BGP router block.
func (p *IOSParser) parseBGP(tokens []Token, start int) (*model.BGPConfig, int) {
	asStr := strings.TrimPrefix(tokens[start].Text, "router bgp ")
	bgp := &model.BGPConfig{}
	if as, err := strconv.Atoi(strings.TrimSpace(asStr)); err == nil {
		bgp.LocalAS = as
	}
	consumed := 1
	baseDepth := tokens[start].Depth

	var currentNeighbor *model.BGPNeighbor
	neighborMap := make(map[string]*model.BGPNeighbor)

	for i := start + 1; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.Depth <= baseDepth && tok.Type != TokenBlockStart {
			break
		}
		text := tok.Text
		switch {
		case strings.HasPrefix(text, "bgp router-id "):
			bgp.RouterID = strings.TrimPrefix(text, "bgp router-id ")
		case strings.HasPrefix(text, "neighbor "):
			parts := strings.Fields(text)
			if len(parts) >= 3 {
				addr := parts[1]
				if _, ok := neighborMap[addr]; !ok {
					n := &model.BGPNeighbor{Address: addr}
					neighborMap[addr] = n
					bgp.Neighbors = append(bgp.Neighbors, *n)
				}
				currentNeighbor = neighborMap[addr]
				attr := strings.Join(parts[2:], " ")
				switch {
				case strings.HasPrefix(attr, "remote-as "):
					if as, err := strconv.Atoi(strings.TrimPrefix(attr, "remote-as ")); err == nil {
						currentNeighbor.RemoteAS = as
					}
				case strings.HasPrefix(attr, "description "):
					currentNeighbor.Description = strings.TrimPrefix(attr, "description ")
				case attr == "next-hop-self":
					currentNeighbor.NextHopSelf = true
				case attr == "shutdown":
					currentNeighbor.Shutdown = true
				case strings.HasPrefix(attr, "update-source "):
					currentNeighbor.UpdateSource = strings.TrimPrefix(attr, "update-source ")
				case strings.HasPrefix(attr, "route-map "):
					rparts := strings.Fields(attr)
					if len(rparts) == 3 {
						if rparts[2] == "in" {
							currentNeighbor.RouteMapIn = rparts[1]
						} else {
							currentNeighbor.RouteMapOut = rparts[1]
						}
					}
				}
				// Sync back.
				_ = currentNeighbor
			}
		case strings.HasPrefix(text, "network "):
			parts := strings.Fields(text)
			net := model.BGPNetwork{}
			if len(parts) >= 2 {
				net.Prefix = parts[1]
			}
			if len(parts) == 4 && parts[2] == "mask" {
				net.Mask = parts[3]
			}
			bgp.Networks = append(bgp.Networks, net)
		}
		consumed++
	}

	// Re-sync neighbors from map.
	updated := make([]model.BGPNeighbor, 0, len(neighborMap))
	for _, n := range neighborMap {
		updated = append(updated, *n)
	}
	bgp.Neighbors = updated

	return bgp, consumed
}

// parseOSPF extracts the OSPF router block.
func (p *IOSParser) parseOSPF(tokens []Token, start int) (*model.OSPFConfig, int) {
	pidStr := strings.TrimPrefix(tokens[start].Text, "router ospf ")
	ospf := &model.OSPFConfig{}
	if pid, err := strconv.Atoi(strings.TrimSpace(pidStr)); err == nil {
		ospf.ProcessID = pid
	}
	consumed := 1
	baseDepth := tokens[start].Depth

	for i := start + 1; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.Depth <= baseDepth && tok.Type != TokenBlockStart {
			break
		}
		text := tok.Text
		switch {
		case strings.HasPrefix(text, "router-id "):
			ospf.RouterID = strings.TrimPrefix(text, "router-id ")
		case strings.HasPrefix(text, "network "):
			parts := strings.Fields(text)
			if len(parts) >= 4 && parts[len(parts)-2] == "area" {
				areaID := parts[len(parts)-1]
				network := parts[1]
				added := false
				for idx, a := range ospf.Areas {
					if a.ID == areaID {
						ospf.Areas[idx].Networks = append(ospf.Areas[idx].Networks, network)
						added = true
						break
					}
				}
				if !added {
					ospf.Areas = append(ospf.Areas, model.OSPFArea{ID: areaID, Networks: []string{network}})
				}
			}
		case strings.HasPrefix(text, "passive-interface default"):
			ospf.DefaultPassive = true
		case strings.HasPrefix(text, "passive-interface "):
			ospf.PassiveInterfaces = append(ospf.PassiveInterfaces, strings.TrimPrefix(text, "passive-interface "))
		case strings.HasPrefix(text, "redistribute "):
			parts := strings.Fields(text)
			if len(parts) >= 2 {
				ospf.Redistributions = append(ospf.Redistributions, model.OSPFRedistribution{Source: parts[1]})
			}
		}
		consumed++
	}

	return ospf, consumed
}

// parseStaticRoute parses an "ip route" line into a StaticRoute.
func (p *IOSParser) parseStaticRoute(text string) model.StaticRoute {
	parts := strings.Fields(strings.TrimPrefix(text, "ip route "))
	route := model.StaticRoute{}
	if len(parts) >= 3 {
		route.Destination = parts[0] + "/" + maskToPrefix(parts[1])
		route.NextHop = parts[2]
	}
	if len(parts) >= 4 {
		if ad, err := strconv.Atoi(parts[3]); err == nil {
			route.AdminDistance = ad
		}
	}
	return route
}

// parseVLAN parses a "vlan <id>" line into a VLAN.
func (p *IOSParser) parseVLAN(text string) model.VLAN {
	idStr := strings.TrimPrefix(text, "vlan ")
	vlan := model.VLAN{}
	if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
		vlan.ID = id
		vlan.State = "active"
	}
	return vlan
}

// splitLines splits raw config bytes on newlines and returns trimmed non-empty lines.
func splitLines(data []byte) []string {
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimRight(line, "\r")
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

// maskToPrefix converts a dotted subnet mask to a CIDR prefix string.
func maskToPrefix(mask string) string {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return mask
	}
	count := 0
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return mask
		}
		for n > 0 {
			count += n & 1
			n >>= 1
		}
	}
	return strconv.Itoa(count)
}
