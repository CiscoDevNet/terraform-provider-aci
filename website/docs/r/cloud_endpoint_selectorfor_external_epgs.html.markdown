---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selectorfor_external_epgs"
sidebar_current: "docs-aci-resource-cloud_endpoint_selectorfor_external_epgs"
description: |-
  Manages Cloud Network Controller Cloud Endpoint Selector for External EPgs
---

# aci_cloud_endpoint_selectorfor_external_epgs

Manages Cloud Network Controller Cloud Endpoint Selector for External EPgs  
<b>Note: This resource is supported in Cloud Network Controller only.</b>

## Example Usage

```hcl
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "foocloud_endpoint_selectorfor_external_epgs" {
		cloud_external_epg_dn = aci_cloud_external_epg.foocloud_external_epg.id
		subnet                = "0.0.0.0/0"
		description           = "sample external ep selector"
		name                  = "ext_ep_selector"
		annotation            = "tag_ext_selector"
		is_shared             = "yes"
		match_expression 	  = "custom:tag=='provbaz'"
		name_alias            = "alias_select"
	}
```

## Argument Reference

- `cloud_external_epg_dn` - (Required) Distinguished name of parent Cloud External EPg object.
- `name` - (Required) Name of Object cloud endpoint selector for external EPgs.
- `subnet` - (Required) Subnet from which EP to select. Any valid CIDR block is allowed here.
- `match_expression` - (Optional) Expressions are not used in cloudExtEPSelector because this selector only match subnets.
- `annotation` - (Optional) Annotation for object cloud endpoint selector for external EPgs.
- `description` - (Optional) Description for object cloud endpoint selector for external EPgs.
- `is_shared` - (Optional) For Selectors set the shared route control. Allowed values are "yes" and "no". Default value is "yes".
- `name_alias` - (Optional) Name alias for object cloud endpoint selector for external EPgs.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Endpoint Selector for External EPgs.

## Importing

An existing Cloud Endpoint Selector for External EPgs can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_endpoint_selectorfor_external_epgs.example <Dn>
```
