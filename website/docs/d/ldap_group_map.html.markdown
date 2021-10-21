---
layout: "aci"
page_title: "ACI: aci_ldap_group_map"
sidebar_current: "docs-aci-data-source-ldap_group_map"
description: |-
  Data source for ACI LDAP Group Map
---

# aci_ldap_group_map #

Data source for ACI LDAP Group Map


## API Information ##

* `Class` - aaaLdapGroupMap
* `Distinguished Named` - uni/userext/ldapext/ldapgroupmap-{name} and uni/userext/duoext/ldapgroupmap-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> LDAP -> LDAP Group Maps & Admin -> AAA -> Authentication -> DUO -> LDAP -> Group Maps


## Example Usage ##

```hcl
data "aci_ldap_group_map" "example" {
  name  = "example"
  type  = "duo"
}
```

## Argument Reference ##

* `name` - (Required) Name of object LDAP Group Map.
* `type` - (Required) Type of object LDAP Group Map. Allowed values are "duo" and "ldap".

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the LDAP Group Map.
* `annotation` - (Optional) Annotation of object LDAP Group Map.
* `name_alias` - (Optional) Name Alias of object LDAP Group Map.
* `description` - (Optional) Description of object LDAP Group Map.
