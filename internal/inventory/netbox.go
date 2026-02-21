package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/0xdevren/netsentry/internal/model"
)

// NetBoxOptions configures the NetBox inventory provider.
type NetBoxOptions struct {
	// BaseURL is the NetBox API base URL (e.g. "https://netbox.example.com").
	BaseURL string
	// Token is the NetBox API authentication token.
	Token string
	// Timeout is the HTTP request timeout.
	Timeout time.Duration
}

// NetBoxInventory fetches the device inventory from a NetBox instance.
type NetBoxInventory struct {
	opts   NetBoxOptions
	client *http.Client
}

// NewNetBoxInventory constructs a NetBoxInventory with the given options.
func NewNetBoxInventory(opts NetBoxOptions) *NetBoxInventory {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &NetBoxInventory{
		opts:   opts,
		client: &http.Client{Timeout: timeout},
	}
}

// netBoxDevice is the JSON shape returned by the NetBox /dcim/devices/ API.
type netBoxDevice struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Platform struct {
		Slug string `json:"slug"`
	} `json:"platform"`
	PrimaryIP struct {
		Address string `json:"address"`
	} `json:"primary_ip4"`
	Site struct {
		Name string `json:"name"`
	} `json:"site"`
}

type netBoxResponse struct {
	Results []netBoxDevice `json:"results"`
}

// List queries NetBox for all devices and returns them as model.Device values.
func (n *NetBoxInventory) List(ctx context.Context) ([]model.Device, error) {
	url := n.opts.BaseURL + "/api/dcim/devices/?limit=1000"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("netbox: build request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+n.opts.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("netbox: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("netbox: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("netbox: read body: %w", err)
	}

	var nbResp netBoxResponse
	if err := json.Unmarshal(body, &nbResp); err != nil {
		return nil, fmt.Errorf("netbox: parse response: %w", err)
	}

	devices := make([]model.Device, 0, len(nbResp.Results))
	for _, d := range nbResp.Results {
		devices = append(devices, model.Device{
			ID:           fmt.Sprintf("%d", d.ID),
			Hostname:     d.Name,
			ManagementIP: d.PrimaryIP.Address,
			Site:         d.Site.Name,
		})
	}
	return devices, nil
}

// Get retrieves a single device by its ID from NetBox.
func (n *NetBoxInventory) Get(ctx context.Context, id string) (model.Device, error) {
	url := n.opts.BaseURL + "/api/dcim/devices/" + id + "/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Device{}, fmt.Errorf("netbox: build request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+n.opts.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return model.Device{}, fmt.Errorf("netbox: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Device{}, fmt.Errorf("netbox: device %s not found (status %d)", id, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Device{}, fmt.Errorf("netbox: read body: %w", err)
	}

	var d netBoxDevice
	if err := json.Unmarshal(body, &d); err != nil {
		return model.Device{}, fmt.Errorf("netbox: parse response: %w", err)
	}

	return model.Device{
		ID:           fmt.Sprintf("%d", d.ID),
		Hostname:     d.Name,
		ManagementIP: d.PrimaryIP.Address,
		Site:         d.Site.Name,
	}, nil
}
