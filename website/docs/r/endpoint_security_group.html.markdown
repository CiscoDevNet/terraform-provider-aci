---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-resource-endpoint_security_group"
description: |-
  Manages ACI Endpoint Security Group
---

# aci_endpoint_security_group #
Manages ACI Endpoint Security Group

## Example Usage ##

```hcl
resource "aci_endpoint_security_group" "example" {

  application_profile_dn  = "${aci_application_profile.example.id}"

  name  = "example"
  annotation  = "example"
  exception_tag  = "example"
  flood_on_encap  = "example"
  match_t  = "example"
  name_alias  = "example"
  pc_enf_pref  = "example"
  pref_gr_memb  = "example"
  prio  = "example"
  userdom  = "example"
}
```
## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) name of Object endpoint_security_group.
* `annotation` - (Optional) annotation for object endpoint_security_group.
* `exception_tag` - (Optional) exception_tag for object endpoint_security_group.
* `flood_on_encap` - (Optional) flood_on_encap for object endpoint_security_group.
* `match_t` - (Optional) match criteria
* `name_alias` - (Optional) name_alias for object endpoint_security_group.
* `pc_enf_pref` - (Optional) enforcement preference
* `pref_gr_memb` - (Optional) pref_gr_memb for object endpoint_security_group.
* `prio` - (Optional) qos priority class id
* `userdom` - (Optional) userdom for object endpoint_security_group.

* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_scope` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Endpoint Security Group.

## Importing ##

An existing Endpoint Security Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_security_group.example <Dn>
```