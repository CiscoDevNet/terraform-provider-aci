---
layout: "aci"
page_title: "ACI: aci_access_credentialtomanagethecloudresources"
sidebar_current: "docs-aci-data-source-access_credentialtomanagethecloudresources"
description: |-
  Data source for ACI Access Credential to manage the cloud resources
---

# aci_access_credentialtomanagethecloudresources #

Data source for ACI Access Credential to manage the cloud resources
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudCredentials
* `Distinguished Name` - uni/tn-{name}/credentials-{name}

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_access_credentialtomanagethecloudresources" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of object Access Credential to manage the cloud resources.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Access Credential to manage the cloud resources.
* `annotation` - (Optional) Annotation of object Access Credential to manage the cloud resources.
* `name_alias` - (Optional) Name Alias of object Access Credential to manage the cloud resources.
* `client_id` - (Optional) Client ID. The client ID (option code 61).
* `email` - (Optional) Credentials email address. The email address of the locally-authenticated user.
* `http_proxy` - (Optional) Http Proxy to connect to cloud provider. 
* `key` - (Optional) Secret Key. The key or password used to uniquely identify this configuration object.
* `key_id` - (Optional) Acces Key ID. The authentication key ID.
* `rsa_private_key` - (Optional) RSA Private Key. RSA Secret Key
