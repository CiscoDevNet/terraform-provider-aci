---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_domain"
sidebar_current: "docs-aci-resource-epg_to_domain"
description: |-
  Manages ACI epg to Domain
---

# aci_epg_to_domain #

Manages ACI epg to Domain

## Example Usage ##

```hcl
resource "aci_epg_to_domain" "example" {
  application_epg_dn    = aci_application_epg.example.id
  tdn                   = aci_fc_domain.foofc_domain.id
  annotation            = "annotation"
  binding_type          = "none"
  allow_micro_seg       = "false"
  delimiter             = ""
  encap                 = "vlan-5"
  encap_mode            = "auto"
  epg_cos               = "Cos0"
  epg_cos_pref          = "disabled"
  instr_imedcy          = "lazy"
  enhanced_lag_policy   = "uni/vmmp-VMware/dom-aci_terraform_lab/vswitchpolcont/enlacplagp-lab_lacp"
  netflow_dir           = "both"
  netflow_pref          = "disabled"
  num_ports             = "0"
  port_allocation       = "none"
  primary_encap         = "unknown"
  primary_encap_inner   = "unknown"
  res_imedcy            = "lazy"
  secondary_encap_inner = "unknown"
  switching_mode        = "native"
  vmm_allow_promiscuous = "accept"
  vmm_forged_transmits  = "reject"
  vmm_mac_changes       = "accept"
  custom_epg_name       = "epg_lab"
}
```
## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tdn` - (Required) Distinguished Name of Target Domain object.
* `annotation` - (Optional) Annotation for object Domain.
* `binding_type` - (Optional) Binding type for object Domain.
Allowed values: "none", "staticBinding", "dynamicBinding", "ephemeral". Default value: "none"
* `allow_micro_seg` - (Optional) Boolean flag for allow micro segment. default value will be "false".
"true" maps to class_pref="useg" and "false maps to class_pref="encap" 
* `custom_epg_name` - (Optional) Custom EPG name used as name of the VMM port group for the domain.
* `enhanced_lag_policy` - (Optional) Distinguished Name of the Enhanced LACP LAG Policy (class lacpEnhancedLagPol) to associate with the VMM domain.
* `delimiter` - (Optional) Delimiter for object Domain.
* `encap` - (Optional) Port encapsulation.
* `encap_mode` - (Optional) Encap mode for object Domain.
Allowed values: "auto", "vlan", "vxlan". Default value: "auto"
* `epg_cos` - (Optional) Epg cos for object Domain.
Allowed values: "Cos0", "Cos1", "Cos2", "Cos3", "Cos4", "Cos5", "Cos6", "Cos7". Default value: "Cos0"
* `epg_cos_pref` - (Optional) Epg cos pref for object Domain.
Allowed values: "disabled", "enabled". Default value: "disabled"
* `instr_imedcy` - (Optional) Determines when policies are pushed to cam.
Allowed values: "immediate", "lazy". Default value: "lazy"
* `lag_policy_name` - (Optional) **Deprecated** Lag policy name for object Domain. Use `enhanced_lag_policy` instead.
* `netflow_dir` - (Optional) Netflow dir for object Domain.
Allowed values: "ingress", "egress", "both". Default value: "both"
* `netflow_pref` - (Optional) Netflow pref for object Domain.
Allowed values: "disabled", "enabled"
* `num_ports` - (Optional) Number of ports existing operationally in module. Default value: "0"
* `port_allocation` - (Optional) Port allocation for object Domain.
Allowed values: "none", "elastic", "fixed". Default value: "none"
* `primary_encap` - (Optional) Primary encap for object Domain. Default value: "unknown".
* `primary_encap_inner` - (Optional) Primary encap inner for object Domain. Default value: "unknown".
* `res_imedcy` - (Optional) Policy resolution.
Allowed values: "immediate", "lazy", "pre-provision". Default value: "lazy"
* `secondary_encap_inner` - (Optional) Secondary encap inner for object Domain.Default value: "unknown".
* `switching_mode` - (Optional) Switching mode for object domain.
Allowed values: "native", "AVE". Default value: "native"

* `vmm_allow_promiscuous` - (Optional) Allow promiscuous for object Vmm security policy.
* `vmm_forged_transmits` - (Optional) Forged transmits for object Vmm security policy.
* `vmm_mac_changes` - (Optional) Mac changes for object Vmm security policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Domain.

* `vmm_id` - which is set to the Dn of the VMM Security Policy.

## Importing ##

An existing Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_domain.example <Dn>
```
