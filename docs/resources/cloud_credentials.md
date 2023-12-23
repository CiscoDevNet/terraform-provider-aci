---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_credentials"
sidebar_current: "docs-aci-resource-aci_cloud_credentials"
description: |-
  Manages Cloud Network Controller Cloud Credential to manage the cloud resources
---

# aci_cloud_credentials #

Manages Cloud Network Controller Cloud Credential to manage the cloud resources
Note: This resource is supported in Cloud Network Controller only.

## API Information ##

* `Class` - cloudCredentials
* `Distinguished Name` - uni/tn-{name}/credentials-{name}

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> Tenants  -> {tenant_name}


## Example Usage ##

```hcl
resource "aci_cloud_credentials" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
  client_id = "cloud_client_id"
  email = "abc@email.com"
  http_proxy = "proxy"
  key = "secret_key"
  key_id = cloud_key_id"
  rsa_private_key = "rsa_key"
  cloud_rs_ad = aci_resource.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Cloud Credential object to manage the cloud resources.
* `annotation` - (Optional) Annotation of the Cloud Credential object to manage the cloud resources.
* `client_id` - (Optional) Client ID of the Cloud Credential object.
* `email` - (Optional) email address of the locally-authenticated user.
* `http_proxy` - (Optional) HTTP Proxy to connect to the cloud provider. 
* `key` - (Optional) The Secret key or password used to uniquely identify the cloud resource configuration object.
* `key_id` - (Optional) The Access key ID of the cloud resource.
* `rsa_private_key` - (Optional)  RSA Secret Key of the cloud resource.

* `relation_cloud_rs_ad` - (Optional) Represents the relation to an Azure Active Directory (class cloudAD). Type String.



## Importing ##

An existing Cloud Credential to manage the cloud resources can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_credentials.example <Dn>
```