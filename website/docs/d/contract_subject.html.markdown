---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-data-source-contract_subject"
description: |-
  Data source for ACI Contract Subject
---

# aci_contract_subject

Data source for ACI Contract Subject

## Example Usage

```hcl
data "aci_contract_subject" "dev_subject" {
  contract_dn  = aci_contract.example.id
  name         = "foo_subject"
}
```

## Argument Reference

- `contract_dn` - (Required) Distinguished name of parent (Contract) object.
- `name` - (Required) name of Object contract_subject.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Contract Subject.
- `annotation` - (Optional) Annotation for the object contract subject.
- `description` - (Optional) Description for the object contract subject.
- `cons_match_t` - (Optional) The subject match criteria across consumers.
- `name_alias` - (Optional) Name alias for the object contract subject.
- `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
- `prov_match_t` - (Optional) The subject match criteria across consumers.
- `rev_flt_ports` - (Optional) Enables filter to apply on ingress and egress traffic.
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.

- `apply_both_directions` - (Optional) . By default set to "yes".
- ` consumer_to_provider` - (Optional) Filter Chain For Consumer to Provider. Class vzInTerm.
    - `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
    - `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
    - `relation_vz_rs_in_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).
- `provider_to_consumer` - (Optional) Filter Chain For Provider to Consumer. Class vzOutTerm
    - `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
    - `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
    - `relation_vz_rs_out_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).
