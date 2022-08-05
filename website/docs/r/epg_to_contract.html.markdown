---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-resource-epg_to_contract"
description: |-
  Manages  ACI EPG to contract relationship.
---

# aci_epg_to_contract

Manages ACI EPG to contract relationship.

## Example Usage

```hcl
resource "aci_epg_to_contract" "example" {
  application_epg_dn = aci_application_epg.demo.id
  contract_dn        = aci_contract.demo_contract.id
  contract_type      = "provider"
  annotation         = "terraform"
  match_t            = "AtleastOne"
  prio               = "unspecified"
}
```

## Argument Reference

- `application_epg_dn` - (Required) Distinguished name of Parent epg.
- `contract_dn` - (Required) Distinguished name of contract to attach.
- `contract_type` - (Required) Type of relationship. Allowed values are "consumer" and "provider".
- `annotation` - (Optional) Annotation for EPg to contract relationship.
- `match_t` - (Optional) Provider matching criteria. Allowed values: "All", "AtleastOne", "AtmostOne", "None". Default value: "AtleastOne". This attribute is supported only for resources with "contract_type" with value "provider"
- `prio` - (Optional) Priority of relation object. Allowed values: "unspecified", "level1", "level2", "level3", "level4", "level5", "level6". Default value: "unspecified".

## Attribute Reference

- `id` - Attribute id set to the Dn of the provider/consumer contract.
