---
layout: "aci"
page_title: "ACI: aci_active_directory"
sidebar_current: "docs-aci-data-source-active_directory"
description: |-
  Data source for ACI Active Directory
---

# aci_active_directory #

Data source for ACI Active Directory


## API Information ##

* `Class` - cloudAD
* `Distinguished Name` - uni/tn-{name}/ad-{id}

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_active_directory" "example" {
  tenant_dn  = aci_tenant.example.id
  active_directory_id  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `active_directory_id` - (Required) Active_directory_id of object Active Directory.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Active Directory.
* `annotation` - (Optional) Annotation of object Active Directory.
* `name_alias` - (Optional) Name Alias of object Active Directory.
* `active_directory_id` - (Optional) AD ID. An object identifier.
