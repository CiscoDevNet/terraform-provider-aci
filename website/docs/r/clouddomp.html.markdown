---
layout: "aci"
page_title: "ACI: aci_cloud_domain_profile"
sidebar_current: "docs-aci-resource-cloud_domain_profile"
description: |-
  Manages ACI Cloud Domain Profile
---

# aci_cloud_domain_profile #
Manages ACI Cloud Domain Profile
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
resource "aci_cloud_domain_profile" "foocloud_domain_profile" {
  annotation  = "tag_domp"
  name_alias  = "alias_domp"
  site_id     = "0"
}

```
## Argument Reference ##
* `annotation` - (Optional) annotation for object cloud_domain_profile.
* `name_alias` - (Optional) name_alias for object cloud_domain_profile.
* `site_id` - (Optional) site_id for object cloud_domain_profile. Allowed value range is "0" to "1000". Default is "0".



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Domain Profile.

## Importing ##

An existing Cloud Domain Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_domain_profile.example <Dn>
```