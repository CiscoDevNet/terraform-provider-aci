---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_ldap_group_map_rule_to_group_map"
sidebar_current: "docs-aci-data-source-ldap_group_map_rule_to_group_map"
description: |-
  Data source for ACI LDAP Group Map Rule to Group Map Ref
---

# aci_ldap_group_map_rule_to_group_map #
Data source for ACI LDAP Group Map Rule to Group Map Ref


## API Information ##
* `Class` - aaaLdapGroupMapRuleRef
* `Distinguished Name` - uni/userext/ldapext/ldapgroupmap-{name}/ldapgroupmapruleref-{name} & uni/userext/ldapext/ldapgroupmap-{name}/ldapgroupmapruleref-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> LDAP -> LDAP Group Maps -> Rules & Admin -> AAA -> Authentication -> DUO -> LDAP -> Group Maps -> Rules

## Example Usage ##
```hcl
data "aci_ldap_group_map_rule_to_group_map" "example" {
  ldap_group_map_dn  = aci_ldap_group_map.example.id
  name  = "example"
}
```

## Argument Reference ##
* `ldap_group_map_dn` - (Required) Distinguished name of parent LDAP Group Map object.
* `name` - (Required) Name of object LDAP Group Map Rule to Group Map Ref.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the LDAP Group Map Rule to Group Map Ref.
* `annotation` - (Optional) Annotation of object LDAP Group Map Rule to Group Map Ref.
* `name_alias` - (Optional) Name Alias of object LDAP Group Map Rule to Group Map Ref.
* `description` - (Optional) Description of object LDAP Group Map Rule to Group Map Ref.
