---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-data-source-aci_epg_to_contract"
description: |-
  Data source for ACI EPG to contract relationship.
---

# aci_epg_to_contract

Data source for ACI EPG to contract relationship.

## Example Usage

```hcl
data "aci_epg_to_contract" "example" {
  application_epg_dn = aci_application_epg.example.id
  contract_dn        = aci_contract.example.id
  contract_type      = "consumer"
}
```

## Argument Reference

- `application_epg_dn` - (Required) Distinguished name of Parent EPG. Type: String.
- `contract_name` - (Deprecated) Name of the Contract object to attach. Type: String.
- `contract_dn` - (Optional) Distinguished name of the Contract object to attach. Type: String.
- `contract_type` - (Required) Type of the EPG to contract relationship object. Allowed values are "consumer" and "provider". Type: String.

## Attribute Reference

- `id` - Attribute id set to the Dn of the provider/consumer contract. Type: String.
- `annotation` - (Read-Only) Annotation of the EPG to contract relationship object. Type: String.
- `match_t` - (Read-Only) Matching criteria of the EPG to contract relationship object, only supported for `contract_type` "provider". Type: String.
- `prio` - (Read-Only) Priority of the EPG to contract relationship object. Type: String.
