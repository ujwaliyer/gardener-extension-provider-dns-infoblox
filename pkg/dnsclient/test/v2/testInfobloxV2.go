package test

import (
	"fmt"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

var connV2 *ibclient.Connector

func GetInfoBloxInstanceV2() *ibclient.Connector {

	if connV2 == nil {
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
		connV2 = connection
		if err != nil {
			fmt.Println(err)
			// return err
		}
	}
	return connV2
}

// func main() {
// 	connec := GetInfoBloxInstance()
// 	// 	connec1 := getInfoBloxInstance()
// 	// 	fmt.Println(connec1)
// 	// 	defer connec.Logout()
// 	objMgr := ibclient.NewObjectManager(connec, "VMWare", "")
// 	// 	// fmt.Println(objMgr.GetLicenceInfo())
// 	// 	fmt.Println(objMgr.GetLicense())
// 	// 	fmt.Println(objMgr.GetAllMembers())
// 	// 	// fmt.Println(objMgr.GetARecordByRef())
// 	// 	// fmt.Println(objMgr.GetCapacityReport())
// 	// 	fmt.Println(objMgr.GetLicense())

// 	// 	// Fetches grid information
// 	fmt.Println(objMgr.GetZoneAuth())
// }
