---
layout: "aci"
page_title: "ACI: aci_ldap_group_map_rule"
sidebar_current: "docs-aci-data-source-ldap_group_map_rule"
description: |-
  Data source for ACI LDAP Group Map Rule
---

# aci_ldap_group_map_rule #

Data source for ACI LDAP Group Map Rule


## API Information ##

* `Class` - aaaLdapGroupMapRule
* `Distinguished Named` - uni/userext/duoext/ldapgroupmaprule-{name}
                          uni/userext/ldapext/ldapgroupmaprule-{name}
## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> LDAP -> LDAP Group Map Rules & 
               Admin -> AAA -> Authentication -> DUO -> LDAP -> Group Map Rules 


## Example Usage ##

```hcl
data "aci_ldap_group_map_rule" "example" {
  name  = "example"
  type  = "duo"
}
```

## Argument Reference ##

* `name` - (Required) Name of object LDAP Group Map Rule.
* `type` - (Required) Type of object LDAP Group MAp Rule. Allowed Values are "duo" and "ldap".

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the LDAP Group Map Rule.
* `annotation` - (Optional) Annotation of object LDAP Group Map Rule.
* `name_alias` - (Optional) Name Alias of object LDAP Group Map Rule.
* `groupdn` - (Optional) LDAP Group DN to compare with LDAP search query for user's membership. 
* `description` - (Optional) Description of object LDAP Group Map Rule.