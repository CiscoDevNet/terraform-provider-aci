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
data "aci_vlan_encapsulationfor_vxlan_traffic" "dev_vlan_traffic" {
  attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.example.id
}
```
## Argument Reference ##
* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Vlan Encapsulation for Vxlan Traffic.
* `annotation` - (Optional) Annotation for object vlan encapsulation for vxlan traffic.
* `name_alias` - (Optional) Name alias for object vlan encapsulation for vxlan traffic.
* `description`- (Optional) Description for object vlan encapsulation for vxlan traffic.
