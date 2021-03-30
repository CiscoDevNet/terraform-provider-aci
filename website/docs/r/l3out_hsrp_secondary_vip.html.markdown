---
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_secondary_vip"
sidebar_current: "docs-aci-resource-l3out_hsrp_secondary_vip"
description: |-
  Manages ACI L3out HSRP Secondary VIP
---

# aci_l3out_hsrp_secondary_vip

Manages ACI L3out HSRP Secondary VIP

## Example Usage

```hcl
resource "aci_l3out_hsrp_secondary_vip" "example" {

  l3out_hsrp_interface_group_dn  = "${aci_l3out_hsrp_interface_group.example.id}"
  ip  = "10.0.0.1"
  annotation  = "example"
  config_issues = "GroupMac-Conflicts-Other-Group"
  name_alias  = "example"

}
```

## Argument Reference

- `l3out_hsrp_interface_group_dn` - (Required) Distinguished name of parent hsrp group profile object.
- `ip` - (Required) IP of Object L3out HSRP Secondary VIP.
- `annotation` - (Optional) Annotation for object L3out HSRP Secondary VIP.
- `description` - (Optional) Description for object L3out HSRP Secondary VIP.
- `config_issues` - (Optional) Configuration Issues.  
  Allowed values: "GroupMac-Conflicts-Other-Group", "GroupName-Conflicts-Other-Group", "GroupVIP-Conflicts-Other-Group", "Multiple-Version-On-Interface", "Secondary-vip-conflicts-if-ip", "Secondary-vip-subnet-mismatch", "group-vip-conflicts-if-ip", "group-vip-subnet-mismatch", "none". Default value: "none".
- `ip` - (Optional) IP address.
- `name_alias` - (Optional) Name alias for object L3out HSRP Secondary VIP.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out HSRP Secondary VIP.

## Importing

An existing L3out HSRP Secondary VIP can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_hsrp_secondary_vip.example <Dn>
```
