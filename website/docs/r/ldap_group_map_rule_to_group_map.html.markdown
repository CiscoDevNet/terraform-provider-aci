---
layout: "aci"
page_title: "ACI: aci_ldap_group_map_rule_to_group_map"
sidebar_current: "docs-aci-resource-ldap_group_map_rule_to_group_map"
description: |-
  Manages ACI LDAP Group Map Rule to Group Map Ref
---

# aci_ldap_group_map_rule_to_group_map #
Manages ACI LDAP Group Map Rule to Group Map Ref

## API Information ##
* `Class` - aaaLdapGroupMapRuleRef
* `Distinguished Named` - uni/userext/ldapext/ldapgroupmap-{name}/ldapgroupmapruleref-{name} & uni/userext/ldapext/ldapgroupmap-{name}/ldapgroupmapruleref-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> LDAP -> LDAP Group Maps -> Rules & Admin -> AAA -> Authentication -> DUO -> LDAP -> Group Maps -> Rules


## Example Usage ##
```hcl
resource "aci_ldap_group_map_rule_to_group_map" "example" {
  ldap_group_map_dn  = aci_ldap_group_map.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias_value"
  description = "from terraform"
}
```

## Argument Reference ##
* `ldap_group_map_dn` - (Required) Distinguished name of parent LDAP Group Map object.
* `name` - (Required) Name of object LDAP Group Map Rule to Group Map Ref.
* `annotation` - (Optional)  Annotation of object LDAP Group Map Rule to Group Map Ref.
* `name_alias` - (Optional) Name Alias of object LDAP Group Map Rule to Group Map Ref.
* `description` - (Optional) Description of object LDAP Group Map Rule to Group Map Ref.

## Importing ##

An existing LDAP Group Map Rule to Group Map Ref can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ldap_group_map_rule_to_group_map.example <Dn>
```