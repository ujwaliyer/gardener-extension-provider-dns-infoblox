package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
)

var _ = Describe("NewDnsclient", func() {

	Before(func() {
		c := &InfobloxConfig{
			Host:            "10.16.198.191",
			Port:            "443",
			SSLVerify:       false,
			Version:         2.10,
			View:            "",
			PoolConnections: 10,
			RequestTimeout:  20,
			CaCert:          nil,
			MaxResults:      100,
			ProxyURL:        nil,
		}
		dnsC := dnsInfoBlox.NewDNSClient(nil, "admin", "infoblox")
		// Accessing struct fields using the dot operator
		fmt.Println("Car Name: ", c.Host)
		fmt.Println("Car Color: ", c.Port)

		// conn := dnsInfoBlox.GetInfoBloxInstance()
		// objMgr := ibclient.NewObjectManager(connec, "VMWare", "")
	})
	Context("with connect api ", func() {
		It("should get the Zone Auth", func() {
			zoneAuth := objMgr.GetZoneAuth()
			fmt.Println(conn)
			Expect().To(Equal(conn1))
		})
	})
})
