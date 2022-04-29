---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-resource-contract_subject"
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
			prio = "level2"
			target_dscp = "AF41"
    		relation_vz_rs_in_term_graph_att = aci_l4_l7_service_graph_template.service_graph2.id
		}
		provider_to_consumer  ={
			prio = "level3"
			target_dscp = "AF32"
    		relation_vz_rs_out_term_graph_att = aci_l4_l7_service_graph_template.service_graph2.id
		}
	}
```

## Argument Reference

- `contract_dn` - (Required) Distinguished name of parent Contract object.
- `name` - (Required) Name of Object contract subject.
- `annotation` - (Optional) Annotation for object contract subject.
- `description` - (Optional) Description for object contract subject.
- `cons_match_t` - (Optional) The subject match criteria across consumers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
- `name_alias` - (Optional) Name alias for object contract subject.
- `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
- `prov_match_t` - (Optional) The subject match criteria across consumers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
- `rev_flt_ports` - (Optional) Enables filter to apply on ingress and egress traffic. Allowed values are "yes" and "no". Default is "yes".
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".

- `relation_vz_rs_subj_graph_att` - (Optional) Relation to class vnsAbsGraph. Cardinality - N_TO_ONE. Type - String.
- `relation_vz_rs_sdwan_pol` - (Optional) Relation to class extdevSDWanSlaPol. Cardinality - N_TO_ONE. Type - String.
- `relation_vz_rs_subj_filt_att` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - Set of String.

- `apply_both_directions` - (Optional) . By default set to "yes".
- ` consumer_to_provider` - (Optional) Filter Chain For Consumer to Provider. Class vzInTerm.
    - `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
    - `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
	- `relation_vz_rs_in_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).
- `provider_to_consumer` - (Optional) Filter Chain For Provider to Consumer. Class vzOutTerm
    - `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server.
    - `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
	- `relation_vz_rs_out_term_graph_att` - (Optional) Relation to class L4-L7 Service Graph (vnsAbsGraph).


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract Subject.

## Importing

An existing Contract Subject can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_contract_subject.example <Dn>
```
