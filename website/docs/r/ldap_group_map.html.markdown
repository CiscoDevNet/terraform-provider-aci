---
layout: "aci"
page_title: "ACI: aci_ldap_group_map"
sidebar_current: "docs-aci-resource-ldap_group_map"
description: |-
  Manages ACI LDAP Group Map
---

# aci_ldap_group_map #

Manages ACI LDAP Group Map

## API Information ##

* `Class` - aaaLdapGroupMap
* `Distinguished Named` - uni/userext/ldapext/ldapgroupmap-{name} and uni/userext/duoext/ldapgroupmap-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_ldap_group_map" "example" {
  name        = "example"
  type        = "ldap"
  name_alias  = "demo_ldap_group_map"
  description = "From Terraform"
  annotation  = "orchestrator:terraform"
}
```

## Argument Reference ##


* `name` - (Required) Name of object LDAP Group Map.
* `type` - (Required) Type of object LDAP Group Map. Allowed values are "duo" and "ldap".
* `annotation` - (Optional) Annotation of object LDAP Group Map.
* `description` - (Optional) Description of object LDAP Group Map.
* `name_alias` - (Optional) Annotation of object LDAP Group Map.



## Importing ##

An existing LDAPGroupMap can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ldap_group_map.example <Dn>
```