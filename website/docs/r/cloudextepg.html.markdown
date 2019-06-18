---
layout: "aci"
page_title: "ACI: aci_cloud_external_e_pg"
sidebar_current: "docs-aci-resource-cloud_external_e_pg"
description: |-
  Manages ACI Cloud External EPg
---

# aci_cloud_external_e_pg #
Manages ACI Cloud External EPg
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
resource "aci_cloud_external_e_pg" "example" {

  cloud_applicationcontainer_dn  = "${aci_cloud_applicationcontainer.example.id}"

  name  = "example"
  annotation  = "example"
  exception_tag  = "example"
  flood_on_encap  = "example"
  match_t  = "example"
  name_alias  = "example"
  pref_gr_memb  = "example"
  prio  = "example"
  route_reachability  = "example"
}
```
## Argument Reference ##
* `cloud_applicationcontainer_dn` - (Required) Distinguished name of parent CloudApplicationcontainer object.
* `name` - (Required) name of Object cloud_external_e_pg.
* `annotation` - (Optional) annotation for object cloud_external_e_pg.
* `exception_tag` - (Optional) exception_tag for object cloud_external_e_pg.
* `flood_on_encap` - (Optional) flood_on_encap for object cloud_external_e_pg.
* `match_t` - (Optional) match criteria
* `name_alias` - (Optional) name_alias for object cloud_external_e_pg.
* `pref_gr_memb` - (Optional) pref_gr_memb for object cloud_external_e_pg.
* `prio` - (Optional) qos priority class id
* `route_reachability` - (Optional) route_reachability for object cloud_external_e_pg.

* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_cloud_rs_cloud_e_pg_ctx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud External EPg.

## Importing ##

An existing Cloud External EPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_external_e_pg.example <Dn>
```