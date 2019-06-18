---
layout: "aci"
page_title: "ACI: aci_filter"
sidebar_current: "docs-aci-resource-filter"
description: |-
  Manages ACI Filter
---

# aci_filter #
Manages ACI Filter

## Example Usage ##

```hcl
	resource "aci_filter" "foofilter" {
		tenant_dn   = "${aci_tenant.example.id}"
		description = "%s"
		name        = "demo_filter"
		annotation  = "tag_filter"
		name_alias  = "alias_filter"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object filter.
* `annotation` - (Optional) annotation for object filter.
* `name_alias` - (Optional) name_alias for object filter.

* `relation_vz_rs_filt_graph_att` - (Optional) Relation to class vnsInTerm. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vz_rs_fwd_r_flt_p_att` - (Optional) Relation to class vzAFilterableUnit. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vz_rs_rev_r_flt_p_att` - (Optional) Relation to class vzAFilterableUnit. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Filter.

## Importing ##

An existing Filter can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_filter.example <Dn>
```