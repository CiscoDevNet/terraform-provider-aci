---
layout: "aci"
page_title: "ACI: aci_contract"
sidebar_current: "docs-aci-resource-contract"
description: |-
  Manages ACI Contract
---

# aci_contract #
Manages ACI Contract

## Example Usage ##

```hcl
	resource "aci_contract" "foocontract" {
		tenant_dn   = "${aci_tenant.dev_tenant.id}"
		description = "%s"
		name        = "demo_contract"
		annotation  = "tag_contract"
		name_alias  = "alias_contract"
		prio        = "level1"
		scope       = "tenant"
		target_dscp = "unspecified"
		filter {
  			filter_name = "example1"
  			description = "first filter from contract resource"
  			annotation = "tag_filter"
  			name_alias = "example1"
		}
		filter {
  			filter_name = "example2"
  			description = "second filter from contract resource"
  			annotation = "tag_filter"
  			name_alias = "example2"
		}

	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object contract.
* `annotation` - (Optional) annotation for object contract.
* `name_alias` - (Optional) name_alias for object contract.
* `prio` - (Optional) priority level of the service contract.  Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified".
* `scope` - (Optional)  Represents the scope of this contract. If the scope is set as application-profile, the epg can only communicate with epgs in the same application-profile. Allowed values are "global", "tenant", "application-profile" and "context". Default is "context".

* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".	

* `relation_vz_rs_graph_att` - (Optional) Relation to class vnsAbsGraph. Cardinality - N_TO_ONE. Type - String.
                

* `filter` - (Optional) to manage filters from the contract resource. It has the attributes like filter_name, annotation, description and name_alias.
* `filter.filter_name` - (Required) Name of the filter object.
* `filter.description` - (Optional) Description for the filter object.
* `filter.annotation` - (Optional) Annotation for filter object.
* `filter.name_alias` - (Optional) Name alias for filter object.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract.
* `filter.id` - exports this attribute for filter object. Set to the Dn for the filter managed by the contract.

## Importing ##

An existing Contract can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract.example <Dn>
```