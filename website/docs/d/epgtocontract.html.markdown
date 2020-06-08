---
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-data-source-epg_to_contract"
description: |-
  Data source for ACI EPG to contract relationship.
---

# aci_epg_to_contract #
Data source for ACI EPG to contract relationship.

## Example Usage ##

```hcl
data "aci_epg_to_contract" "example" {
    application_epg_dn = "${aci_application_epg.demo.id}"
    contract_name  = "example"
    contract_type = "consumer"
}
```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of Parent epg.
* `contract_name` - (Required) Name of the contract.
* `contract_type` - (Required) Type of relationship. Allowed values are `consumer` and `provider`.



## Attribute Reference

* `id` - Attribute id set to the Dn of the provider/consumer contract.
* `annotation` - (Optional) annotation for object.
* `match_t` - (Optional) Provider matching criteria.
* `prio` - (Optional) Priority of relation object.
