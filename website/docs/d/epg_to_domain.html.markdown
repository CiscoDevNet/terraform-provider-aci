---
layout: "aci"
page_title: "ACI: aci_epg_to_domain"
sidebar_current: "docs-aci-data-source-epg_to_domain"
description: |-
  Data source for ACI epg to Domain
---

# aci_domain #
Data source for ACI epg to Domain

## Example Usage ##

```hcl

data "aci_epg_to_domain" "temp" {
  application_epg_dn  = "${aci_application_epg.epg2.id}"
  tdn                = "${aci_vmm_domain.example.id}"
}

```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tdn` - (Required) Vmm domain instance.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Domain.
* `annotation` - (Optional) annotation for object domain.
* `binding_type` - (Optional) binding_type for object domain.
* `allow_micro_seg` - (Optional) boolean flag for allow micro segment. default value will be "false".
"true" maps to class_pref="useg" and "false maps to class_pref="encap"
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
