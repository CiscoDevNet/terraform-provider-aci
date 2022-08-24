---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_credentials"
sidebar_current: "docs-aci-data-source-aci_cloud_credentials"
description: |-
  Data source for Cloud ACI Access Credential to manage the cloud resources
---

# aci_cloud_credentials #

Data source for Cloud ACI Access Credential to manage the cloud resources
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudCredentials
* `Distinguished Name` - uni/tn-{tenant_name}/credentials-{name}

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> Tenants  -> {tenant_name}



## Example Usage ##

```hcl
data "aci_cloud_credentials" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Access Credential object used to manage the cloud resources.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Access Credential to manage the cloud resources.
* `annotation` - (Optional) Annotation of the Access Credential object to manage the cloud resources.
* `name_alias` - (Optional) Name Alias of the Access Credential object to manage the cloud resources.
* `client_id` - (Optional) Client ID of the Access Credential object.
* `email` - (Optional) email address of the locally-authenticated user.
* `http_proxy` - (Optional) Http Proxy to connect to cloud provider. 
* `key` - (Optional) The Secret key or password used to uniquely identify the cloud resource configuration object.
* `key_id` - (Optional) The Access key ID of the cloud resource.
* `rsa_private_key` - (Optional)  RSA Secret Key of the cloud resource.
