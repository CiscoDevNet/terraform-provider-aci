---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_duo_provider_group"
sidebar_current: "docs-aci-data-source-duo_provider_group"
description: |-
  Data source for ACI Duo Provider Group
---

# aci_duo_provider_group #
Data source for ACI Duo Provider Group


## API Information ##
* `Class` - aaaDuoProviderGroup
* `Distinguished Name` - uni/userext/duoext/duoprovidergroup-{name}

## Example Usage ##

```hcl
data "aci_duo_provider_group" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object Duo Provider Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Duo Provider Group.
* `annotation` - (Optional) Annotation of object Duo Provider Group.
* `auth_choice` - (Optional) Authentication choice of object Duo Provider Group. 
* `provider_type` - (Optional) Type of the Auth Provider. 
* `ldap_group_map_ref` - (Optional) Reference to LDAP Group Map containing user's group membership info.
* `sec_fac_auth_methods` - (Optional) Second factor authentication methods. 
* `name_alias` - (Optional) Name alias of object Duo Provider Group.
* `description` - (Optional) Description of object Duo Provider Group.
