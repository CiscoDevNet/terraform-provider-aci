---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-data-source-aci_contract_subject"
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

- `contract_dn` - (Required) Distinguished name of parent Contract object.
- `name` - (Required) name of Object contract_subject.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Contract Subject.
- `annotation` - (Optional) Annotation for the contract subject object.
- `description` - (Optional) Description for the contract subject object.
- `cons_match_t` - (Optional) The subject match criteria across consumers.
- `name_alias` - (Optional) Name alias for the contract subject object.
- `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
- `prov_match_t` - (Optional) The subject match criteria across providers.
- `rev_flt_ports` - (Optional) Enables filter to apply on ingress and egress traffic.
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.

- `apply_both_directions` - (Optional) Defines if a contract subject allows traffic matching the defined filters for both consumer-to-provider and provider-to-consumer together or if each direction can be defined separately. By default set to "yes".
  - When set to "yes", the filters defined in the subject are applied for both directions.
  - When set to "no", the filters for each direction (consumer-to-provider and provider-to-consumer) are defined independently.
- ` consumer_to_provider` - (Optional) Filter Chain for Consumer to Provider. Class vzInTerm. Type - Block.
    - `prio` - (Optional) The priority level to assign to traffic matching the consumer to provider flows.
    - `target_dscp` - (Optional) The target DSCP to assign to traffic matching the consumer to provider flows.
	- `relation_vz_rs_in_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).
	-`relation_vz_rs_filt_att` - (Optional) Relation to class Filters (vzRsFiltAtt).
      - `action` - (Optional) The action required when the condition is met. Allowed values are "deny", "permit", and the default value is "permit".
      - `directives` - (Optional) Directives of the Contract Subject Filter object for Consumer to Provider. Allowed values are "log", "no_stats", "none", and the default value is "none".
      - `priority_override` - (Optional) Priority override of the Consumer to Provider Filter object. It is only used when action is deny. Allowed values are "default", "level1", "level2", "level3", and the default value is "default".
      - `filter_dn` - (Required) Distinguished name of the Filter object for Consumer to Provider.
- `provider_to_consumer` - (Optional) Filter Chain For Provider to Consumer. Class vzOutTerm. Type - Block.
    - `prio` - (Optional) The priority level to assign to traffic matching the provider to consumer flows.
    - `target_dscp` - (Optional) The target DSCP to assign to traffic matching the provider to consumer flows.
	- `relation_vz_rs_out_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).
	-`relation_vz_rs_filt_att` - (Optional) (Optional) Relation to class Filters (vzRsFiltAtt).
      - `action` - (Optional) The action required when the condition is met. Allowed values are "deny", "permit", and the default value is "permit".
      - `directives` - (Optional) Directives of the Contract Subject Filter object for Provider to Consumer. Allowed values are "log", "no_stats", "none", and the default value is "none".
      - `priority_override` - (Optional) Priority override of the Provider to Consumer Filter object. It is only used when action is deny. Allowed values are "default", "level1", "level2", "level3", and the default value is "default".
      - `filter_dn` - (Required) Distinguished name of the Filter object for Provider to Consumer.
