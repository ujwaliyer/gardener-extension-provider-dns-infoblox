package functional

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net"
)

var _ = Describe("DnsLookup", func() {
	Context("DNS Lookups----", func() {

		It("Rverse Lookup :", func() {
			namesAssociated, err := net.LookupAddr("10.16.198.140")
			Expect(namesAssociated).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("TXT Lookup :", func() {
			txtRecord, err := net.LookupTXT("owner.vsphere.garden.internal")
			Expect(txtRecord).NotTo(BeNil())
			Expect(err).To(BeNil())

		})
	})
})
