---
layout: "aci"
page_title: "ACI: aci_fabric_wide_settings"
sidebar_current: "docs-aci-resource-fabric_wide_settings"
description: |-
  Manages ACI Fabric-Wide Settings Policy
---

# aci_fabric_wide_settings #
Manages ACI Fabric-Wide Settings Policy

## API Information ##
* `Class` - infraSetPol
* `Distinguished Named` - uni/infra/settings

## GUI Information ##
* `Location` - System -> System Settings -> Fabric-Wide Settings 

## Example Usage ##

```hcl
resource "aci_fabric_wide_settings" "example" {
  name = "example"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
  disable_ep_dampening = "yes"
  enable_mo_streaming = "yes"
  enable_remote_leaf_direct = "yes"
  enforce_subnet_check = "yes"
  opflexp_authenticate_clients = "yes"
  opflexp_use_ssl = "yes"
  restrict_infra_vlan_traffic = "yes"
  unicast_xr_ep_learn_disable = "yes"
  validate_overlapping_vlans = "yes"
}
```

## NOTE ##
Users can use the resource of type `aci_fabric_wide_settings` to change the configuration of the object Fabric Wide Settings. Users cannot create more than one instance of object Fabric Wide Settings.

## Argument Reference ##
* `name` - (Optional) Name of object Fabric-Wide Settings Policy.
* `annotation` - (Optional) Annotation of object Fabric-Wide Settings Policy.
* `description` - (Optional) Description of object Fabric-Wide Settings Policy.
* `name_alias` - (Optional) Name alias of object Fabric-Wide Settings Policy.
* `disable_ep_dampening` - (Optional) Disable Ep Dampening knob of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `enable_mo_streaming` - (Optional) Enable MO streaming of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `enable_remote_leaf_direct` - (Optional) Enable remote leaf direct communication of object Fabric-Wide Settings Policy.  Allowed values are "yes" and "no". 
* `enforce_subnet_check` - (Optional) Enforce subnet check of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `opflexp_authenticate_clients` - (Optional) Opflexp Client Certificates for authentication of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `opflexp_use_ssl` - (Optional) SSL transport for Opflexp indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `restrict_infra_vlan_traffic` - (Optional) Intra Leaf Communication traffic indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `unicast_xr_ep_learn_disable` - (Optional) Disable xrLeanrs indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `validate_overlapping_vlans` - (Optional) Validate Overlapping VLANS indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".


## Importing ##

An existing Fabric-WideSettings Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_wide_settings.example <Dn>
```