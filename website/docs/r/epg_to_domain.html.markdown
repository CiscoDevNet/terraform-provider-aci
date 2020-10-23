---
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

  application_epg_dn    = "${aci_application_epg.example.id}"
  tdn                   = "${aci_vmm_domain.example.id}"
  vmm_allow_promiscuous = "accept"
  vmm_forged_transmits  = "reject"
  vmm_mac_changes       = "accept"
}

```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tdn` - (Required) Distinguished Name of Target Domain object.
* `annotation` - (Optional) annotation for object domain.
* `binding_type` - (Optional) binding_type for object domain.
Allowed values: "none", "staticBinding", "dynamicBinding", "ephemeral",
* `allow_micro_seg` - (Optional) boolean flag for allow micro segment. default value will be "false".
"true" maps to class_pref="useg" and "false maps to class_pref="encap" 
* `delimiter` - (Optional) delimiter for object domain.
* `encap` - (Optional) port encapsulation
* `encap_mode` - (Optional) encap_mode for object domain.
Allowed values: "auto", "vlan", "vxlan"
* `epg_cos` - (Optional) epg_cos for object domain.
Allowed values: "Cos0", "Cos1", "Cos2", "Cos3", "Cos4", "Cos5", "Cos6", "Cos7"
* `epg_cos_pref` - (Optional) epg_cos_pref for object domain.
Allowed values: "disabled", "enabled"
* `instr_imedcy` - (Optional) determines when policies are pushed to cam.
Allowed values: "immediate", "lazy"
* `lag_policy_name` - (Optional) lag_policy_name for object domain.
* `netflow_dir` - (Optional) netflow_dir for object domain.
Allowed values: "ingress", "egress", "both"
* `netflow_pref` - (Optional) netflow_pref for object domain.
Allowed values: "disabled", "enabled"
* `num_ports` - (Optional) number of ports existing operationally in module
* `port_allocation` - (Optional) port_allocation for object domain.
* `primary_encap` - (Optional) primary_encap for object domain.
* `primary_encap_inner` - (Optional) primary_encap_inner for object domain.
* `res_imedcy` - (Optional) policy resolution.
Allowed values: "immediate", "lazy", "pre-provision"
* `secondary_encap_inner` - (Optional) secondary_encap_inner for object domain.
* `switching_mode` - (Optional) switching_mode for object domain.
Allowed values: "native", "AVE"


* `vmm_allow_promiscuous` - (Optional) allow_promiscuous for object vmm_security_policy.
* `vmm_forged_transmits` - (Optional) forged_transmits for object vmm_security_policy.
* `vmm_mac_changes` - (Optional) mac_changes for object vmm_security_policy.



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