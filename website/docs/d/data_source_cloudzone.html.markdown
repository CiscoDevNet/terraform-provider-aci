---
layout: "aci"
page_title: "ACI: aci_cloud_availability_zone"
sidebar_current: "docs-aci-data-source-cloud_availability_zone"
description: |-
  Data source for ACI Cloud Availability Zone
---

# aci_cloud_availability_zone #
Data source for ACI Cloud Availability Zone
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_availability_zone" "example" {

  cloud_providers_region_dn  = "${aci_cloud_providers_region.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `cloud_providers_region_dn` - (Required) Distinguished name of parent CloudProvidersRegion object.
* `name` - (Required) name of Object cloud_availability_zone.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Availability Zone.
* `annotation` - (Optional) annotation for object cloud_availability_zone.
* `name_alias` - (Optional) name_alias for object cloud_availability_zone.
