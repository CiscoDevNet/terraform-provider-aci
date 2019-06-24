---
layout: "aci"
page_title: "ACI: aci_cloud_cidr_pool"
sidebar_current: "docs-aci-resource-cloud_cidr_pool"
description: |-
  Manages ACI Cloud CIDR Pool
---

# aci_cloud_cidr_pool #
Manages ACI Cloud CIDR Pool
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
	resource "aci_cloud_cidr_pool" "foocloud_cidr_pool" {
		cloud_context_profile_dn = "${aci_cloud_context_profile.foocloud_context_profile.id}"
		description              = "cloud CIDR"
		addr                     = "10.0.1.10/28"
		annotation               = "tag_cidr"
		name_alias               = "%s"
		primary                  = "yes"
	}
```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `addr` - (Required) CIDR IPv4 block.
* `annotation` - (Optional) annotation for object cloud_cidr_pool.
* `name_alias` - (Optional) name_alias for object cloud_cidr_pool.
* `primary` - (Optional) Flag to specify whether CIDR is primary CIDR or not. Allowed values are "yes" and "no". Default is "no". Only one primary CIDR is supported under a cloud context profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud CIDR Pool.

## Importing ##

An existing Cloud CIDR Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_cidr_pool.example <Dn>
```