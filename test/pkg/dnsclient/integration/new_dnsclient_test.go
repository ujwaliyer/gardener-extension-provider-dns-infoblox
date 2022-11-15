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
	const default_zone = "sujindar.com"
	const user = "admin"
	const password = "btprpc_infoblox"
	var value string
	const a_record_name = "example.com"
	const txt_record_name = "abcd-efgh"
	var cname_record_name = "txt" + "." + a_record_name
	var id_addr = []string{"10.16.2.13"}
	const dns_view = "default"
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
			err := dnsClient.CreateOrUpdateRecordSet(nil, dns_view, value, a_record_name, "A", id_addr, 30)
			Expect(err).NotTo(BeNil())
		})
		It("Should create TXT record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, dns_view, value, txt_record_name+"."+value, "TXT", id_addr, 30)
			Expect(err).To(BeNil())
		})

		It("Should create CNAME record :", func() {
			err := dnsClient.CreateOrUpdateRecordSet(nil, dns_view, value, cname_record_name, "CNAME", id_addr, 30)
			Expect(err).NotTo(BeNil())
		})

		It("Should delete TXT record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, txt_record_name+"."+value, "TXT")
			Expect(err).To(BeNil())
		})
		It("Should delete A record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, a_record_name, "A")
			Expect(err).NotTo(BeNil())
		})
		It("Should delete CNAME record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, cname_record_name, "CNAME")
			Expect(err).NotTo(BeNil())
		})
	})
})
