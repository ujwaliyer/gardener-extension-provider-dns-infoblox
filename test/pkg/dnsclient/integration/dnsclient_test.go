package integration

import (
	"fmt"
	ibclient "github.com/infobloxopen/infoblox-go-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testInfoBlox "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test"
	// "testing"
	// . "testInfoblox"
)

var _ = Describe("Dnsclient", func() {

	var objMgr *ibclient.ObjectManager
	const hostAdd string = "abc.example.com"
	var refString string
	var refCNAMEString string
	// var refString string = "zone_auth/ZG5zLnpvbmUkLl9kZWZhdWx0LmNvbS5leGFtcGxlLmFiYw:abc.example.com/default"

	BeforeEach(func() {
		conn := testInfoBlox.GetInfoBloxInstance()
		objMgr = ibclient.NewObjectManager(conn, "VMWare", "")
		objMgr.OmitCloudAttrs = true
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
			canonical := "test-canonical.domain.com"
			dnsView := "default"
			recordName := "test.example.com"
			useTtl := false
			ttl := uint32(0)
			comment := "test CNAME record creation"
			ea := make(ibclient.EA)
			aCNAMERecord, err := objMgr.CreateCNAMERecord(dnsView, canonical, recordName, useTtl, ttl, comment, ea)
			refCNAMEString = aCNAMERecord.Ref
			Expect(err).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		It("should delete expected CNAME record Ref to DeleteObject", func() {
			dCNAMERef, err = objMgr.DeleteCNAMERecord(refCNAMEString)
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
