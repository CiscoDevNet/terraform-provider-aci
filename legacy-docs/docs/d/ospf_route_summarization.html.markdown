---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_ospf_route_summarization"
sidebar_current: "docs-aci-data-source-ospf_route_summarization"
description: |-
  Data source for ACI OSPF Route Summarization
---

# aci_ospf_route_summarization

Data source for ACI OSPF Route Summarization

## Example Usage

```hcl
data "aci_ospf_route_summarization" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "ospf_route_summarization_1"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of object OSPF route summarization.

## Attribute Reference

- `id` - Attribute ID set to the Dn of the OSPF Route Summarization.
- `annotation` - (Optional) Annotation for object OSPF route summarization.
- `cost` - (Optional) The OSPF Area cost for the default summary LSAs. The Area cost is used with NSSA and stub area types only.
- `inter_area_enabled` - (Optional) Inter area enabled flag for object OSPF route summarization.
- `name_alias` - (Optional) Name alias for object OSPF route summarization.
- `tag` - (Optional) The color of a policy label.
- `description` - Description for the object of the OSPF Route Summarization.
