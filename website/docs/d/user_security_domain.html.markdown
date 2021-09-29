---
layout: "aci"
page_title: "ACI: aci_user_security_domain"
sidebar_current: "docs-aci-data-source-user_security_domain"
description: |-
  Data source for ACI User Security Domain
---

# aci_user_domain #
Data source for ACI User Domain

## API Information ##
* `Class` - aaaUserDomain
* `Distinguished Named` - uni/userext/user-{name}/userdomain-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Users -> User -> Add User Domain

## Example Usage ##

```hcl
data "aci_user_domain" "example" {
  local_user_dn  = aci_local_user.example.id
  name  = "example"
}
```

## Argument Reference ##
* `local_user_dn` - (Required) Distinguished name of parent LocalUser object.
* `name` - (Required) name of object User Domain.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the User Domain.
* `annotation` - (Optional) Annotation of object User Security Domain.
* `name_alias` - (Optional) Name Alias of object User SecDomain.
* `description` - (Optional) Description of object User Security Domain.
