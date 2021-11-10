---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_interface_group"
sidebar_current: "docs-aci-resource-l3out_hsrp_interface_group"
description: |-
  Manages ACI L3out HSRP Interface Group
---

# aci_l3out_hsrp_interface_group

Manages ACI L3out HSRP Interface Group

## Example Usage

```hcl
resource "aci_l3out_hsrp_interface_group" "example" {
  l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.example.id
  name                            = "one"
  annotation                      = "example"
  config_issues                   = "GroupMac-Conflicts-Other-Group"
  group_af                        = "ipv4"
  group_id                        = "20"
  group_name                      = "test"
  ip                              = "10.22.30.40"
  ip_obtain_mode                  = "admin"
  mac                             = "02:10:45:00:00:56"
  name_alias                      = "example"
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `name` - (Required) Name of L3out HSRP interface group object.
- `annotation` - (Optional) Annotation for L3out HSRP interface group object.
- `description` - (Optional) Description for L3out HSRP interface group object.
- `config_issues` - (Optional) Configuration issues for L3out HSRP interface group object. Allowed values are "GroupMac-Conflicts-Other-Group", "GroupName-Conflicts-Other-Group", "GroupVIP-Conflicts-Other-Group", "Multiple-Version-On-Interface", "Secondary-vip-conflicts-if-ip", "Secondary-vip-subnet-mismatch", "group-vip-conflicts-if-ip", "group-vip-subnet-mismatch" and "none". Default value is "none".
- `group_af` - (Optional) Group type for L3out HSRP interface group object. Allowed values are "ipv4" and "ipv6". Default value is "ipv4".
- `group_id` - (Optional) Group id for L3out HSRP interface group object.
- `group_name` - (Optional) Group name for L3out HSRP interface group object.
- `ip` - (Optional) IP address for L3out HSRP interface group object.
- `ip_obtain_mode` - (Optional) IP obtain mode for L3out HSRP interface group object. Allowed values are "admin", "auto" and "learn". Default value is "admin".
- `mac` - (Optional) MAC address for L3out HSRP interface group object.
- `name_alias` - (Optional) Name alias for L3out HSRP interface group object.

- `relation_hsrp_rs_group_pol` - (Optional) Relation to class hsrpGroupPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out HSRP interface group.

## Importing

An existing L3out HSRP interface group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_hsrp_interface_group.example <Dn>
```
