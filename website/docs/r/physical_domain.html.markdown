---
layout: "aci"
page_title: "ACI: aci_physical_domain"
sidebar_current: "docs-aci-resource-physical_domain"
description: |-
  Manages ACI Physical Domain
---

# aci_physical_domain #
Manages ACI Physical Domain

## Example Usage ##

```hcl
resource "aci_physical_domain" "example" {


  name  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object physical_domain.
* `annotation` - (Optional) annotation for object physical_domain.
* `name_alias` - (Optional) name_alias for object physical_domain.

* `relation_infra_rs_vlan_ns` - (Optional) Relation to class fvnsVlanInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vip_addr_ns` - (Optional) Relation to class fvnsAddrInst. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_dom_vxlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Physical Domain.

## Importing ##

An existing Physical Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_physical_domain.example <Dn>
```