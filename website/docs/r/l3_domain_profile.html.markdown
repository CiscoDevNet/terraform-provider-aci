---
layout: "aci"
page_title: "ACI: aci_l3_domain_profile"
sidebar_current: "docs-aci-resource-l3_domain_profile"
description: |-
  Manages ACI L3 Domain Profile
---

# aci_l3_domain_profile #
Manages ACI L3 Domain Profile

## Example Usage ##

```hcl
resource "aci_l3_domain_profile" "example" {
  name  = "l3_domain_profile_1"
  annotation  = "l3_domain_profile_tag"
  name_alias  = "alias_name"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object l3 domain profile.
* `annotation` - (Optional) Annotation for object l3 domain profile.
* `name_alias` - (Optional) Name alias for object l3 domain profile.

* `relation_infra_rs_vlan_ns` - (Optional) Relation to class fvnsVlanInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vip_addr_ns` - (Optional) Relation to class fvnsAddrInst. Cardinality - N_TO_ONE. Type - String.
                
* `relation_extnw_rs_out` - (Optional) Relation to class infraAccGrp. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_dom_vxlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3 Domain Profile.

## Importing ##

An existing L3 Domain Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l3_domain_profile.example <Dn>
```