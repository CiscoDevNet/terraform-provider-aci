---
subcategory: "Monitoring"
layout: "aci"
page_title: "ACI: aci_span_source_group"
sidebar_current: "docs-aci-resource-span_source_group"
description: |-
  Manages ACI SPAN Source Group
---

# aci_span_source_group

Manages ACI SPAN Source Group

## Example Usage

```hcl
resource "aci_span_source_group" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "example"
  admin_st    = "enabled"
  annotation  = "tag_span"
  description = "from terraform"
  name_alias  = "alias_span"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object SPAN Source Group.
- `admin_st` - (Optional) administrative state of the object or policy.
  Allowed values: "enabled", "disabled". Default value: "enabled".
- `annotation` - (Optional) Specifies a annotation of the policy definition.
- `description` - (Optional) Specifies a description of the policy definition.
- `name_alias` - (Optional) Specifies a alias-name of the policy definition.

- `relation_span_rs_src_grp_to_filter_grp` - (Optional) Relation to class spanFilterGrp. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the SPAN Source Group.

## Importing

An existing SPAN Source Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_span_source_group.example <Dn>
```
