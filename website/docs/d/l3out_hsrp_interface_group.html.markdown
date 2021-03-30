---
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_interface_group"
sidebar_current: "docs-aci-data-source-l3out_hsrp_interface_group"
description: |-
  Data source for ACI HSRP Interface Group
---

# aci_l3out_hsrp_interface_group #
Data source for ACI HSRP Interface Group

## Example Usage ##

```hcl
data "aci_l3out_hsrp_interface_group" "check" {
  l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.example.id
  name                            = "one"
}
```


## Argument Reference ##

* `l3out_hsrp_interface_profile_dn` - (Required) Distinguished name of parent L3out HSRP interface profile object.
* `name` - (Required) Name of L3out HSRP interface group object.



## Attribute Reference

* `id` - Attribute id set to the Dn of L3out HSRP interface group object.
* `annotation` - Annotation for L3out HSRP interface group object.
* `description` - Description for L3out HSRP interface group object.
* `config_issues` - Configuration issues for L3out HSRP interface group object.
* `group_af` - Group type for L3out HSRP interface group object.
* `group_id` - Group id for L3out HSRP interface group object.
* `group_name` - Group Name for L3out HSRP interface group object.
* `ip` - IP address for L3out HSRP interface group object.
* `ip_obtain_mode` - IP obtain mode for L3out HSRP interface group object.
* `mac` - MAC address for L3out HSRP interface group object.
* `name_alias` - Name alias for L3out HSRP interface group object.
