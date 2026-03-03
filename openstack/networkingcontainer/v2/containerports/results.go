package containerports

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a port resource.
func (r commonResult) Extract() (*ContainerPort, error) {
	var s ContainerPort
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "container_port")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Port.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Port.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Port.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AddressPair contains the IP Address and the MAC address.
type AddressPair struct {
	IPAddress  string `json:"ip_address,omitempty"`
	MACAddress string `json:"mac_address,omitempty"`
}

type ContainerPort struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	Name           string    `json:"name"`
	NetworkID      string    `json:"network_id"`
	MacAddress     string    `json:"mac_address"`
	ParentPort     string    `json:"parent_port"`
	AdminStateUp   bool      `json:"admin_state_up"`
	Status         string    `json:"status"`
	SegmentationID int       `json:"segmentation_id"`
	DeviceOwner    string    `json:"device_owner"`
	FixedIPs       []FixedIP `json:"fixed_ips"`
}

func (r *ContainerPort) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	return nil
}

type FixedIP struct {
	SubnetID   string     `json:"subnet_id"`
	IPAddress  string     `json:"ip_address,omitempty"`
	SubnetInfo SubnetInfo `json:"subnet_info"`
}

type SubnetInfo struct {
	Mask       string `json:"mask"`
	GatewayIP  string `json:"gateway_ip"`
	GatewayMac string `json:"gateway_mac"`
}

// PortPage is the page returned by a pager when traversing over a collection
// of network ports.
type PortPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of ports has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r PortPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"ports_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a PortPage struct is empty.
func (r PortPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractPorts(r)
	return len(is) == 0, err
}

// ExtractPorts accepts a Page struct, specifically a PortPage struct,
// and extracts the elements into a slice of Port structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPorts(r pagination.Page) ([]ContainerPort, error) {
	var s []ContainerPort
	err := ExtractPortsInto(r, &s)
	return s, err
}

func ExtractPortsInto(r pagination.Page, v interface{}) error {
	return r.(PortPage).Result.ExtractIntoSlicePtr(v, "ports")
}
