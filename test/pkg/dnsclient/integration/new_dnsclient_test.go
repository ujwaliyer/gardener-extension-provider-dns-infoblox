package integration

import (
	// "fmt"
	// ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
)

var _ = Describe("NewDnsclient", func() {
	var dnsClient dnsInfoBlox.DNSClient
	var zone map[string]string
	var default_zone = "sujindar.com"
	var user = "admin"
	var password = "btprpc_infoblox"
	var key string
	var id_addr = []string{"10.16.2.13"}
	BeforeEach(func() {
		Host := "10.16.198.17"
		dnsC, err := dnsInfoBlox.NewDNSClient(nil, user, password, Host)
		dnsClient = dnsC
		Expect(dnsC).NotTo(BeNil())
		Expect(err).To(BeNil())
	})
	Context("DNSClient go testing", func() {
		It("GetManaged zone :", func() {
			zones, err := dnsClient.GetManagedZones(nil)
			Ω(zones).Should(ContainElement(ContainSubstring(default_zone), &zone))
			for k := range zone {
				key = k
			}
			Expect(err).To(BeNil())
		})
		It("Should not create A record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", key, "example.com", "A", id_addr, 30)
			Expect(err).NotTo(BeNil())
		})
		It("Should create TXT record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", key, "abcd-efgh"+"."+default_zone, "TXT", id_addr, 30)
			Expect(err).To(BeNil())
		})

		It("Should create CNAME record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", key, "txt.example.com", "CNAME", id_addr, 30)
			Expect(err).NotTo(BeNil())
		})

		It("Should create TXT record :", func() {
			err := dnsClient.DeleteRecordSet(nil, key, "abcd-efgh"+"."+default_zone, "TXT")
			Expect(err).To(BeNil())
		})
	})
})
