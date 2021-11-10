---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selector"
sidebar_current: "docs-aci-data-source-cloud_endpoint_selector"
description: |-
  Data source for ACI Cloud Endpoint Selector
---

# aci_cloud_endpoint_selector

Data source for ACI Cloud Endpoint Selector  
<b>Note: This resource is supported in Cloud APIC only.</b>

## Example Usage

```hcl
data "aci_cloud_endpoint_selector" "dev_ep_select" {
  cloud_epg_dn  = aci_cloud_epg.dev_epg.id
  name          = "dev_ep_select"
}
```

## Argument Reference

- `cloud_epg_dn` - (Required) Distinguished name of parent Cloud EPg object.
- `name` - (Required) Name of Object cloud endpoint selector.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Cloud Endpoint Selector.
- `annotation` - (Optional) Annotation for object cloud endpoint selector.
- `description` - (Optional) Description for object cloud endpoint selector.
- `match_expression` - (Optional) Match expression for the endpoint selector to select EP on criteria.
- `name_alias` - (Optional) Name alias for object cloud endpoint selector.
