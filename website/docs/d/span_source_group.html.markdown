---
layout: "aci"
page_title: "ACI: aci_span_source_group"
sidebar_current: "docs-aci-data-source-span_source_group"
description: |-
  Data source for ACI SPAN Source Group
---

# aci_span_source_group

Data source for ACI SPAN Source Group

## Example Usage

```hcl
data "aci_span_source_group" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object SPAN Source Group.

## Attribute Reference

- `id` - Attribute id set to the Dn of the SPAN Source Group.
- `admin_st` - (Optional) administrative state of the object or policy.
  Allowed values: "enabled", "disabled". Default value: "enabled".
- `annotation` - (Optional) Specifies a annotation of the policy definition.
- `description` - (Optional) Specifies a description of the policy definition.
- `name_alias` - (Optional) Specifies a alias-name of the policy definition.
