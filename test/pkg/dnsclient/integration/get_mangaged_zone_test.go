package integration

import (
	// "fmt"
	// ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
	cfg "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test/config"
)

var _ = Describe("GetManagedZone", func() {
	var zone map[string]string
	var default_zone string
	var user string
	var password string
	const a_record_name = "example.com"
	const txt_record_name = "abcd-efgh"
	const dns_view = "default"
	Context("zone ::::----", func() {
		It("Should not create A record :", func() {
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
			Expect(err).To(BeNil())
		})
	})
})
