package integration

import (
	"fmt"
	// ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
)

var _ = Describe("NewDnsclient", func() {

	BeforeEach(func() {
		// infobloxConfig := dnsInfoBlox.InfobloxConfig{}
		// Host := "10.16.198.191"
		// Port := 443
		// SSLVerify := false
		// Version := "2.10"
		// View := ""
		// PoolConnections := 10
		// RequestTimeout := 20
		// // CaCert := nil
		// MaxResults := 100
		// ProxyURL := nil
		// c := dnsInfoBlox.InfobloxConfig{
		// 	&Host,
		// 	&Port,
		// 	&SSLVerify,
		// 	&Version,
		// 	&View,
		// 	&PoolConnections,
		// 	&RequestTimeout,
		// 	nil,
		// 	MaxResults,
		// 	nil,
		// }
		// c := dnsInfoBlox.&InfobloxConfig{
		// 	Host:            "10.16.198.191",
		// 	Port:            "443",
		// 	SSLVerify:       false,
		// 	Version:         2.10,
		// 	View:            "",
		// 	PoolConnections: 10,
		// 	RequestTimeout:  20,
		// 	CaCert:          nil,
		// 	MaxResults:      100,
		// 	ProxyURL:        nil,
		// }
		// dnsC, err := dnsInfoBlox.NewDNSClient(nil, "admin", "infoblox", Host)
		// Accessing struct fields using the dot operator
		// fmt.Println("Car Name: ", c)
		// fmt.Println("Car Color: ", dnsC)
		// fmt.Println(err)

		// conn := dnsInfoBlox.GetInfoBloxInstance()
		// objMgr := ibclient.NewObjectManager(dnsC, "VMWare", "")
		// fmt.Println("Car Color: ", objMgr)
	})
	// Context("with connect api ", func() {
	It("should fill details and return nil ----", func() {
		err := dnsInfoBlox.fillDefaultDetails()
		Expect(err).To(BeNil())
	})
	// })
})
