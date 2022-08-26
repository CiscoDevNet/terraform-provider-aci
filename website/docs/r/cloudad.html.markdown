---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_ad"
sidebar_current: "docs-aci-resource-cloud_ad"
description: |-
  Manages ACI Cloud Active Directory
---

# aci_cloud_ad #

Manages ACI Cloud Active Directory
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudAD
* `Distinguished Name` - uni/tn-{name}/ad-{id}

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> Tenants  -> {tenant_name}


## Example Usage ##

```hcl
resource "aci_cloud_ad" "example" {
  tenant_dn  = aci_tenant.example.id
  active_directory_id  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `active_directory_id` - (Required) Id of the Azure Active Directory.
* `annotation` - (Optional) Annotation of the object Active Directory.
* `name` - (Optional) Name of the Active Directory object.


## Importing ##

An existing Active Directory can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_ad.example <Dn>
```