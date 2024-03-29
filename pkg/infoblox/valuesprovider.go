// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package infoblox

import (
	"github.com/gardener/gardener/extensions/pkg/controller/common"
	"github.com/go-logr/logr"
	// types "github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/apis/config/v1alpha1/types"
	//"github.com/ujwaliyer/gardener-extension-provider-dns-infoblox/pkg/controller/dnsrecord/actuator"
)

type actuator struct {
	common.ClientContext
	logger logr.Logger
}

// func (a *actuator) GetControlPlaneChartValues(
// 	//_ context.Context,
// 	dns *extensionsv1alpha1.DNSRecord,
// 	//cluster *extensionscontroller.Cluster,
// 	//secretsReader secretsmanager.Reader,
// 	//checksums map[string]string,
// 	//scaledDown bool,
// ) (map[string]interface{}, error) {
// 	// Decode providerConfig
// 	cpConfig := &types.ProviderConfigManager{}
// 	if dns.Spec.ProviderConfig != nil {
// 		if _, _, err := a.Decoder().Decode(dns.Spec.ProviderConfig.Raw, nil, cpConfig); err != nil {
// 			return nil, fmt.Errorf("could not decode providerConfig of dnsrecord '%s': %w", kutil.ObjectName(dns), err)
// 		}
// 	}

// 	return cpConfig
// }

/*
func ExtractCredentials(cpConfig struct) (*Credentials, error) {

	var hostName, view *string

	alt hostName = pointer.String(Host)
	alt view = pointer.String(View)


	hostName := getOptional(secret, UserName, altUserNameKey)
	password := getOptional(secret, Password, altPasswordKey)
	applicationCredentialID := getOptional(secret, ApplicationCredentialID, altApplicationCredentialID)
	applicationCredentialName := getOptional(secret, ApplicationCredentialName, altApplicationCredentialName)
	applicationCredentialSecret := getOptional(secret, ApplicationCredentialSecret, altApplicationCredentialSecret)
	authURL := getOptional(secret, AuthURL, altAuthURLKey)

	if password != "" {
		if applicationCredentialSecret != "" {
			return nil, fmt.Errorf("cannot specify both '%s' and '%s' in secret %s/%s", Password, ApplicationCredentialSecret, secret.Namespace, secret.Name)
		}
		if userName == "" {
			return nil, fmt.Errorf("'%s' is required if '%s' is given in %s/%s", UserName, Password, secret.Namespace, secret.Name)
		}
	} else {
		if applicationCredentialSecret == "" {
			return nil, fmt.Errorf("must either specify '%s' or '%s' in secret %s/%s", Password, ApplicationCredentialSecret, secret.Namespace, secret.Name)
		}
		if applicationCredentialID == "" {
			if userName == "" || applicationCredentialName == "" {
				return nil, fmt.Errorf("'%s' and '%s' are required if application credentials are used without '%s' in secret %s/%s", ApplicationCredentialName, UserName,
					ApplicationCredentialID, secret.Namespace, secret.Name)
			}
		}
	}

	return &Credentials{
		View:                  view,
		TenantName:                  tenantName,
		Username:                    userName,
		Password:                    password,
		ApplicationCredentialID:     applicationCredentialID,
		ApplicationCredentialName:   applicationCredentialName,
		ApplicationCredentialSecret: applicationCredentialSecret,
		AuthURL:                     string(authURL),
	}, nil
}
*/
