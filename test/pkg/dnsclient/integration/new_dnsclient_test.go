package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
)

var _ = Describe("NewDnsclient", func() {
	var dnsClient dnsInfoBlox.DNSClient
	var zone map[string]string
	BeforeEach(func() {
		Host := "10.16.198.191"
		dnsC, err := dnsInfoBlox.NewDNSClient(nil, "admin", "infoblox", Host)
		dnsClient = dnsC
		Expect(dnsC).NotTo(BeNil())
		Expect(err).To(BeNil())
	})
	Context("DNSClient go testing", func() {
		It("GetManaged zone :", func() {
			zones, err := dnsClient.GetManagedZones(nil)
			Ω(zones).Should(ContainElement(ContainSubstring("btprpc.supportnix2.com"), &zone))
			Expect(err).To(BeNil())
		})
	})
})
