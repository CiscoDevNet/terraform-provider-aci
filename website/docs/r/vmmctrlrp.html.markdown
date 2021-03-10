---
layout: "aci"
page_title: "ACI: aci_vmm_controller"
sidebar_current: "docs-aci-resource-vmm_controller"
description: |-
  Manages ACI VMM Controller
---

# aci_vmm_controller #
Manages ACI VMM Controller

## Example Usage ##

```hcl
resource "aci_vmm_controller" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"
  name  = "example"
  annotation  = "example"
  dvs_version  = "example"
  host_or_ip  = "example"
  inventory_trig_st  = "example"
  mode  = "example"
  msft_config_err_msg  = "example"
  msft_config_issues  = "example"
  n1kv_stats_mode  = "example"
  name_alias  = "example"
  port  = "example"
  root_cont_name  = "example"
  scope  = "example"
  seq_num  = "example"
  stats_mode  = "example"
  vxlan_depl_pref  = "example"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `name` - (Required) name of Object vmm_controller.
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
