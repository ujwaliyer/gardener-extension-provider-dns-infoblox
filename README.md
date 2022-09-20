# gardener-extension-provider-dns-infoblox

This is an extension developed for using the Infoblox DNS service as a provider, instead of using the **external-dns-management** repo.

# local run test suite

after cloning the project
go mod tidy

then

go mod vendor

ginkgo test/pkg/dnsclient/integration/
