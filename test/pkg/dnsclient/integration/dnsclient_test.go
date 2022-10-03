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
			fmt.Println(refString)
			zoneAuth, err := objMgr.GetZoneAuthByRef(refString)
			fmt.Println(znErr)
			// fmt.Println(newZone)
			// fmt.Println(zoneAuth)
			fmt.Println(err)
			Expect(zoneAuth).NotTo(BeNil())
			// Expect(zoneAuth).To(Equal(newZone))
		})

		It("should delete Zone Auth", func() {
			delZone, znErr := objMgr.DeleteZoneAuth(refString)
			fmt.Println(znErr)
			// fmt.Println(delZone)
			Expect(refString).To(Equal(delZone))
		})
	})
})
