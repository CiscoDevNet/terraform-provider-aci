---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_ldap_group_map_rule"
sidebar_current: "docs-aci-resource-ldap_group_map_rule"
description: |-
  Manages ACI LDAP Group Map Rule
---

# aci_ldap_group_map_rule #

Manages ACI LDAP Group Map Rule

## API Information ##

* `Class` - aaaLdapGroupMapRule
* `Distinguished Name` - uni/userext/duoext/ldapgroupmaprule-{name}
                          uni/userext/ldapext/ldapgroupmaprule-{name}
## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> LDAP -> LDAP Group Map Rules & <br>Admin -> AAA -> Authentication -> DUO -> LDAP -> Group Map Rules

## Example Usage ##

```hcl
resource "aci_ldap_group_map_rule" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  groupdn     = "groupdn_example"
  type        = "duo"
  description = "From Terraform"
  name_alias  = "ldap_group_map_rule_alias"
}
```

## Argument Reference ##


* `name` - (Required) Name of object LDAP Group Map Rule.
* `type` - (Required) Type of object LDAP Group Map Rule. Allowed values are "duo" and "ldap". Type: String.
* `annotation` - (Optional) Annotation of object LDAP Group Map Rule.
* `description` - (Optional) Description for object LDAP Group Map Rule. Type: String.
* `groupdn` - (Optional) LDAP Group DN to compare with LDAP search query for user's membership. Type: String.
* `name_alias` - (Optional) Name Alias of object LDAP Group Map Rule. Type: String.

## Importing ##

An existing LDAPGroupMapRule can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ldap_group_map_rule.example <Dn>
```