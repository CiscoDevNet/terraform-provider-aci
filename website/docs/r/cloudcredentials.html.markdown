---
layout: "aci"
page_title: "ACI: aci_access_credentialtomanagethecloudresources"
sidebar_current: "docs-aci-resource-access_credentialtomanagethecloudresources"
description: |-
  Manages ACI Access Credential to manage the cloud resources
---

# aci_access_credentialtomanagethecloudresources #

Manages ACI Access Credential to manage the cloud resources

## API Information ##

* `Class` - cloudCredentials
* `Distinguished Name` - uni/tn-{name}/credentials-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_access_credentialtomanagethecloudresources" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
  client_id = 
  email = 
  http_proxy = 
  key = 
  key_id = 

  rsa_private_key = ""

  cloud_rs_ad = aci_resource.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the object Access Credential to manage the cloud resources.
* `annotation` - (Optional) Annotation of the object Access Credential to manage the cloud resources.
* `client_id` - (Optional) Client ID.The client ID (option code 61).
* `email` - (Optional) Credentials email address.The email address of the locally-authenticated user.
* `http_proxy` - (Optional) Http Proxy to connect to cloud provider.
* `key` - (Optional) Secret Key.The key or password used to uniquely identify this configuration object.
* `key_id` - (Optional) Acces Key ID.The authentication key ID.
* `rsa_private_key` - (Optional) RSA Private Key.RSA Secret Key Allowed values are and default value is "".

* `relation_cloud_rs_ad` - (Optional) Represents the relation to a Attachment to billing account (class cloudAD). (Azure only, relation to active directory) Type: String.



## Importing ##

An existing AccessCredentialtomanagethecloudresources can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_credentialtomanagethecloudresources.example <Dn>
```