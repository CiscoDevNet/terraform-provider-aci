---
layout: "aci"
page_title: "ACI: aci_vmm_controller"
sidebar_current: "docs-aci-resource-vmm_controller"
description: |-
  Manages ACI VMM Controller
---

# aci_vmm_controller #

Manages ACI VMM Controller

## API Information ##

* `Class` - vmmCtrlrP
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/ctrlr-{name}

## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> {controller_name}

## Example Usage ##

```hcl
resource "aci_vmm_controller" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  dvs_version = "unmanaged"
  host_or_ip = "10.10.10.10"
  inventory_trig_st = "untriggered"
  mode = "default"
  msft_config_err_msg = "Error"
  msft_config_issues = ["not-applicable"]
  n1kv_stats_mode = "enabled"
  port = "0"
  root_cont_name = "vmmdc"
  scope = "vm"
  seq_num = "0"
  stats_mode = "disabled"
  vxlan_depl_pref = "vxlan"

  vmm_rs_acc = aci_resource.example.id

  vmm_rs_ctrlr_p_mon_pol = aci_resource.example.id

  vmm_rs_mcast_addr_ns = aci_resource.example.id

  vmm_rs_mgmt_e_pg = aci_resource.example.id

  vmm_rs_to_ext_dev_mgr = [aci_resource.example.id]

  vmm_rs_vmm_ctrlr_p {
    epg_depl_pref = "local"
    target_dn = aci_resource.example.id
  }

  vmm_rs_vxlan_ns = aci_resource.example.id

  vmm_rs_vxlan_ns_def = aci_resource.example.id
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
* `name` - (Required) Name of object VMM Controller.
* `annotation` - (Optional) Annotation of object VMM Controller.
* `dvs_version` - (Optional) Dvs Version.  Allowed values are "5.1", "5.5", "6.0", "6.5", "6.6", "7.0", "unmanaged", and default value is "unmanaged". Type: String.
* `host_or_ip` - (Optional) Hostname or IP Address.
* `inventory_trig_st` - (Optional) Triggered Inventory Sync Status.  Allowed values are "autoTriggered", "triggered", "untriggered", and default value is "untriggered". Type: String.
* `mode` - (Optional) The mode of operation. Allowed values are "cf", "default", "k8s", "n1kv", "nsx", "openshift", "ovs", "rancher", "rhev", "unknown", and default value is "default". Type: String.
* `msft_config_err_msg` - (Optional) Deployment Error Message of Mirosoft Plugin SCVM Controller.
                    It captures error message encountered in SCVMM Controller
                    plugin. This error message represents specific details for bitmask
                    based msftConfigIssues fault.
* `msft_config_issues` - (Optional) msftConfigIssues. Allowed values are "aaacert-invalid", "duplicate-mac-in-inventory", "duplicate-rootContName", "invalid-object-in-inventory", "invalid-rootContName", "inventory-failed", "missing-hostGroup-in-cloud", "missing-rootContName", "not-applicable", "zero-mac-in-inventory", and default value is "not-applicable". Type: List.
* `n1kv_stats_mode` - (Optional) n1kv statistics enable. Allowed values are "disabled", "enabled", "unknown", and default value is "enabled". Type: String.
* `port` - (Optional) Default value is "0".
* `root_cont_name` - (Optional) Top level container name. Type: String.
* `scope` - (Optional) The VMM control policy scope. Allowed values are "MicrosoftSCVMM", "cloudfoundry", "iaas", "kubernetes", "network", "nsx", "openshift", "openstack", "rhev", "unmanaged", "vm", and default value is "vm". Type: String.
* `seq_num` - (Optional) An ISIS link-state packet sequence number. Default value is "0".
* `stats_mode` - (Optional) The statistics mode. Allowed values are "disabled", "enabled", "unknown", and default value is "disabled". Type: String.
* `vxlan_depl_pref` - (Optional) VxLAN Deployment Preference. Allowed values are "nsx", "vxlan", and default value is "vxlan". Type: String.

* `relation_vmm_rs_acc` - (Optional) Represents the relation to a User Access Profile (class vmmUsrAccP). A source relation to the user account profile. Type: String.

* `relation_vmm_rs_ctrlr_p_mon_pol` - (Optional) Represents the relation to a Monitoring Policy (class monInfraPol). A source relation to the monitoring policy model for the infra semantic scope. Type: String.

* `relation_vmm_rs_mcast_addr_ns` - (Optional) Represents the relation to a Multicast Addr Pool (class fvnsMcastAddrInstP). A source relation to the policy definition of the multicast IP address ranges. Type: String.

* `relation_vmm_rs_mgmt_e_pg` - (Optional) Represents the relation to a Management EPg (class fvEPg). A source relation to a set of endpoints. Type: String.

* `relation_vmm_rs_to_ext_dev_mgr` - (Optional) Represents the relation to an External device mgr profile (class extdevMgrP). Association to External Device Controller Profile Type: List.

* `relation_vmm_rs_vmm_ctrlr_p` - (Optional) A block representing the relation to a Vmm Controller Profile (class vmmCtrlrP). It's a source relation to the VMM controller profile. The VMM controller profile is a policy pertaining to a single VM management domain that also corresponds to a single policy enforcement domain. A cluster of VMware VCs forms such a domain. Type: Block.
  * `epg_depl_pref` - (Optional) Allowed values are "both", "local", and default value is "local". Type: String.
  * `target_dn` - (Required) The distinguished name of the target controller. Type: String.

* `relation_vmm_rs_vxlan_ns` - (Optional) Represents the relation to a VXLAN Pool (class fvnsVxlanInstP). It's a source relation to the VxLAN namespace policy definition. Type: String.

* `relation_vmm_rs_vxlan_ns_def` - (Optional) Represents the relation to an Unresolvable Relation to VxLAN Pool (class fvnsAInstP). A source relation to the namespace policy is used for managing the encap (VXLAN, NVGRE, VLAN) ranges. Type: String.

## Importing ##

An existing VMMController can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import vmm_controller.example <Dn>
```