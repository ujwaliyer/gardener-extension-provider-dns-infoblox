package integration

import (
	// "fmt"
	// ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	dnsInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient"
	cfg "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test/config"
)

var _ = Describe("NewDnsclient", func() {
	var dnsClient dnsInfoBlox.DNSClient
	// var cnfg cfg.Config
	var zone map[string]string
	var default_zone string
	var user string
	var password string
	var value string
	const a_record_name = "example1.com"
	const txt_record_name = "abcd-efgh1"
	var cname_record_name = "txt1" + "." + a_record_name
	var id_addr = []string{"10.16.2.14"}
	const dns_view = "default"
	BeforeEach(func() {
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
		dnsClient = dnsC
		Expect(dnsC).NotTo(BeNil())
		Expect(err).To(BeNil())
	})
	Context("DNSClient go testing", func() {
		It("GetManaged zone :", func() {
			zones, err := dnsClient.GetManagedZones(nil)
			Ω(zones).Should(ContainElement(ContainSubstring(default_zone), &zone))
			for k := range zone {
				value = k
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
			Expect(err).To(BeNil())
		})
		It("Should delete CNAME record :", func() {
			err := dnsClient.DeleteRecordSet(nil, value, cname_record_name, "CNAME")
			Expect(err).NotTo(BeNil())
		})
	})
})
