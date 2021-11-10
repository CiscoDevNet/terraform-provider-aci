---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selector"
sidebar_current: "docs-aci-resource-cloud_endpoint_selector"
description: |-
  Manages ACI Cloud Endpoint Selector
---

# aci_cloud_endpoint_selector

Manages ACI Cloud Endpoint Selector  
<b>Note: This resource is supported in Cloud APIC only.</b>

## Example Usage

```hcl
	resource "aci_cloud_endpoint_selector" "foocloud_endpoint_selector" {
		cloud_epg_dn     = aci_cloud_epg.foocloud_epg.id
		description      = "sample ep selector"
		name             = "ep_select"
		annotation       = "tag_ep"
		match_expression = "custom:Name=='admin-ep2'"
		name_alias       = "alias_ep"
	}
```

## Argument Reference

- `cloud_epg_dn` - (Required) Distinguished name of parent Cloud EPg object.
- `name` - (Required) Name of Object cloud endpoint selector.
- `annotation` - (Optional) Annotation for object cloud endpoint selector.
- `description` - (Optional) Description for object cloud endpoint selector.
- `match_expression` - (Optional) Match expression for the endpoint selector to select EP on criteria.
- `name_alias` - (Optional) Name alias for object cloud endpoint selector.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Endpoint Selector.

## Importing

An existing Cloud Endpoint Selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_endpoint_selector.example <Dn>
```
