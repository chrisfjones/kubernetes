package ports

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/racker/perigee"
)

// AdminState gives users a solid type to work with for create and update
// operations. It is recommended that users use the `Up` and `Down` enums.
type AdminState *bool

// Convenience vars for AdminStateUp values.
var (
	iTrue  = true
	iFalse = false

	Up   AdminState = &iTrue
	Down AdminState = &iFalse
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToPortListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the port attributes you want to see returned. SortKey allows you to sort
// by a particular port attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status       string `q:"status"`
	Name         string `q:"name"`
	AdminStateUp *bool  `q:"admin_state_up"`
	NetworkID    string `q:"network_id"`
	TenantID     string `q:"tenant_id"`
	DeviceOwner  string `q:"device_owner"`
	MACAddress   string `q:"mac_address"`
	ID           string `q:"id"`
	DeviceID     string `q:"device_id"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
}

// ToPortListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPortListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// ports. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those ports that are owned by the tenant
// who submits the request, unless the request is submitted by a user with
// administrative rights.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToPortListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PortPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific port based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToPortCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents the attributes used when creating a new port.
type CreateOpts struct {
	NetworkID      string
	Name           string
	AdminStateUp   *bool
	MACAddress     string
	FixedIPs       interface{}
	DeviceID       string
	DeviceOwner    string
	TenantID       string
	SecurityGroups []string
}

// ToPortCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToPortCreateMap() (map[string]interface{}, error) {
	p := make(map[string]interface{})

	if opts.NetworkID == "" {
		return nil, errNetworkIDRequired
	}
	p["network_id"] = opts.NetworkID

	if opts.DeviceID != "" {
		p["device_id"] = opts.DeviceID
	}
	if opts.DeviceOwner != "" {
		p["device_owner"] = opts.DeviceOwner
	}
	if opts.FixedIPs != nil {
		p["fixed_ips"] = opts.FixedIPs
	}
	if opts.SecurityGroups != nil {
		p["security_groups"] = opts.SecurityGroups
	}
	if opts.TenantID != "" {
		p["tenant_id"] = opts.TenantID
	}
	if opts.AdminStateUp != nil {
		p["admin_state_up"] = &opts.AdminStateUp
	}
	if opts.Name != "" {
		p["name"] = opts.Name
	}
	if opts.MACAddress != "" {
		p["mac_address"] = opts.MACAddress
	}

	return map[string]interface{}{"port": p}, nil
}

// Create accepts a CreateOpts struct and creates a new network using the values
// provided. You must remember to provide a NetworkID value.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToPortCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Response
	_, res.Err = perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{201},
	})

	return res
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToPortUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents the attributes used when updating an existing port.
type UpdateOpts struct {
	Name           string
	AdminStateUp   *bool
	FixedIPs       interface{}
	DeviceID       string
	DeviceOwner    string
	SecurityGroups []string
}

// ToPortUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToPortUpdateMap() (map[string]interface{}, error) {
	p := make(map[string]interface{})

	if opts.DeviceID != "" {
		p["device_id"] = opts.DeviceID
	}
	if opts.DeviceOwner != "" {
		p["device_owner"] = opts.DeviceOwner
	}
	if opts.FixedIPs != nil {
		p["fixed_ips"] = opts.FixedIPs
	}
	if opts.SecurityGroups != nil {
		p["security_groups"] = opts.SecurityGroups
	}
	if opts.AdminStateUp != nil {
		p["admin_state_up"] = &opts.AdminStateUp
	}
	if opts.Name != "" {
		p["name"] = opts.Name
	}

	return map[string]interface{}{"port": p}, nil
}

// Update accepts a UpdateOpts struct and updates an existing port using the
// values provided.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToPortUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("PUT", updateURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200, 201},
	})
	return res
}

// Delete accepts a unique ID and deletes the port associated with it.
func Delete(c *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(c, id), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}
