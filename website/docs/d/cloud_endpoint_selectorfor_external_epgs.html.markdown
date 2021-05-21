---
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selectorfor_external_epgs"
sidebar_current: "docs-aci-data-source-cloud_endpoint_selectorfor_external_epgs"
description: |-
  Data source for ACI Cloud Endpoint Selector for External EPgs
---

# aci_cloud_endpoint_selectorfor_external_epgs

Data source for ACI Cloud Endpoint Selector for External EPgs  
<b>Note: This resource is supported in Cloud APIC only.</b>

## Example Usage

```hcl
data "aci_cloud_endpoint_selectorfor_external_epgs" "foo_ep_selector" {
  cloud_external_epg_dn  = aci_cloud_external_epg.ext_epg.id
  name                   = dev_ext_ep_select
  subnet                 = "10.0.0.1/16"
}
```

## Argument Reference

- `cloud_external_epg_dn` - (Required) Distinguished name of parent Cloud External EPg object.
- `name` - (Required) Name of Object cloud endpoint selector for external EPgs.
- `subnet` - (Required) Subnet from which EP to select. Any valid CIDR block is allowed here.


## Attribute Reference

- `id` - Attribute id set to the Dn of the Cloud Endpoint Selector for External EPgs.
- `annotation` - (Optional) Annotation for object cloud endpoint selector for external EPgs.
- `description` - (Optional) Description for object cloud endpoint selector for external EPgs.
- `is_shared` - (Optional) For Selectors set the shared route control.
- `name_alias` - (Optional) Name alias for object cloud endpoint selector for external EPgs.
