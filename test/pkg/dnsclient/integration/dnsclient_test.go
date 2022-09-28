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

	Before(func() {
		conn := testInfoBlox.GetInfoBloxInstance()
		objMgr := ibclient.NewObjectManager(connec, "VMWare", "")
	})
	Context("with connect api ", func() {
		It("should get the Zone Auth", func() {
			zoneAuth := objMgr.GetZoneAuth()
			fmt.Println(conn)
			Expect().To(Equal(conn1))
		})
	})
})
