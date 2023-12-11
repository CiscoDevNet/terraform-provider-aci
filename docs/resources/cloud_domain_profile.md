---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_domain_profile"
sidebar_current: "docs-aci-resource-aci_cloud_domain_profile"
description: |-
  Manages Cloud Network Controller Cloud Domain Profile
---

# aci_cloud_domain_profile

Manages Cloud Network Controller Cloud Domain Profile  
<b>Note: This resource is supported in Cloud Network Controller only.</b>

## Example Usage

```hcl
resource "aci_cloud_domain_profile" "foocloud_domain_profile" {
  annotation  = "tag_domp"
  name_alias  = "alias_domp"
  description = "from terraform"
  site_id     = "0"
}

```

## Argument Reference

- `annotation` - (Optional) Specifies the Annotation of the cloud domain profile.
- `description` - (Optional) Specifies the Description of the cloud domain profile.
- `name_alias` - (Optional) Specifies the alias-name of the cloud domain profile.
- `site_id` - (Optional) Site-ID of the cloud domain profile. The allowed value range is "0" to "1000". Default is "0".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Domain Profile.

## Importing

An existing Cloud Domain Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_domain_profile.example <Dn>
```
