---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_login_domain"
sidebar_current: "docs-aci-resource-login_domain"
description: |-
  Manages ACI Login Domain
---

# aci_login_domain #

Manages ACI Login Domain

## API Information ##

* `Class` - aaaLoginDomain and aaaDomainAuth
* `Distinguished Name` - uni/userext/logindomain-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> AAA -> Policy -> Login Domains 


## Example Usage ##

```hcl
resource "aci_login_domain" "example" {
  name             = "example"
  annotation       = "orchestrator:terraform"
  provider_group   = "example" 
  realm            = "local"
  realm_sub_type   = "default"
  description      = "From Terraform"
  name_alias       = "login_domain_alias"
}
```

## Argument Reference ##


* `name` - (Required) Name of object Login Domain.
* `annotation` - (Optional) Annotation of object Login Domain.
* `provider_group` - (Optional) Provider Group. An AAA configuration provider group is a group of remote servers supporting the same AAA protocol that will be used for authentication and authorization. When a provider group is specified, only the servers within that group will be used for authentication and authorization. If no provider group is specified, all servers supporting the realm of AAA protocols will be used for authentication and authorization. (Note: Attribute provider_group will be set only for the value "none" of attribute "realm", for other values server will not allow to set "provider_group" attribute.) 
* `realm` - (Optional) Realm. The security method for processing authentication requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database. Allowed values are "ldap", "local", "none", "radius", "rsa", "saml", "tacacs". Default value is "local". Type: String.
* `realm_sub_type` - (Optional) Realm subtype that can be default or Duo. Allowed values are "default", "duo". Default value is "default". Type: String. (Note: attribute realm_sub_type is supported for version 5 and above of APIC)
* `name_alias` - (Optional) Name Alias of object Login Domain. Type: String.
* `description` - (Optional) Description of object Login Domain. Type: String.


## Importing ##

An existing LoginDomain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_login_domain.example <Dn>
```