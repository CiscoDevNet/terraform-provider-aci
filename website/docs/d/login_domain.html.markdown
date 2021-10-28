---
layout: "aci"
page_title: "ACI: aci_login_domain"
sidebar_current: "docs-aci-data-source-login_domain"
description: |-
  Data source for ACI Login Domain
---

# aci_login_domain #

Data source for ACI Login Domain


## API Information ##

* `Class` - aaaLoginDomain and aaaDomainAuth
* `Distinguished Named` - uni/userext/logindomain-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> AAA -> Policy -> Login Domains



## Example Usage ##

```hcl
data "aci_login_domain" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object Login Domain.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Login Domain.
* `annotation` - (Optional) Annotation of object Login Domain.
* `name_alias` - (Optional) Name Alias of object Login Domain.
* `description` - (Optional) Description of object Login Domain.
* `provider_group` - (Optional) Provider Group. An AAA configuration provider group is a group of remote servers supporting the same AAA protocol that will be used for authentication and authorization. When a provider group is specified, only the servers within that group will be used for authentication and authorization. If no provider group is specified, all servers supporting the realm of AAA protocols will be used for authentication and authorization.
* `realm` - (Optional) Realm. The security method for processing authentication requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database.
* `realm_sub_type` - (Optional) Realm subtype of object Login Domain.
