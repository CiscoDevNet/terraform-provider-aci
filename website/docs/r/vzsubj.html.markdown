---
layout: "aci"
page_title: "ACI: aci_contract_subject"
sidebar_current: "docs-aci-resource-contract_subject"
description: |-
  Manages ACI Contract Subject
---

# aci_contract_subject #
Manages ACI Contract Subject

## Example Usage ##

```hcl
	resource "aci_contract_subject" "foocontract_subject" {
		contract_dn   = "${aci_contract.example.id}"
		description   = "%s"
		name          = "demo_subject"
		annotation    = "tag_subject"
		cons_match_t  = "AtleastOne"
		name_alias    = "alias_subject"
		prio          = "level1"
		prov_match_t  = "AtleastOne"
		rev_flt_ports = "yes"
		target_dscp   = "CS0"
	}
```
## Argument Reference ##
* `contract_dn` - (Required) Distinguished name of parent Contract object.
* `name` - (Required) name of Object contract_subject.
* `annotation` - (Optional) annotation for object contract_subject.
* `cons_match_t` - (Optional) The subject match criteria across consumers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
* `name_alias` - (Optional) name_alias for object contract_subject.
* `prio` - (Optional) The priority level of a sub application running behind an endpoint group, such as an Exchange server. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
* `prov_match_t` - (Optional) The subject match criteria across consumers. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
* `rev_flt_ports` - (Optional) enables filter to apply on ingress and egress traffic. Allowed values are "yes" and "no". Default is "yes".
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".

* `relation_vz_rs_subj_graph_att` - (Optional) Relation to class vnsAbsGraph. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vz_rs_sdwan_pol` - (Optional) Relation to class extdevSDWanSlaPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vz_rs_subj_filt_att` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract Subject.

## Importing ##

An existing Contract Subject can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract_subject.example <Dn>
```