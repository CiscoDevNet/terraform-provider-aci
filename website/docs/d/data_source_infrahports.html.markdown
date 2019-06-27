---
layout: "aci"
page_title: "ACI: aci_access_port_selector"
sidebar_current: "docs-aci-data-source-access_port_selector"
description: |-
  Data source for ACI Access Port Selector
---

# aci_access_port_selector #
Data source for ACI Access Port Selector

## Example Usage ##

```hcl
data "aci_access_port_selector" "dev_acc_port_select" {
  leaf_interface_profile_dn  = "${aci_leaf_interface_profile.example.id}"
  name                       = "foo_acc_port_select"
  access_port_selector_type  = "ALL"
}
```
## Argument Reference ##
* `leaf_interface_profile_dn` - (Required) Distinguished name of parent LeafInterfaceProfile object.
* `name` - (Required) name of Object access_port_selector.
* `access_port_selector_type` - (Required) access_port_selector_type of Object access_port_selector.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Port Selector.
* `annotation` - (Optional) annotation for object access_port_selector.
* `name_alias` - (Optional) name_alias for object access_port_selector.
* `access_port_selector_type` - (Optional) host port selector type.
