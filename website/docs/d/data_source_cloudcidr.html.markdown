---
layout: "aci"
page_title: "ACI: aci_cloud_cidr_pool"
sidebar_current: "docs-aci-data-source-cloud_cidr_pool"
description: |-
  Data source for ACI Cloud CIDR Pool
---

# aci_cloud_cidr_pool #
Data source for ACI Cloud CIDR Pool
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_cidr_pool" "example" {

  cloud_context_profile_dn  = "${aci_cloud_context_profile.example.id}"

  addr  = "example"
}
```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `addr` - (Required) addr of Object cloud_cidr_pool.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud CIDR Pool.
* `addr` - (Optional) peer address
* `annotation` - (Optional) annotation for object cloud_cidr_pool.
* `name_alias` - (Optional) name_alias for object cloud_cidr_pool.
* `primary` - (Optional) primary for object cloud_cidr_pool.
