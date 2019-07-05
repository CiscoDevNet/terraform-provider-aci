---
layout: "aci"
page_title: "ACI: aci_subnet"
sidebar_current: "docs-aci-resource-subnet"
description: |-
  Manages ACI Subnet
---

# aci_subnet #
Manages ACI Subnet

## Example Usage ##

```hcl
resource "aci_l3_ext_subnet" "example" {

  external_network_instance_profile_dn  = "${aci_external_network_instance_profile.example.id}"

  ip  = "example"
  aggregate  = "example"
  annotation  = "example"
  name_alias  = "example"
  scope  = "example"
}
```
## Argument Reference ##
* `external_network_instance_profile_dn` - (Required) Distinguished name of parent ExternalNetworkInstanceProfile object.
* `ip` - (Required) ip of Object subnet.
* `aggregate` - (Optional) aggregate for object subnet.
* `annotation` - (Optional) annotation for object subnet.
* `ip` - (Optional) ip address
* `name_alias` - (Optional) name_alias for object subnet.
* `scope` - (Optional) capability domain

* `relation_l3ext_rs_subnet_to_profile` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_M. Type - Set of Map.
                
* `relation_l3ext_rs_subnet_to_rt_summ` - (Optional) Relation to class rtsumARtSummPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Subnet.

## Importing ##

An existing Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_subnet.example <Dn>
```