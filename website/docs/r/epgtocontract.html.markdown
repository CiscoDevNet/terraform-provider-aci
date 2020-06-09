---
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-resource-epg_to_contract"
description: |-
  Manages  ACI EPG to contract relationship.
---

# aci_epg_to_contract #
Manages ACI EPG to contract relationship.

## Example Usage ##

```hcl
resource "aci_epg_to_contract" "example" {
    application_epg_dn = "${aci_application_epg.demo.id}"
    contract_dn  = "${aci_contract.demo_contract.id}"
    contract_type = "consumer"
}
```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of Parent epg.
* `contract_dn` - (Required) Distinguished name of contract to attach.
* `contract_type` - (Required) Type of relationship. Allowed values are `consumer` and `provider`.
* `annotation` - (Optional) annotation for object.
* `match_t` - (Optional) Provider matching criteria.
* `prio` - (Optional) Priority of relation object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the provider/consumer contract.

