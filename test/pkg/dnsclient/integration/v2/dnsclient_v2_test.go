package integration

import (
	"fmt"
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testInfoBloxV2 "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test/v2"
	// "testing"
	// . "testInfoblox"
)

var _ = Describe("Dnsclient", func() {

	var objMgr ibclient.IBObjectManager
	const hostAdd string = "example.com"
	var refString string
	var refCNAMEString string
	var refTXTString string
	var dnsView string = "default"
	var refAString string

	BeforeEach(func() {
		conn := testInfoBloxV2.GetInfoBloxInstanceV2()
		objMgr = ibclient.NewObjectManager(conn, "VMWare", "")
	})
	Context("with connect api ", func() {

		It("should create Zone Auth", func() {
			ea := make(ibclient.EA)
			newZone, znErr := objMgr.CreateZoneAuth(hostAdd, ea)
			refString = newZone.Ref
			zoneAuth, err := objMgr.GetZoneAuthByRef(refString)
			Expect(znErr).To(BeNil())
			Expect(zoneAuth).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("should create CNAME record Object", func() {
			canonical := "test-canonical.example.com"
			recordName := "test.example.com"
			useTtl := false
			ttl := uint32(0)
			comment := "test CNAME record creation"
			ea := make(ibclient.EA)
			aCNAMERecord, err := objMgr.CreateCNAMERecord(dnsView, canonical, recordName, useTtl, ttl, comment, ea)
			refCNAMEString = aCNAMERecord.Ref
			fmt.Println(refCNAMEString)
			Expect(aCNAMERecord).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("should create TXT record Object", func() {
			recordName := "text.example.com"
			aTXTRecord, err := objMgr.CreateTXTRecord(recordName, "domain is assigned to the user", 0, dnsView)
			refTXTString = aTXTRecord.Ref
			Expect(aTXTRecord).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("should create A record Object", func() {
			recordName := "example.com"
			netView := ""
			useTtl := false
			ttl := uint32(0)
			comment := "test A record creation"
			ea := make(ibclient.EA)
			aRecord, err := objMgr.CreateARecord(netView, dnsView, recordName, "10.16.0.0/8", "10.16.1.2", ttl, useTtl, comment, ea)
			refAString = aRecord.Ref
			Expect(aRecord).NotTo(BeNil())
			Expect(err).To(BeNil())
		})

		It("should delete expected A record Ref to DeleteObject", func() {
			dARef, err := objMgr.DeleteARecord(refAString)
			Expect(dARef).To(Equal(refAString))
			Expect(err).To(BeNil())
		})

		It("should delete expected TXT record Ref to DeleteObject", func() {
			dTXTRef, err := objMgr.DeleteTXTRecord(refTXTString)
			Expect(dTXTRef).To(Equal(refTXTString))
			Expect(err).To(BeNil())
		})

		It("should delete expected CNAME record Ref to DeleteObject", func() {
			dCNAMERef, err := objMgr.DeleteCNAMERecord(refCNAMEString)
			Expect(dCNAMERef).To(Equal(refCNAMEString))
			Expect(err).To(BeNil())
		})

		It("should delete Zone Auth", func() {
			delZone, znErr := objMgr.DeleteZoneAuth(refString)
			Expect(znErr).To(BeNil())
			Expect(refString).To(Equal(delZone))
		})
	})
})
