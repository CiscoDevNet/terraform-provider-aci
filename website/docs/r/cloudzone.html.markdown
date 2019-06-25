---
layout: "aci"
page_title: "ACI: aci_cloud_availability_zone"
sidebar_current: "docs-aci-resource-cloud_availability_zone"
description: |-
  Manages ACI Cloud Availability Zone
---

# aci_cloud_availability_zone #
Manages ACI Cloud Availability Zone
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
	resource "aci_cloud_availability_zone" "foocloud_availability_zone" {
		cloud_providers_region_dn = "${aci_cloud_providers_region.example.id}"
		description               = "sample aws availability zone"
		name                      = "us-east-1a"
		annotation                = "tag_zone_a"
		name_alias                = "alias_zone"
	}
```
## Argument Reference ##
* `cloud_providers_region_dn` - (Required) Distinguished name of parent CloudProvidersRegion object.
* `name` - (Required) name of Object cloud_availability_zone. Should match the Availability zone name in AWS cloud.
* `annotation` - (Optional) annotation for object cloud_availability_zone.
* `name_alias` - (Optional) name_alias for object cloud_availability_zone.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Availability Zone.

## Importing ##

An existing Cloud Availability Zone can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_availability_zone.example <Dn>
```