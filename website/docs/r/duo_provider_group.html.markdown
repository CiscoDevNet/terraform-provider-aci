---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_duo_provider_group"
sidebar_current: "docs-aci-resource-duo_provider_group"
description: |-
  Manages ACI Duo Provider Group
---

# aci_duo_provider_group #
Manages ACI Duo Provider Group

## API Information ##
* `Class` - aaaDuoProviderGroup
* `Distinguished Named` - uni/userext/duoext/duoprovidergroup-{name}

## Example Usage ##
```hcl
resource "aci_duo_provider_group" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  auth_choice = "CiscoAVPair"
  ldap_group_map_ref = "DuoEmpGroupMap"
  provider_type = "radius"
  sec_fac_auth_methods = ["auto"]
  name_alias = "example_name_alias"
  description = "from terraform"
}
```

## Argument Reference ##
* `name` - (Required) Name of object Duo Provider Group.
* `annotation` - (Optional) Annotation of object Duo Provider Group.
* `auth_choice` - (Optional) Authentication choice of object Duo Provider Group. Allowed values are "CiscoAVPair" and "LdapGroupMap". Default value is "CiscoAVPair". Type: String.
* `provider_type` - (Optional) Type of the Auth Provider. Allowed values are "ldap" and "radius". Default value is "radius". Type: String.
* `ldap_group_map_ref` - (Optional) Reference to LDAP Group Map containing user's group membership info.
* `sec_fac_auth_methods` - (Optional) Second factor authentication methods. Allowed values are "auto", "passcode", "phone" and "push". Default value is "auto". Type: List.
* `name_alias` - (Optional) Name alias of object Duo Provider Group.
* `description` - (Optional) Description of object Duo Provider Group.


## Importing ##
An existing Duo Provider Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_duo_provider_group.example <Dn>
```