---
layout: "aci"
page_title: "ACI: aci_vlan_encapsulationfor_vxlan_traffic"
sidebar_current: "docs-aci-resource-vlan_encapsulationfor_vxlan_traffic"
description: |-
  Manages ACI Vlan Encapsulation for Vxlan Traffic
---

# aci_vlan_encapsulationfor_vxlan_traffic #
Manages ACI Vlan Encapsulation for Vxlan Traffic

## Example Usage ##

```hcl
resource "aci_vlan_encapsulationfor_vxlan_traffic" "example" {
  attachable_access_entity_profile_dn  = "${aci_attachable_access_entity_profile.example.id}"
  annotation                           = "tag_traffic"
  name_alias                           = "alias_traffic"
}
```
## Argument Reference ##
* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.
* `annotation` - (Optional) Annotation for object vlan encapsulation for vxlan traffic.
* `name_alias` - (Optional) Name alias for object vlan encapsulation for vxlan traffic.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Vlan Encapsulation for Vxlan Traffic.

## Importing ##

An existing Vlan Encapsulation for Vxlan Traffic can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vlan_encapsulationfor_vxlan_traffic.example <Dn>
```