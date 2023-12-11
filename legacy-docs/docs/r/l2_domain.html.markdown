---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_l2_domain"
sidebar_current: "docs-aci-resource-aci_l2_domain"
description: |-
  Manages ACI L2 Domain
---

# aci_l2_domain #
Manages ACI L2 Domain

## Example Usage ##

```hcl
resource "aci_l2_domain" "fool2_domain" {
  name  = "l2_domain_1"
  annotation  = "l2_domain_tag"
  name_alias  = "example"
}
```
## Argument Reference ##
* `name` - (Required) Name of object L2 Domain.
* `annotation` - (Optional) Annotation for object L2 Domain.
* `name_alias` - (Optional) Name alias for object L2 Domain.

* `relation_infra_rs_vlan_ns` - (Optional) Relation to class fvnsVlanInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vip_addr_ns` - (Optional) Relation to class fvnsAddrInst. Cardinality - N_TO_ONE. Type - String.
                
* `relation_extnw_rs_out` - (Optional) Relation to class infraAccGrp. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_dom_vxlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the L2 Domain .

## Importing ##

An existing L2 Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l2_domain.example <Dn>
```
