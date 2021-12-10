---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_user_security_domain_role"
sidebar_current: "docs-aci-resource-user_security_domain_role"
description: |-
  Manages ACI User Security Domain Role
---

# aci_user_security_domain_role #

Manages ACI User Security Domain Role

## API Information ##

* `Class` - aaaUserRole
* `Distinguished Named` - uni/userext/user-{name}/userdomain-{name}/role-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Users -> User -> Add User Domain -> Roles


## Example Usage ##

```hcl
resource "aci_user_security_domain_role" "example" {
  user_domain_dn  = aci_user_security_domain.example.id
  annotation      = "orchestrator:terraform"
  name            = "example"
  priv_type       = "readPriv"
  name_alias      = "user_role_alias"
  description     = "From Terraform"
}
```

## Argument Reference ##

* `user_domain_dn` - (Required) Distinguished name of parent UserDomain object.
* `name` - (Required) Name of object User Security Domain Role.
* `annotation` - (Optional) Annotation of object User Security Domain Role. Type: String.
* `name_alias` - (Optional) Name Alias of object User Security Domain Role. Type: String.
* `description` - (Optional) Description of object User Security Domain Role. Type: String.
* `priv_type` - (Optional) Privilege Type.The privilege type for a user role. Allowed values are "readPriv", "writePriv". Default value is "readPriv". Type: String.


## Importing ##

An existing User Security Domain Role can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_user_security_domain_role.example <Dn>
```