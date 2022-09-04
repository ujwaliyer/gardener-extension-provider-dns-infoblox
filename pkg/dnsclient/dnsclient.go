package dnsclient

import (
	"context"
	"fmt"
	"strconv"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

const (
	type_A     = "A"
	type_CNAME = "CNAME"
	type_AAAA  = "AAAA"
	type_TXT   = "TXT"
)

// type Record interface{}
type Record ibclient.IBObject

var RecordSet []Record

type RecordA ibclient.RecordA
type RecordAAAA ibclient.RecordAAAA
type RecordCNAME ibclient.RecordCNAME
type RecordTXT ibclient.RecordTXT

// var _ Record = (*RecordA)(nil)
// var _ Record = (*RecordAAAA)(nil)
// var _ Record = (*RecordCNAME)(nil)
// var _ Record = (*RecordTXT)(nil)

type RecordNS ibclient.RecordNS

type DNSClient interface {
	GetManagedZones(ctx context.Context) (map[string]string, error)
	CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error
	DeleteRecordSet(ctx context.Context, managedZone, name, recordType string) error
}

type dnsClient struct {
	client ibclient.IBConnector
}

type InfobloxConfig struct {
	Host            *string `json:"host,omitempty"`
	Port            *int    `json:"port,omitempty"`
	SSLVerify       *bool   `json:"sslVerify,omitempty"`
	Version         *string `json:"version,omitempty"`
	View            *string `json:"view,omitempty"`
	PoolConnections *int    `json:"httpPoolConnections,omitempty"`
	RequestTimeout  *int    `json:"httpRequestTimeout,omitempty"`
	CaCert          *string `json:"caCert,omitempty"`
	MaxResults      int     `json:"maxResults,omitempty"`
	ProxyURL        *string `json:"proxyUrl,omitempty"`
}

// NewDNSClient creates a new dns client based on the Infoblox config provided
func NewDNSClient(username string, password string) (DNSClient, error) {

	infobloxConfig := &InfobloxConfig{}

	// define hostConfig
	hostConfig := ibclient.HostConfig{
		host:     *infobloxConfig.Host,
		Port:     *infobloxConfig.Port,
		Version:  *infobloxConfig.Version,
		Username: username,
		Password: password,
	}

	verify := "true"
	if infobloxConfig.SSLVerify != nil {
		verify = strconv.FormatBool(*infobloxConfig.SSLVerify)
	}

	// define transportConfig
	transportConfig := ibclient.NewTransportConfig(verify, infobloxConfig.RequestTimeout, infobloxConfig.PoolConnections)

	var requestBuilder ibclient.HttpRequestBuilder = &ibclient.WapiRequestBuilder{}

	dns_client, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, &ibclient.WapiHttpRequestor)
	if err != nil {
		fmt.Errorf(err)
	}

	return &dnsClient{
		client: dns_client,
	}, nil
}

// GetManagedZones returns a map of all managed zone DNS names mapped to their IDs, composed of the project ID and
// their user assigned resource names.
func (c *dnsClient) GetManagedZones(ctx context.Context) (map[string]string, error) {
	zones := make(map[string]string)

	result, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	for _, z := range result {
		zones[z.Name] = z.ID
	}

	return zones, nil
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *dnsClient) CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error {
	records, err := c.getRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}
	for _, ip_addr := range ip_addrs {
		if _, ok := records[ip_addr]; ok {
			// entry already exists
			delete(records, ip_addr)
			continue
		}
		if err := c.createRecord(name, view, zone, ip_addr, ttl, record_type); err != nil {
			return err
		}
		delete(records, ip_addr)
	}

	// delete undefined data
	for _, record := range records {
		if err := c.deleteRecord(ctx, zone, record.ID, name, record.name); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRecordSet deletes the resource recordset with the given name and record type
// in the managed zone with the given name or ID.
func (c *dnsClient) DeleteRecordSet(ctx context.Context, zone, name, record_type string) error {
	records, err := c.getRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type != record_type {
			continue
		}
		if err := c.deleteRecord(ctx, zone, record.ID, name, record.name); err != nil {
			return err
		}
	}

	return nil
}

// create DNS record for the Infoblox DDI setup
func (c *dnsClient) createRecord(name string, view string, zone string, ip_addr string, ttl int64, record_type string) Record {

	dns_objmgr, err := ibclient.NewObjectManager(c, "VMWare", "")

	var record ibclient.IBObject

	switch record_type {
	case type_A:
		record = dns_objmgr.CreateARecord()
		record.View = view
		record.Name = name
		record.IpV4Addr = ip_addr
		record.Ttl = ttl
	case type_AAAA:
		record = dns_objmgr.CreateAAAARecord()
		record.View = view
		record.Name = name
		record.IpV6Addr = ip_addr
		record.Ttl = ttl
	case type_CNAME:
		record = dns_objmgr.CreateCNAMERecord()
		record.View = view
		record.Name = name
		record.Canonical = ip_addr
		record.Ttl = ttl
	case type_TXT:

	}

	dns_record := ibclient.CreateObject(record.(ibclient.IBObject))

	return dns_record

}

func (c *dnsClient) deleteRecord(record Record, zone string) error {

	_, err := c.client.DeleteObject(record.Ref)

	if err != nil {
		return fmt.Errorf(err)
	}

	return nil

}

func (c *dnsClient) getZoneID(ctx context.Context, name string) (string, error) {
	zones, err := c.GetManagedZones(ctx)
	if err != nil {
		return "", err
	}
	zoneID, ok := zones[name]
	if !ok {
		return "", fmt.Errorf("No zone found for %s", name)
	}
	return zoneID, nil
}

func (c *dnsClient) getRecordSet(name, record_type string, zone string) (map[string]Record, error) {

	results, err := c.client.GetObject()

	if record_type != type_TXT {
		return nil, fmt.Errorf("record type %s not supported for GetRecord", record_type)
	}

	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	records := make(map[string]Record, len(results))
	for _, record := range results {
		records[record.name] = record
	}
	return records, nil
}
