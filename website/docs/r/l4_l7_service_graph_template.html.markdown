---
layout: "aci"
page_title: "ACI: aci_l4_l7_service_graph_template"
sidebar_current: "docs-aci-resource-l4_l7_service_graph_template"
description: |-
  Manages ACI L4-L7 Service Graph Template
---

# aci_l4_l7_service_graph_template

Manages ACI L4-L7 Service Graph Template

## Example Usage

```hcl

resource "aci_l4_l7_service_graph_template" "example" {
  tenant_dn                         = aci_tenant.tenentcheck.id
  name                              = "first"
  annotation                        = "example"
  name_alias                        = "example"
  description                       = "from terraform"
  l4_l7_service_graph_template_type = "legacy"
  ui_template_type                  = "ONE_NODE_ADC_ONE_ARM"
}

```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name for L4-L7 service graph template object.
- `annotation` - (Optional) Annotation for L4-L7 service graph template object.
- `description` - (Optional) Description for L4-L7 service graph template object.
- `name_alias` - (Optional) Name alias for L4-L7 service graph template object.
- `l4-l7_service_graph_template_type` - (Optional) Component type for L4-L7 service graph template object. Allowed values are "legacy" and "cloud". Default value is "legacy".
- `ui_template_type` - (Optional) UI template type for L4-L7 service graph template object. Allowed values are "ONE_NODE_ADC_ONE_ARM", "ONE_NODE_ADC_ONE_ARM_L3EXT", "ONE_NODE_ADC_TWO_ARM", "ONE_NODE_FW_ROUTED", "ONE_NODE_FW_TRANS", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM_L3EXT", "TWO_NODE_FW_ROUTED_ADC_TWO_ARM", "TWO_NODE_FW_TRANS_ADC_ONE_ARM", "TWO_NODE_FW_TRANS_ADC_ONE_ARM_L3EXT", "TWO_NODE_FW_TRANS_ADC_TWO_ARM" and "UNSPECIFIED". Default value is "UNSPECIFIED".
- `term_cons_name` - (Optional) Name of consumer terminal node. Default value is "T1".
- `term_prov_name` - (Optional) Name of provider terminal node. Default value is "T2".

## Attribute Reference

- `id` - Dn of the L4-L7 service graph template.
- `term_cons_dn` - Dn of consumer terminal node for L4-L7 service graph template.
- `term_prov_dn` - Dn of provider terminal node for L4-L7 service graph template.

## Importing

An existing L4-L7 Service Graph Template can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l4_l7_service_graph_template.example <Dn>
```
