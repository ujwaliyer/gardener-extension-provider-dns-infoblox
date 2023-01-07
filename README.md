# Infoblox-Controller-Deployment-Guide

## Contents

1. An Introduction to the extension
   1. Why it was created
   2. What basic parameters/requirements were considered during creation
   3. What it does solve
2. Development
   1. What sort of implementations/paths to current implementation were taken in to consideration
   2. Major components:
     * **DNS Client** package
     * **DNS Controller** package
   4. Notes on what needs to be kept in mind
3. Tool Versions
4. How to Deploy
   1. How deployment works in general with gardener extensions
   2. Pre-requisites for Setup
   3. Steps to deploy


## Introduction to the Extension

### 1. Why was the extension created
The extension was created to enable gardener to support the use of Infoblox as a DNS provider. 

### 2. What Basic Parameters were considered during the creation
 - The environment the Infoblox DNS extension would be deployed in 
 - The functionality and api/function endpoints desired for correct functioning of the DNS service

### 3. What Does it solve?
The extension allows using Infoblox as a DNS provider with gardener in multiple environments (including air-gapped setups) allowing the user to perform relevant **CRUD** functions on DNS records/recordsets as per their requirement ensuring limited/minimum downtime.

## Development

### 1. What sort of implementations/paths to current implementation were taken in to consideration

