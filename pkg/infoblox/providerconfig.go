package infoblox

// import (
// 	"fmt"

// 	"github.com/gardener/gardener/extensions/pkg/controller/common"
// 	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
// 	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
// 	"github.com/go-logr/logr"
// 	types "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/apis/config/v1alpha1"
// )

// type valuesProvider struct {
// 	common.ClientContext
// 	logger logr.Logger
// }

// type providerConfig struct {
// 	view            string
// 	version         int
// 	poolConnections int
// 	sslVerify       bool
// }

// func (vp *valuesProvider) GetDNSRecordValues(
// 	//_ context.Context,
// 	dns *extensionsv1alpha1.DNSRecord,
// 	//cluster *extensionscontroller.Cluster,
// 	//secretsReader secretsmanager.Reader,
// 	//checksums map[string]string,
// 	//scaledDown bool,
// ) (*types.ProviderConfigManager, error) {
// 	// Decode providerConfig
// 	cpConfig := &types.ProviderConfigManager{}
// 	if dns.Spec.ProviderConfig != nil {
// 		if _, _, err := vp.Decoder().Decode(dns.Spec.ProviderConfig.Raw, nil, cpConfig); err != nil {
// 			return cpConfig, fmt.Errorf("could not decode providerConfig of dnsrecord '%s': %w", kutil.ObjectName(dns), err)
// 		}
// 	}

// 	return cpConfig, nil
// }
