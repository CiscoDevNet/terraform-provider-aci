---
layout: "aci"
page_title: "ACI: aci_active_directory"
sidebar_current: "docs-aci-resource-active_directory"
description: |-
  Manages ACI Active Directory
---

# aci_active_directory #

Manages ACI Active Directory

## API Information ##

* `Class` - cloudAD
* `Distinguished Name` - uni/tn-{name}/ad-{id}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_active_directory" "example" {
  tenant_dn  = aci_tenant.example.id
  active_directory_id  = "example"
  active_directory_id = 

}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `active_directory_id` - (Required) Active_directory_id of the object Active Directory.
* `annotation` - (Optional) Annotation of the object Active Directory.
* `active_directory_id` - (Optional) AD ID.An object identifier.


## Importing ##

An existing ActiveDirectory can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_active_directory.example <Dn>
```