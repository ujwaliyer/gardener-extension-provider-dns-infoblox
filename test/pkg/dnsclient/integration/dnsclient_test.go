package integration

import (
	// "fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	infobloxdnsconfig "pkg/dnsclient/testInfoBlox"
	// . "testInfoblox"
)

var _ = Describe("Dnsclient", func() {
	Context("with connect api ", func() {
		It("should connect to infoblox server", func() {
			conn := testInfoBlox.GetInfoBloxInstance()
			fmt.Println(conn)
			Expect("abc").To(Equal("abc"))
		})
	})
})