The [external-dns-management](https://github.com/gardener/external-dns-management) package was previously used to provide the functionality to use Infoblox as a DNS provider. As per recent design decisions w.r.t gardener development post v1.44, the functionality to include this package via the **dns-external** tag was removed.

The extension uses the aforementioned package for infoblox as a reference to manage DNS-related functions via **DNSRecord**, which is native to Gardener code-base and is the accepted method for handling DNS transactions w.r.t. gardener and related extensions.


### 2. Major Components:

#### The **DNS Client**:
  * The DNS client serves as interface for the extension to interact with the Infoblox setup. 
  * It makes use of the [infoblox-client-sdk](https://github.com/infobloxopen/infoblox-go-client), which in turn leverages the Infoblox WebAPI for performing CRUD functions relevant for DNS records management.

#### The **DNS Controller**: 
  * The controller runs the functionality of the extension as a whole. 
  * It leverages the DNS client for the creation, updation and deletion of DNS records in accordance with the **DNSRecord** object. 
  * It manages the DNS functionality for a Gardener installation, with type-definitions and parameters defined in the *dnsrecord.yaml* file 

## Tool Versions:

The extension has been developed using the following tools with versions as listed below: 

  * go - version 1.19.1
  * Infoblox WebAPI - v2.10 and above
  * garden-setup - v3.41.0
  * gardener/gardener - v1.56.1
  * Infoblox DDI Client - v8.4.7-395215
  * Testing tools:
    * Gingko - v2.1.6
    * Gomega - v1.20.1

## Testing:

The extension was tested using the BDD tooling of Gingko and Gomega with an API led approach.

## How to Deploy:

### 1. How deployment works in general with gardener extensions:

Since the DNS extension is fundamental to the functionality of Gardener, the extension needs to be installed as part of **garden-setup** itself.

### 2. Pre-requisites: System setup

* Container Registry:
  A container repository needs to be setup in order to host the docker image for the extension.
* Tools:
  Certain tools need to be installed in the environment for the extension to work:
  1. helm - the extension uses helm for creation of certain Kubernetes objects, relevant to the extension.
  2. git - the extension requires certain packages to be imported, which requires git to be present

### 3. Steps to deploy:

In order for Gardener to create DNS records using Infoblox, a Shoot has to provide credentials with sufficient permissions to the desired infoblox zones.

Every shoot can either reference these credentials in the shoot manifest using a custom domain (see [Example shoot manifest](#example-shoot-manifest))

This `Secret` must look as follows:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: domain-infoblox
  namespace: garden-dev
type: infoblox-dns
data:
  USERNAME: base64(username)
  PASSWORD: base64(password)
  HOST: base64(host)
```

#### Example `Shoot` manifest

Please find below an example `Shoot` manifest:

```yaml
apiVersion: core.gardener.cloud/v1alpha1
kind: Shoot
metadata:
  name: test-infoblox
  namespace: garden-dev
spec:
  dns:
    domain: my-domain.example.com
    providers:
    - type: infoblox-dns
      secretName: domain-infoblox

  infobloxProfileName: xxx
  region: xxx
  secretBindingName: xxx
  provider:
    type: xxx
    infrastructureConfig:
      ...
    controlPlaneConfig:
      ...
    workers:
    ...
  networking:
    ...
  kubernetes:
    version: 
```

#### Creating the docker image

In case of using a private container registry, the image would need to be generated and pushed. 

* The tar file has been added in the delivery assets folder
* The tar can be converted back to a docker image using the **docker load** command. Refer the [link](https://docs.docker.com/engine/reference/commandline/load/) for relevant steps.
* Rename the docker image and push it to the registry

#### Updating the helm charts

Once the previous step is done, make sure to update the value for the providerConfig in **example/controller-registration.yaml** file as per following steps:

1. Update the container registry details in the **charts/gardener-extension-provider-dns-infoblox/values.yaml** file as required:

```yaml
image:
  repository: <repository-details>
  tag: <tag>
  pullPolicy: IfNotPresent

```
2. Once that is done, convert the **charts/** folder in a zip file
3. Convert the zip file generated to a base64 string using a suitable tool
4. Add the string to the controller-registration yaml under the chart key:

```yaml
apiVersion: core.gardener.cloud/v1beta1
kind: ControllerDeployment
metadata:
  name: provider-dns-infoblox
type: helm
providerConfig:
  chart: 
  <add-string-here>
```


A few changes need to be made in order to include the infoblox extension as part of the install, since the extension isn't officially included in the main **garden-setup** repo:

1. Update the acre.yaml:

  Following segments need to be added in the **acre.yaml** for handling the infoblox configuration as follows:
  
  ```yaml
  credentials:
      infoblox-dns:
        USERNAME: <username>
        PASSWORD: <password>
        HOST: <IP address for Infoblox Server>
        SSL_VERIFY: true/false  
  ```
  Make sure to enable the extension by adding the following entry under **extensions** key:

  ```yaml
  extensions:
    provider-dns-infoblox:
    active: true
  ```

  Add the IP address for the Infoblox server in the landscape.iaas.seeds.dnsServers section, in case we are using separate DNS servers. Refer the segment below as an example:

  ```yaml
  landscape:
    iaas:
      seeds:
        dnsServers:
        - x.x.x.x/x
        - x.x.x.x/x
  ```
**Note** - vSphere at the moment can only support up to 2 DNS server entries. 

Finally, add the following entries for the setup to use infoblox credentials:

```yaml
  dns:                                    # optional for gcp/aws/azure/openstack, default values based on `landscape.iaas`
    type: infoblox-dns 
    credentials: (( .credentials.infoblox-dns ))   # credentials for the dns provider
```

2. Update the **acre.yaml** in the crop folder, to add the entries for infoblox extension:

  ```yaml
  provider-dns-infoblox:
          <<: (( merge ))
          tag: (( .dependency_versions.versions.gardener.extensions.provider-dns-infoblox.version ))
          repo: (( .dependency_versions.versions.gardener.extensions.provider-dns-infoblox.repo ))
          chart_path: charts/gardener-extension-provider-dns-infoblox
          image_tag: (( ~~ ))
          image_repo: (( ~~ ))

  ```

3. Update the **crop/components/gardener/extensions/deployment.yaml** file to include the section for Infoblox:

  ```yaml
  ########################################
  provider-dns-infoblox:
  ########################################
      <<: (( &template ))
      extensionName: provider-dns-infoblox
      type: helm
      providerConfig:
        chart: (( encoded_chart ))
        values:
          <<: (( valuesOverwrite ))
          image:
            repository: (( version.image_repo || ~~ ))
            tag: (( version.image_tag || ~~ ))
          resources:
            limits:
              memory: 1Gi
      resources:
      - kind: DNSRecord
        type: infoblox-dns
        primary: true
  ```


4. Optional: Update the **crop/dependency-versions.yaml** file to include the repo and release version for the infoblox repo:

  ```yaml
  "provider-dns-infoblox": {
    "repo": "<repo-link>",  # add the repo link
    "version": "<version>" # add the release version for the repo
  },
  ```

 In case of using a repository other than Github (eg. Gitlab, Bitbucket etc.), the values for repo and version would need to be modified accordingly.

Once this is done, we can proceed with the standard process for garden-setup installation in order to install Gardener in the respective environment.


### References:

Packages used for developing the extension:

* <https://github.com/gardener/gardener/extensions/pkg/controller>
* <https://github.com/gardener/gardener/extensions/pkg/controller/common>
* <https://github.com/gardener/gardener/extensions/pkg/controller/dnsrecord>
* <https://github.com/gardener/gardener/pkg/apis/core/v1beta1>
* <https://github.com/gardener/gardener/pkg/apis/extensions/v1alpha1>
* <https://github.com/gardener/gardener/pkg/apis/extensions/v1alpha1/helper>
* <https://github.com/gardener/gardener/pkg/controllerutils/reconciler>
* <https://github.com/gardener/gardener/pkg/utils/kubernetes>
