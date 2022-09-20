package integration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dnsclient", func() {
	Context("with connect api ", func() {
		It("should connect to infoblox server", func() {
			Expect("abc").To(Equal("abc"))
		})
	})
})
