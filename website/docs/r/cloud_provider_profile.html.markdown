---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_provider_profile"
sidebar_current: "docs-aci-resource-cloud_provider_profile"
description: |-
  Manages ACI Cloud Provider Profile
---

# aci_cloud_provider_profile #
Manages ACI Cloud Provider Profile
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
	resource "aci_cloud_provider_profile" "foocloud_provider_profile" {
		vendor      = "aws"
		annotation  = "tag_aws_prof"
	}
```
## Argument Reference ##
* `vendor` - (Required) vendor of Object cloud_provider_profile. Currently only supported vendor is "aws".
* `annotation` - (Optional) annotation for object cloud_provider_profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Provider Profile.

## Importing ##

An existing Cloud Provider Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_provider_profile.example <Dn>
```