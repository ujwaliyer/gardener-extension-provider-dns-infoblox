package dnsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	ibclient "github.com/infobloxopen/infoblox-go-client"
	ibclient2 "github.com/infobloxopen/infoblox-go-client/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	type_A     = "record:a"
	type_CNAME = "record:cname"
	type_AAAA  = "record:aaaa"
	type_TXT   = "record:txt"
)

// type Record interface{}
type Record ibclient.IBObject

type RecordSet []Record

type RecordA ibclient.RecordA
type RecordAAAA ibclient2.RecordAAAA
type RecordCNAME ibclient.RecordCNAME
type RecordTXT ibclient.RecordTXT

var _ Record = (*RecordA)(nil)
var _ Record = (*RecordAAAA)(nil) // doesn't resolve in v2
var _ Record = (*RecordCNAME)(nil)
var _ Record = (*RecordTXT)(nil)

type RecordNS ibclient2.RecordNS

type DNSClient interface {
	GetManagedZones(ctx context.Context) ([]string, error)
	CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error
	DeleteRecordSet(ctx context.Context, managedZone, name, recordType string) error
}

type dnsClient struct {
	client ibclient.Connector
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
func NewDNSClient(ctx context.Context, username string, password string) (DNSClient, error) {

	infobloxConfig := &InfobloxConfig{}

	// define hostConfig
	hostConfig := ibclient.HostConfig{
		Host:     *infobloxConfig.Host,
		Port:     strconv.Itoa(*infobloxConfig.Port),
		Version:  *infobloxConfig.Version,
		Username: username,
		Password: password,
	}

	verify := "true"
	if infobloxConfig.SSLVerify != nil {
		verify = strconv.FormatBool(*infobloxConfig.SSLVerify)
	}

	if infobloxConfig.CaCert != nil && verify == "true" {
		tmpfile, err := ioutil.TempFile("", "cacert")
		if err != nil {
			return nil, fmt.Errorf("cannot create temporary file for cacert: %w", err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write([]byte(*infobloxConfig.CaCert)); err != nil {
			return nil, fmt.Errorf("cannot write temporary file for cacert: %w", err)
		}
		if err := tmpfile.Close(); err != nil {
			return nil, fmt.Errorf("cannot close temporary file for cacert: %w", err)
		}
		verify = tmpfile.Name()
	}

	// define transportConfig
	transportConfig := ibclient.NewTransportConfig(verify, *infobloxConfig.RequestTimeout, *infobloxConfig.PoolConnections)

	var requestBuilder ibclient.HttpRequestBuilder = &ibclient.WapiRequestBuilder{}

	dns_client, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, &ibclient.WapiHttpRequestor{})
	if err != nil {
		fmt.Println(err)
	}

	// todo: set correct type for dns_client to create dns_object
	dns_object := ibclient2.CreateObject(dns_client.(ibclient.IBObject))

	return &dnsClient{
		client: dns_object,
	}, nil
}

// get DNS client from secret reference
// todo: rewrite client parameter
// func (c *dnsClient) NewDNSClientFromSecretRef(ctx context.Context, cl client.Client, secretRef corev1.SecretReference) (DNSClient, error) {
func NewDNSClientFromSecretRef(ctx context.Context, c client.Client, secretRef corev1.SecretReference) (DNSClient, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, c, &secretRef)
	if err != nil {
		return nil, err
	}

	username, ok := secret.Data["username"]
	if !ok {
		return nil, fmt.Errorf("No username found")
	}

	password, ok := secret.Data["password"]
	if !ok {
		return nil, fmt.Errorf("No password found")
	}

	return NewDNSClient(ctx, string(username), string(password))

}

// GetManagedZones returns a map of all managed zone DNS names mapped to their IDs, composed of the project ID and
// their user assigned resource names.
func (c *dnsClient) GetManagedZones(ctx context.Context) ([]string, error) {

	// get all zones
	objMgr := ibclient.NewObjectManager(c, "VMWare", "")

	// todo: getzoneauth only supported in v1; record creation needs v2; what to do?
	all_zones, err := objMgr.GetZoneAuth()
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	var zone_list []string

	for _, zone := range all_zones {
		zone_list = append(zone_list, zone.Fqdn)
	}

	return zone_list, nil
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *dnsClient) CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error {
	records, err := c.getRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}

	// need second look at logic
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

	// dns_objmgr := ibclient.NewObjectManager(c, "VMWare", "")

	var record Record

	switch record_type {
	case type_A:
		record = ibclient.NewEmptyRecordA()
		record.View = view
		record.Name = name
		record.IpV4Addr = ip_addr
		record.Ttl = ttl
	case type_AAAA:
		record := ibclient2.NewEmptyRecordAAAA()
		record.View = view
		record.Name = name
		record.IpV6Addr = ip_addr
		record.Ttl = ttl
	case type_CNAME:
		record := ibclient2.NewEmptyRecordCNAME()
		record.View = view
		record.Name = name
		record.Canonical = ip_addr
		record.Ttl = ttl
	case type_TXT:
		if n, err := strconv.Unquote(value); err == nil && !strings.Contains(value, " ") {
			value = n
		}
		record = (*RecordTXT)(ibclient.NewRecordTXT(ibclient.RecordTXT{
			Name: name,
			Text: ip_addr,
			View: c.View,
		}))
	}

	dns_record := ibclient2.CreateObject(record.(ibclient.IBObject))
	return dns_record

}

func (c *dnsClient) deleteRecord(record Record, zone string) error {

	_, err := c.client.DeleteObject(record.Ref)

	if err != nil {
		return fmt.Errorf(err)
	}

	return nil

}

func (c *dnsClient) getRecordSet(name, record_type string, zone string) ([]RecordSet, error) {

	results, err := c.client.GetObject()

	if record_type != type_TXT {
		return nil, fmt.Errorf("record type %s not supported for GetRecord", record_type)
	}

	execRequest := func(forceProxy bool) ([]byte, error) {
		rt := ibclient.NewRecordTXT(ibclient.RecordTXT{})
		urlStr := c.RequestBuilder.BuildUrl(ibclient.GET, rt.ObjectType(), "", rt.ReturnFields(), &ibclient.QueryParams{})
		urlStr += "&name=" + dnsName
		if forceProxy {
			urlStr += "&_proxy_search=GM"
		}
		req, err := http.NewRequest("GET", urlStr, new(bytes.Buffer))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(c.HostConfig.Username, c.HostConfig.Password)

		return c.Requestor.SendRequest(req)
	}

	resp, err := execRequest(false)
	if err != nil {
		// Forcing the request to redirect to Grid Master by making forcedProxy=true
		resp, err = execRequest(true)
	}
	if err != nil {
		return nil, err
	}

	rs := []RecordTXT{}
	err = json.Unmarshal(resp, &rs)
	if err != nil {
		return nil, err
	}

	rs2 := []RecordSet{}
	for _, r := range rs {
		rs2 = append(rs2, r.Copy())
	}
	return rs2, nil
}
