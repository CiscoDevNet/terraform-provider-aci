---
layout: "aci"
page_title: "ACI: aci_vmm_controller"
sidebar_current: "docs-aci-data-source-vmm_controller"
description: |-
  Data source for ACI VMM Controller
---

# aci_vmm_controller #
Data source for ACI VMM Controller

## Example Usage ##

```hcl
data "aci_vmm_controller" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `name` - (Required) name of Object vmm_controller.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VMM Controller.
* `annotation` - (Optional) annotation for object vmm_controller.
* `dvs_version` - (Optional) 
* `host_or_ip` - (Optional) host or ip
* `inventory_trig_st` - (Optional) 
* `mode` - (Optional) mode of operation
* `msft_config_err_msg` - (Optional) msft_config_err_msg for object vmm_controller.
* `msft_config_issues` - (Optional) 
* `n1kv_stats_mode` - (Optional) n1kv_stats_mode for object vmm_controller.
* `name_alias` - (Optional) name_alias for object vmm_controller.
* `port` - (Optional) service port number for LDAP service
* `root_cont_name` - (Optional) top level container name
* `scope` - (Optional) scope
* `seq_num` - (Optional) isis lsp sequence number
* `stats_mode` - (Optional) statistics mode
* `vxlan_depl_pref` - (Optional) vxlan_depl_pref for object vmm_controller.
