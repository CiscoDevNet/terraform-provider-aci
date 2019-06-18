---
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-data-source-contract_subject"
description: |-
  Data source for ACI Contract Subject
---

# aci_contract_subject #
Data source for ACI Contract Subject

## Example Usage ##

```hcl
data "aci_contract_subject" "example" {
  contract_dn  = "${aci_contract.example.id}"
  name         = "example"
}
```
## Argument Reference ##
* `contract_dn` - (Required) Distinguished name of parent Contract object.
* `name` - (Required) name of Object contract_subject.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Contract Subject.
* `annotation` - (Optional) annotation for object contract_subject.
* `cons_match_t` - (Optional) The subject match criteria across consumers.
* `name_alias` - (Optional) name_alias for object contract_subject.
* `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
* `prov_match_t` - (Optional) The subject match criteria across consumers.
* `rev_flt_ports` - (Optional) enables filter to apply on ingress and egress traffic.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
