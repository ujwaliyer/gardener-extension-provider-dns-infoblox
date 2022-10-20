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
	ibclient1 "github.com/infobloxopen/infoblox-go-client"
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	raw "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/infoblox"
)

type DNSClient interface {
	GetManagedZones(ctx context.Context) (map[string]string, error)
	CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error
	DeleteRecordSet(ctx context.Context, zone, name, recordType string) error
}

type dnsClient struct {
	client ibclient.IBConnector
}

type RecordSet []raw.Base_Record

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

func (ib *InfobloxConfig) fillDefaultDetails() error {
	if ib.RequestTimeout == nil {
		*ib.RequestTimeout = 60
	} else {
		return fmt.Errorf("no valid value for RequestTimeout")
	}

	if ib.Version == nil {
		*ib.Version = "2.10"
	} else {
		return fmt.Errorf("no valid value for API version")
	}

	if ib.View == nil {
		*ib.View = "default"
	} else {
		return fmt.Errorf("no valid value for view")
	}

	if ib.PoolConnections == nil {
		*ib.PoolConnections = 60
	} else {
		return fmt.Errorf("no valid value for no. of pool connections")
	}

	return nil
}

// NewDNSClient creates a new dns client based on the Infoblox config provided
func NewDNSClient(ctx context.Context, username string, password string, host string) (DNSClient, error) {

	infobloxConfig := &InfobloxConfig{}

	// define hostConfig
	hostConfig := ibclient.HostConfig{
		// Host:     *infobloxConfig.Host,
		Host:     host,
		Port:     strconv.Itoa(*infobloxConfig.Port),
		Version:  *infobloxConfig.Version,
		Username: username,
		Password: password,
	}

	err := infobloxConfig.fillDefaultDetails()
	if err != nil {
		fmt.Println(err)
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
	// dns_object := ibclient.CreateObject(dns_client.(ibclient.IBObject))

	return &dnsClient{
		client: dns_client,
	}, nil
}

// get DNS client from secret reference
// func (c *dnsClient) NewDNSClientFromSecretRef(ctx context.Context, cl client.Client, secretRef corev1.SecretReference) (DNSClient, error) {
func NewDNSClientFromSecretRef(ctx context.Context, c client.Client, secretRef corev1.SecretReference) (DNSClient, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, c, &secretRef)
	if err != nil {
		return nil, err
	}

	username, ok := secret.Data["username"]
	if !ok {
		return nil, fmt.Errorf("no username found")
	}

	password, ok := secret.Data["password"]
	if !ok {
		return nil, fmt.Errorf("no password found")
	}

	// placeholder for host details using providerConfig
	host, ok := secret.Data["host"]
	if !ok {
		return nil, fmt.Errorf("no host details found")
	}

	return NewDNSClient(ctx, string(host), string(username), string(password))

}

// GetManagedZones returns a map of all managed zone DNS names mapped to their IDs, composed of the project ID and
// their user assigned resource names.
func (c *dnsClient) GetManagedZones(ctx context.Context) (map[string]string, error) {

	// dummy code for testing; might remove later
	infobloxConfig := &InfobloxConfig{}

	// define hostConfig
	hostConfig := ibclient1.HostConfig{
		Host:    *infobloxConfig.Host,
		Port:    strconv.Itoa(*infobloxConfig.Port),
		Version: *infobloxConfig.Version,
	}
	verify := "false"

	// define transportConfig
	transportConfig := ibclient1.NewTransportConfig(verify, *infobloxConfig.RequestTimeout, *infobloxConfig.PoolConnections)

	var requestBuilder ibclient1.HttpRequestBuilder = &ibclient1.WapiRequestBuilder{}

	dns_client1, err := ibclient1.NewConnector(hostConfig, transportConfig, requestBuilder, &ibclient1.WapiHttpRequestor{})
	if err != nil {
		fmt.Println(err)
	}

	// get all zones; need separate connector for using this function
	objMgr := ibclient1.NewObjectManager(dns_client1, "VMWare", "")

	// todo: getzoneauth only supported in v1; record creation needs v2; what to do?
	all_zones, err := objMgr.GetZoneAuth()
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	// var zone_list []string
	zone_list := make(map[string]string)

	for _, zone := range all_zones {
		// zone_list = append(zone_list, zone.Fqdn)
		zone_list[zone.Ref] = zone.Fqdn
	}

	return zone_list, nil
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *dnsClient) CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, ip_addrs []string, ttl int64) error {
	records, err := c.GetRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}

	err_del := c.DeleteRecordSet(ctx, zone, name, record_type)
	if err_del != nil {
		fmt.Println(err_del)
	}

	for _, r := range records {
		r0 := c.NewRecord(r.GetDNSName(), view, zone, r.GetValue(), int64(r.GetTTL()), r.GetType())
		err := c.CreateRecord(r0, zone)
		if err != nil {
			fmt.Println(err)
		}
	}

	return err
}

