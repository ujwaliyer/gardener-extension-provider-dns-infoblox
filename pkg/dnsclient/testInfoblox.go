package dnsclient

import (
	"fmt"

	ibclient "github.com/infobloxopen/infoblox-go-client"
)

var conn *ibclient.Connector

func GetInfoBloxInstance() *ibclient.Connector {

	if conn == nil {
		hostConfig := ibclient.HostConfig{
			Host:     "10.16.198.191",
			Version:  "2.10",
			Port:     "443",
			Username: "admin",
			Password: "infoblox",
		}
		transportConfig := ibclient.NewTransportConfig("false", 60, 10)
		requestBuilder := &ibclient.WapiRequestBuilder{}
		requestor := &ibclient.WapiHttpRequestor{}
		connection, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
		conn = connection
		if err != nil {
			fmt.Println(err)
			// return err
		}
	}
	return conn
}

// func main() {
// 	connec := getInfoBloxInstance()
// 	connec1 := getInfoBloxInstance()
// 	fmt.Println(connec1)
// 	defer connec.Logout()
// 	objMgr := ibclient.NewObjectManager(connec, "VMWare", "")
// 	// fmt.Println(objMgr.GetLicenceInfo())
// 	fmt.Println(objMgr.GetLicense())
// 	fmt.Println(objMgr.GetAllMembers())
// 	// fmt.Println(objMgr.GetARecordByRef())
// 	// fmt.Println(objMgr.GetCapacityReport())
// 	fmt.Println(objMgr.GetLicense())

// 	// Fetches grid information
// 	fmt.Println(objMgr.GetZoneAuth())
// }
