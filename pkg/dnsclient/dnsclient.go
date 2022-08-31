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

type DNSClient interface {
	GetManagedZones(ctx context.Context) (map[string]string, error)
	CreateOrUpdateRecordSet(ctx context.Context, managedZone, name, recordType string, rrdatas []string, ttl int64) error
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

	client, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, &ibclient.WapiHttpRequestor)
	if err != nil {
		fmt.Errorf(err)
	}

	return &dnsClient{
		client: client,
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
func (c *dnsClient) CreateOrUpdateRecordSet(ctx context.Context, zoneID, name, recordType string, rrdatas []string, ttl int64) error {
	records, err := c.getRecordSet(ctx, name, zoneID)
	if err != nil {
		return err
	}
	for _, rrdata := range rrdatas {
		if _, ok := records[rrdata]; ok {
			// entry already exists
			delete(records, rrdata)
			continue
		}
		if err := c.createRecord(ctx, zoneID, name, recordType, rrdata, ttl); err != nil {
			return err
		}
		delete(records, rrdata)
	}

	// delete undefined data
	for _, record := range records {
		if err := c.deleteRecord(ctx, zoneID, record.ID, name, record.Content); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRecordSet deletes the resource recordset with the given name and record type
// in the managed zone with the given name or ID.
func (c *dnsClient) DeleteRecordSet(ctx context.Context, zoneID, name, recordType string) error {
	records, err := c.getRecordSet(ctx, name, zoneID)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Type != recordType {
			continue
		}
		if err := c.deleteRecord(ctx, zoneID, record.ID, name, record.Content); err != nil {
			return err
		}
	}

	return nil
}

// create DNS record for the Infoblox DDI setup
func (c *dnsClient) createRecord(name string, zone, ip_addr string, ttl int64, record_type string) Record {

	// create a DNS record based on record type

	// create a DNS object
	obj_dnsrecord, err := c.CreateObject()

}

func (c *dnsClient) deleteRecord(ctx context.Context, zoneID, recordID, name, rrdata string) error {
	err := c.api.DeleteDNSRecord(ctx, zoneID, recordID)
	if err != nil {
		return fmt.Errorf("Unable to set dns record for %s to %s: %w", name, rrdata, err)
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

func (c *dnsClient) getRecordSet(ctx context.Context, name, zoneID string) map[string]cloudflare.DNSror {
	// results, err := c.api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
	// 	Name: name,
	// })

	//results, err := c.client.

	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	records := make(map[string]cloudflare.DNSRecord, len(results))
	for _, record := range results {
		records[record.Content] = record
	}
	return records, nil
}
