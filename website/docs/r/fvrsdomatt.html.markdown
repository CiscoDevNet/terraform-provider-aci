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
* `class_pref` - (Optional) class_pref for object domain.
* `delimiter` - (Optional) delimiter for object domain.
* `encap` - (Optional) port encapsulation
* `encap_mode` - (Optional) encap_mode for object domain.
* `epg_cos` - (Optional) epg_cos for object domain.
* `epg_cos_pref` - (Optional) epg_cos_pref for object domain.
* `instr_imedcy` - (Optional) determines when policies are pushed to cam
* `lag_policy_name` - (Optional) lag_policy_name for object domain.
* `netflow_dir` - (Optional) netflow_dir for object domain.
* `netflow_pref` - (Optional) netflow_pref for object domain.
* `num_ports` - (Optional) number of ports existing operationally in module
* `port_allocation` - (Optional) port_allocation for object domain.
* `primary_encap` - (Optional) primary_encap for object domain.
* `primary_encap_inner` - (Optional) primary_encap_inner for object domain.
* `res_imedcy` - (Optional) policy resolution
* `secondary_encap_inner` - (Optional) secondary_encap_inner for object domain.
* `switching_mode` - (Optional) switching_mode for object domain.


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