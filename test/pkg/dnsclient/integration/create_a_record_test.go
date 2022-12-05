package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
	cfg "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test/config"
)

var _ = Describe("CreateARecord", func() {
	var zone map[string]string
	var default_zone string
	var user string
	var password string
	var value string
	const a_record_name = "example.com"
	const txt_record_name = "abcd-efgh"
	var id_addr = []string{"10.16.2.13"}
	const dns_view = "default"
	Context("zone ::::----", func() {
		It("GetManaged zone :", func() {
			config := cfg.GetConfig()
			Host := config.Host
			user = config.Username
			password = config.Password
			default_zone = config.DefaultZone

			Expect(user).NotTo(BeEmpty())
			Expect(user).NotTo(BeNil())
			Expect(user).NotTo(Equal(""))

			Expect(password).NotTo(BeEmpty())
			Expect(password).NotTo(BeNil())
			Expect(password).NotTo(Equal(""))

			Expect(default_zone).NotTo(BeEmpty())
			Expect(default_zone).NotTo(BeNil())
			Expect(default_zone).NotTo(Equal(""))

			Expect(Host).NotTo(BeEmpty())
			Expect(Host).NotTo(BeNil())
			Expect(Host).NotTo(Equal(""))

			dnsC, err := dnsInfoBlox.NewDNSClient(nil, user, password, Host)
			Expect(err).To(BeNil())

			zones, err := dnsC.GetManagedZones(nil)
			Ω(zones).Should(ContainElement(ContainSubstring(default_zone), &zone))
			for k := range zone {
				value = k
			}
			Expect(err).To(BeNil())

			err2 := dnsC.CreateOrUpdateRecordSet(nil, dns_view, value, a_record_name, "A", id_addr, 30)
			Expect(err2).NotTo(BeNil())
		})
	})
})
