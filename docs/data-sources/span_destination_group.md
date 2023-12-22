---
subcategory: "Monitoring"
layout: "aci"
page_title: "ACI: aci_span_destination_group"
sidebar_current: "docs-aci-data-source-aci_span_destination_group"
description: |-
  Data source for ACI SPAN Destination Group
---

# aci_span_destination_group

Data source for ACI SPAN Destination Group

## Example Usage

```hcl
data "aci_span_destination_group" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of object SPAN destination group.

## Attribute Reference

- `id` - Attribute id set to the Dn of the SPAN Destination Group.
- `annotation` - (Optional) Annotation of object SPAN destination group.
- `description` - (Optional) Specifies a description of the policy definition.
- `name_alias` - (Optional) Name alias of object SPAN destination group.
