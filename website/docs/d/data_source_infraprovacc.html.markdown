---
layout: "aci"
page_title: "ACI: aci_vlan_encapsulationfor_vxlan_traffic"
sidebar_current: "docs-aci-data-source-vlan_encapsulationfor_vxlan_traffic"
description: |-
  Data source for ACI Vlan Encapsulation for Vxlan Traffic
---

# aci_vlan_encapsulationfor_vxlan_traffic #
Data source for ACI Vlan Encapsulation for Vxlan Traffic

## Example Usage ##

```hcl
data "aci_vlan_encapsulationfor_vxlan_traffic" "example" {
  attachable_access_entity_profile_dn  = "${aci_attachable_access_entity_profile.example.id}"
}
```
## Argument Reference ##
* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Vlan Encapsulation for Vxlan Traffic.
* `annotation` - (Optional) annotation for object vlan_encapsulationfor_vxlan_traffic.
* `name_alias` - (Optional) name_alias for object vlan_encapsulationfor_vxlan_traffic.
