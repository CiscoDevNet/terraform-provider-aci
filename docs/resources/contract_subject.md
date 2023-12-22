---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-resource-aci_contract_subject"
description: |-
  Manages ACI Contract Subject
---

# aci_contract_subject

Manages ACI Contract Subject

## Example Usage

```hcl
	resource "aci_contract_subject" "foocontract_subject" {
		contract_dn   = aci_contract.example.id
		description   = "from terraform"
		name          = "demo_subject"
		annotation    = "tag_subject"
		cons_match_t  = "AtleastOne"
		name_alias    = "alias_subject"
		prio          = "level1"
		prov_match_t  = "AtleastOne"
		rev_flt_ports = "yes"
		target_dscp   = "CS0"
	}

	resource "aci_contract_subject" "foocontract_subject_2" {
		contract_dn   = aci_contract.example.id
		description   = "from terraform"
		name          = "demo_subject"
		annotation    = "tag_subject"
		cons_match_t  = "AtleastOne"
		name_alias    = "alias_subject"
		prio          = "level1"
		prov_match_t  = "AtleastOne"
		rev_flt_ports = "no"
		target_dscp   = "CS0"
		apply_both_directions = "no"
		consumer_to_provider = {
			prio                             = "level2"
			target_dscp                      = "AF41"
			relation_vz_rs_in_term_graph_att = aci_l4_l7_service_graph_template.service_graph.id
			relation_vz_rs_filt_att {
				action = "deny"
				directives = ["log", "no_stats"]
				priority_override = "level2"
				filter_dn = aci_filter.test_filter.id
			}
		}
		provider_to_consumer  ={
			prio = "level3"
			target_dscp = "AF32"
    		relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph2.id
			relation_vz_rs_filt_att {
				action = "permit"
				directives = ["log", "no_stats"]
				priority_override = "default"
				filter_dn = aci_filter.tf_filter.id
			}
		}
	}
```

## Argument Reference

- `contract_dn` - (Required) Distinguished name of parent Contract object.
- `name` - (Required) Name of the contract subject object.
- `annotation` - (Optional) Annotation for the contract subject object.
- `description` - (Optional) Description for the contract subject object.
- `cons_match_t` - (Optional) The subject match criteria across consumers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
- `name_alias` - (Optional) Name alias for the contract subject object.
- `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
- `prov_match_t` - (Optional) The subject match criteria across providers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
- `rev_flt_ports` - (Optional) Enables filter to apply on ingress and egress traffic. Allowed values are "yes" and "no". Default is "yes".
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".

- `relation_vz_rs_subj_graph_att` - (Optional) Relation to class vnsAbsGraph. Cardinality - N_TO_ONE. Type - String.
- `relation_vz_rs_sdwan_pol` - (Optional) Relation to class extdevSDWanSlaPol. Cardinality - N_TO_ONE. Type - String.
- `relation_vz_rs_subj_filt_att` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - Set of String.

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


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract Subject.

## Importing

An existing Contract Subject can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_contract_subject.example <Dn>
```