// DeleteRecordSet deletes the resource recordset with the given name and record type
// in the managed zone with the given name or ID.
func (c *dnsClient) DeleteRecordSet(ctx context.Context, zone, name, record_type string) error {
	records, err := c.GetRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}

	for _, rec := range records {
		if rec.GetId() != "" {
			err := c.DeleteRecord(rec.(raw.Record), zone)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// create DNS record for the Infoblox DDI setup
func (c *dnsClient) NewRecord(name string, view string, zone string, value string, ttl int64, record_type string) (record raw.Record) {

	switch record_type {
	case raw.Type_A:
		r := ibclient.NewEmptyRecordA()
		r.View = view
		r.Name = name
		r.Ipv4Addr = value
		r.Ttl = uint32(ttl)
		record = (*raw.RecordA)(r)
	case raw.Type_AAAA:
		r := ibclient.NewEmptyRecordAAAA()
		r.View = view
		r.Name = name
		r.Ipv6Addr = value
		r.Ttl = uint32(ttl)
		record = (*raw.RecordAAAA)(r)
	case raw.Type_CNAME:
		r := ibclient.NewEmptyRecordCNAME()
		r.View = view
		r.Name = name
		r.Canonical = value
		r.Ttl = uint32(ttl)
		record = (*raw.RecordCNAME)(r)
	case raw.Type_TXT:
		if n, err := strconv.Unquote(value); err == nil && !strings.Contains(value, " ") {
			value = n
		}
		record = (*raw.RecordTXT)(ibclient.NewRecordTXT(ibclient.RecordTXT{
			Name: name,
			Text: value,
			View: view,
		}))
	}

	return
}

func (c *dnsClient) CreateRecord(r raw.Record, zone string) error {

	_, err := c.client.CreateObject(r.(ibclient.IBObject))
	return err

}

func (c *dnsClient) DeleteRecord(record raw.Record, zone string) error {

	_, err := c.client.DeleteObject(record.GetId())

	if err != nil {
		return err
	}

	return nil

}

func (c *dnsClient) GetRecordSet(name, record_type string, zone string) (RecordSet, error) {

	results := c.client.(*ibclient.Connector)

	if record_type != raw.Type_TXT {
		return nil, fmt.Errorf("record type %s not supported for GetRecord", record_type)
	}

	execRequest := func(forceProxy bool) ([]byte, error) {
		rt := ibclient.NewRecordTXT(ibclient.RecordTXT{})
		urlStr := results.RequestBuilder.BuildUrl(ibclient.GET, rt.ObjectType(), "", rt.ReturnFields(), &ibclient.QueryParams{})
		urlStr += "&name=" + name
		if forceProxy {
			urlStr += "&_proxy_search=GM"
		}
		req, err := http.NewRequest("GET", urlStr, new(bytes.Buffer))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(results.HostConfig.Username, results.HostConfig.Password)

		return results.Requestor.SendRequest(req)
	}

	resp, err := execRequest(false)
	if err != nil {
		// Forcing the request to redirect to Grid Master by making forcedProxy=true
		resp, err = execRequest(true)
		// fmt.Println(err)
	}
	if err != nil {
		return nil, err
	}

	rs := []raw.RecordTXT{}
	err = json.Unmarshal(resp, &rs)
	if err != nil {
		return nil, err
	}

	rs2 := RecordSet{}
	for _, r := range rs {
		rs2 = append(rs2, r.Copy())
	}
	return rs2, nil
}
