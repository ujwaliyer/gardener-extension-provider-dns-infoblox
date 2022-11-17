package integration

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	cfg "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/dnsclient/test/config"
	"testing"
)

func TestIntegration(t *testing.T) {
	config, err := cfg.LoadConfig()
	fmt.Println(config)
	fmt.Println(err)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
