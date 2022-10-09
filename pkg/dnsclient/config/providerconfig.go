package config

import (

	"fmt"
	ib_api "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	"github.com/gardener/gardener/extensions/pkg/controller/common"

)

type valueProvider struct {
	common.ClientContext
}

func (vp *valueProvider) GetProviderConfig(_ context.Context, cp *extensionsv1alpha1.InfobloxConfig, cluster *extensionscontroller.Cluster) struct{} {

	// decode InfobloxConfig
	ibConfig := &ib_api.InfobloxConfig{}
	if cp.Spec.ProviderConfig != nil {
		if _, _, err := vp.Decoder().decode(cp.Spec.ProviderConfig.Raw, nil, ibConfig); err != nil {
			return nil, fmt.Errof("Could not decode providerConfig of infoblox '%s': %w", kutil.ObjectName(cp), err)
		}
	}

	// return getInfobloxConfigValues(ibConfig, cp, cluster)
	return ibConfig

}

// func (vp *valueProvider) getInfobloxConfigValues(ibConfig *ib_api.InfobloxConfig, cp *extensionsv1alpha1.InfobloxConfig)