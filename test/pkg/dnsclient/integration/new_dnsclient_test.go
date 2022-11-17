package integration

import (
	// "fmt"
	// ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
	record "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/infoblox"
)

var _ = Describe("NewDnsclient", func() {
	var dnsClient dnsInfoBlox.DNSClient
	var zone map[string]string
	var default_zone = "infobloxbtprpc"
	var user = "admin"
	var password = "btprpc_infoblox"
	var value string
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
			for _, v := range zone {
				value = v
			}
			Expect(err).To(BeNil())
		})
		It("Should not create A record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", value, "examplea", record.Type_A, id_addr, 30)
			Expect(err).NotTo(BeNil())
		})
		It("Should create TXT record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", value, "abcd-efgh"+"."+default_zone, record.Type_TXT, id_addr, 30)
			Expect(err).To(BeNil())
		})

		It("Should create CNAME record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, "default", value, "txt.examplea.com", record.Type_CNAME, id_addr, 30)
			Expect(err).NotTo(BeNil())
		})

		It("Should delete TXT record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, "abcd-efgh"+"."+"infobloxbtprpc", record.Type_TXT)
			Expect(err).To(BeNil())
		})
		It("Should delete A record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, "example"+"."+"infobloxbtprpc", record.Type_A)
			Expect(err).To(BeNil())
		})
		It("Should delete CNAME record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, "def_cname"+"."+"infobloxbtprpc", record.Type_CNAME)
			Expect(err).NotTo(BeNil())
		})
	})
})
