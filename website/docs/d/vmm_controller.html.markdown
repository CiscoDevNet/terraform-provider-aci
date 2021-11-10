---
subcategory: "Virtual Networking"
layout: "aci"
page_title: "ACI: aci_vmm_controller"
sidebar_current: "docs-aci-data-source-vmm_controller"
description: |-
  Data source for ACI VMM Controller
---

# aci_vmm_controller #

Data source for ACI VMM Controller

## API Information ##

* `Class` - vmmCtrlrP
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/ctrlr-{name}

## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> {controller_name}

## Example Usage ##

```hcl
data "aci_vmm_controller" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  name  = "example"
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
* `name` - (Required) name of object VMM Controller.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the VMM Controller.
* `annotation` - (Optional) Annotation of object VMM Controller.
* `name_alias` - (Optional) Name Alias of object VMM Controller.
* `dvs_version` - (Optional) Dvs Version.
* `host_or_ip` - (Optional) Hostname or IP Address.
* `inventory_trig_st` - (Optional) Triggered Inventory Sync Status.
* `mode` - (Optional) The mode of operation.
* `msft_config_err_msg` - (Optional) Deployment Error Message of Mirosoft Plugin SCVM Controller.
                    It captures error message encountered in SCVMM Controller
                    plugin. This error message represents specific details for bitmask
                    based msftConfigIssues fault.
* `msft_config_issues` - (Optional) msftConfigIssues.
* `n1kv_stats_mode` - (Optional) n1kv statistics enable.
* `port` - (Optional) Port. Port
* `root_cont_name` - Top level container name.
* `scope` - (Optional) The VMM control policy scope.
* `seq_num` - (Optional) An ISIS link-state packet sequence number.
* `stats_mode` - (Optional) The statistics mode.
* `vxlan_depl_pref` - (Optional) VxLAN Deployment Preference.
