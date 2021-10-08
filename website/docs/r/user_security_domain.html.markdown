---
layout: "aci"
page_title: "ACI: aci_user_security_domain"
sidebar_current: "docs-aci-resource-user_security_domain"
description: |-
  Manages ACI User Security Domain
---

# aci_user_security_domain #
Manages ACI User Security Domain

## API Information ##
* `Class` - aaaUserDomain
* `Distinguished Named` - uni/userext/user-{name}/userdomain-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Users -> User -> Add User Domain

## Example Usage ##

```hcl
resource "aci_user_security_domain" "example" {
  local_user_dn  = aci_local_user.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
  description = "from Terraform"
}
```

## Argument Reference ##
* `local_user_dn` - (Required) Distinguished name of parent LocalUser object.
* `name` - (Required) Name of object User Security Domain.
* `annotation` - (Optional) Annotation of object User Security Domain.
* `name_alias` - (Optional) Name Alias of object User Security Domain.
* `description` - (Optional) Description of object User Security Domain.

## Importing ##
An existing User Security Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_user_security_domain.example <Dn>
```