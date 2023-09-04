---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-data-source-epg_to_contract"
description: |-
  Data source for ACI EPG to contract relationship.
---

# aci_epg_to_contract

Data source for ACI EPG to contract relationship.

## Example Usage

```hcl
data "aci_epg_to_contract" "example" {
  application_epg_dn = aci_application_epg.demo.id
  contract_dn        = "example"
  contract_type      = "consumer"
}
```

## Argument Reference

- `application_epg_dn` - (Required) Distinguished name of Parent EPG.
- `contract_dn` - (Required) Distinguished name of the Contract object to attach.
- `contract_type` - (Required) Type of the EPG to contract relationship object. Allowed values are "consumer" and "provider".

## Attribute Reference

- `id` - Attribute id set to the Dn of the provider/consumer contract.
- `annotation` - (Read-Only) Annotation of the EPG to contract relationship object.
- `match_t` - (Read-Only) Matching criteria of the EPG to contract relationship object, only supported for `contract_type` "provider".
- `prio` - (Read-Only) Priority of the EPG to contract relationship object.
