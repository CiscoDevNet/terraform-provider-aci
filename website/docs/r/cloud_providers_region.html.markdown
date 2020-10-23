---
layout: "aci"
page_title: "ACI: aci_cloud_providers_region"
sidebar_current: "docs-aci-resource-cloud_providers_region"
description: |-
  Manages ACI Cloud Providers Region
---

# aci_cloud_providers_region #
Manages ACI Cloud Providers Region
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
	resource "aci_cloud_providers_region" "foocloud_providers_region" {
		cloud_provider_profile_dn = "${aci_cloud_provider_profile.example.id}"
		description               = "aws region"
		name                      = "us-east-1"
		annotation                = "tag_region"
		name_alias                = "default_reg"
	}
```
## Argument Reference ##
* `cloud_provider_profile_dn` - (Required) Distinguished name of parent CloudProviderProfile object.
* `name` - (Required) name of Object cloud_providers_region.
* `annotation` - (Optional) annotation for object cloud_providers_region.
* `name_alias` - (Optional) name_alias for object cloud_providers_region.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Providers Region.

## Importing ##

An existing Cloud Providers Region can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_providers_region.example <Dn>
```