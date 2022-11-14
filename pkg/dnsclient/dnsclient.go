package dnsclient

import (
	"bytes"
	"context"
	raw "dnsclient-poc/raw/records"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	raw "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/infoblox"
)

type DNSClient interface {
	GetManagedZones(ctx context.Context) (map[string]string, error)
	CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, values []string, ttl int64) error
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

func assignDefaultValues() (InfobloxConfig, error) {

	port := 443
	view := "default"
	poolConnections := 10
	requestTimeout := 60
	version := "2.10"

	return InfobloxConfig{
		Port:            &port,
		View:            &view,
		PoolConnections: &poolConnections,
		RequestTimeout:  &requestTimeout,
		Version:         &version,
	}, nil

}

// NewDNSClient creates a new dns client based on the Infoblox config provided
func NewDNSClient(ctx context.Context, username string, password string, host string) (DNSClient, error) {

	infobloxConfig, err := assignDefaultValues()
	if err != nil {
		fmt.Println(err)
	}

	// define hostConfig
	hostConfig := ibclient.HostConfig{
		// Host:     *infobloxConfig.Host,
		Host:     host,
		Port:     strconv.Itoa(*infobloxConfig.Port),
		Version:  *infobloxConfig.Version,
		Username: username,
		Password: password,
	}

	// verify := "true"
	verify := "false"
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

	conn := c.client.(*ibclient.Connector)

	rt := ibclient.NewZoneAuth(ibclient.ZoneAuth{})
	urlStr := conn.RequestBuilder.BuildUrl(ibclient.GET, rt.ObjectType(), "", rt.ReturnFields(), &ibclient.QueryParams{})

	req, err := http.NewRequest("GET", urlStr, new(bytes.Buffer))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(conn.HostConfig.Username, conn.HostConfig.Password)

	resp, err := conn.Requestor.SendRequest(req)
	if err != nil {
		fmt.Println(err)
	}

	rs := []ibclient.ZoneAuth{}
	err = json.Unmarshal(resp, &rs)
	if err != nil {
		fmt.Println(err)
	}

	ZoneList := make(map[string]string)

	for _, zone := range rs {
		// zone_list = append(zone_list, zone.Fqdn)
		ZoneList[zone.Ref] = zone.Fqdn
	}

	return ZoneList, nil
}

// CreateOrUpdateRecordSet creates or updates the resource recordset with the given name, record type, rrdatas, and ttl
// in the managed zone with the given name or ID.
func (c *dnsClient) CreateOrUpdateRecordSet(ctx context.Context, view, zone, name, record_type string, values []string, ttl int64) error {
	records, err := c.GetRecordSet(name, record_type, zone)
	if err != nil {
		return err
	}

	for _, r := range records {
		if r.GetDNSName() == name {
			err_del := c.DeleteRecord(r.(raw.Record), zone)
			if err_del != nil {
				return err_del
			}
		}
	}

	for _, value := range values {
		_, err := c.createRecord(name, view, value, ttl, record_type)
		if err != nil {
			return err
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
func (c *dnsClient) createRecord(name string, view string, value string, ttl int64, record_type string) (string, error) {

	var record string
	var err error
	var rec ibclient.IBObject

	switch record_type {
	case raw.Type_A:
		rec = ibclient.NewRecordA(view, "", name, value, uint32(ttl), false, "", nil, "")

	case raw.Type_AAAA:
		rec = ibclient.NewRecordAAAA(view, name, value, false, uint32(ttl), "", nil, "")

	case raw.Type_CNAME:
		rec = ibclient.NewRecordCNAME(view, value, name, true, uint32(ttl), "", nil, "")

	case raw.Type_TXT:
		rec = ibclient.NewRecordTXT(ibclient.RecordTXT{
			Name: name,
			View: view,
			Text: value,
		})
	}

	record, err = c.client.CreateObject(rec)
	if err != nil {
		return "", err
	}

	return record, nil
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

	execRequest := func(forceProxy bool, recordType string) ([]byte, error) {
		var rt ibclient.IBObject

		switch recordType {
		case "A":
			rt = ibclient.NewRecordTXT(ibclient.RecordTXT{})
		case "CNAME":
			rt = ibclient.NewEmptyRecordA()
		}

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

	resp, err := execRequest(false, record_type)
	if err != nil {
		// Forcing the request to redirect to Grid Master by making forcedProxy=true
		resp, err = execRequest(true, record_type)
		// fmt.Println(err)
	}
	if err != nil {
		return nil, err
	}

	switch record_type {
	case raw.Type_TXT:
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
	case raw.Type_A:
		rs := []raw.RecordA{}
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

	return nil, nil

}
