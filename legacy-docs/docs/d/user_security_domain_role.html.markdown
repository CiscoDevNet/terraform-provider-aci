---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_user_security_domain_role"
sidebar_current: "docs-aci-data-source-aci_user_security_domain_role"
description: |-
  Data source for ACI User Security Domain Role
---

# aci_user_security_domain_role #

Data source for ACI User Security Domain Role


## API Information ##

* `Class` - aaaUserRole
* `Distinguished Name` - uni/userext/user-{name}/userdomain-{name}/role-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Users -> User -> Add User Domain -> Roles 



## Example Usage ##

```hcl
data "aci_user_security_domain_role" "example" {
  user_domain_dn  = aci_user_security_domain.example.id
  name            = "example"
}
```

## Argument Reference ##

* `user_domain_dn` - (Required) Distinguished name of parent UserDomain object.
* `name` - (Required) name of object User Security Domain Role.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the User Security Domain Role.
* `annotation` - (Optional) Annotation of object User Security Domain Role.
* `name_alias` - (Optional) Name Alias of object User Security Domain Role.
* `priv_type` - (Optional) Privilege Type. The privilege type for a user role.
* `description` - (Optional) Description of object User Security Domain Role.