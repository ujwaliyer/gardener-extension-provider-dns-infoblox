# Using the Cloudflare dns provider extension with Gardener as operator

The [`core.gardener.cloud/v1beta1.Seed` resource](https://github.com/gardener/gardener/blob/master/example/50-seed.yaml) declares a `dns` field that confgures the default domain that is used for shoots running on that seed.

This document explains the necessary configuration for this provider extension.

## `Seed` resource

This provider extension supports the `spec.dns` field of the seed resource.

### Cloudflare Provider Credentials

In order for Gardener to create dns records using CloudFlare, a Shoot has to provide credentials with sufficient permissions to the desired cloudflare zones.

Create a api token for your cloudflare account as described in https://developers.cloudflare.com/api/tokens/create/ 

This `Secret` must look as follows:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: domain-cloudflare
  namespace: garden-dev
type: Opaque
data:
  apiToken: base64(api token)
```


The creates secret can now be referenced in the seed:

```yaml
---
apiVersion: core.gardener.cloud/v1beta1
kind: Seed
metadata:
  name: my-seed
spec:
  provider:
    type: xxx
    region: xxx
  dns:
    ingressDomain: my-domain.example.com
    provider:
      type: cldouflare-dns
      secretRef:
        name: domain-cloudflare
        namespace: garden-dev
  ...
```
