---
subcategory: "Application Management"
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
  application_epg_dn  = aci_application_epg.epg2.id
  tdn                 =  aci_vmm_domain.example.id
}

```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tdn` - (Required) Vmm domain instance.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Domain.
* `annotation` - (Optional) Annotation for object Domain.
* `binding_type` - (Optional) Binding type for object Domain.
* `allow_micro_seg` - (Optional) Boolean flag for allow micro segment. default value will be "false".
"true" maps to class_pref="useg" and "false maps to class_pref="encap"
* `custom_epg_name` - (Optional) Custom EPG name used as name of the VMM port group for the domain.
* `delimiter` - (Optional) Delimiter for object Domain.
* `encap` - (Optional) Port encapsulation.
* `encap_mode` - (Optional) Encap mode for object Domain.
* `epg_cos` - (Optional) Epg cos for object Domain.
* `epg_cos_pref` - (Optional) Epg cos pref for object Domain.
* `instr_imedcy` - (Optional) Determines when policies are pushed to cam.
* `lag_policy_name` - (Optional) Lag policy name for object Domain.
* `netflow_dir` - (Optional) Netflow dir for object Domain.
* `netflow_pref` - (Optional) Netflow pref for object Domain.
* `num_ports` - (Optional) Number of ports existing operationally in module
* `port_allocation` - (Optional) Port allocation for object Domain.
* `primary_encap` - (Optional) Primary encap for object Domain.
* `primary_encap_inner` - (Optional) Primary encap inner for object Domain.
* `res_imedcy` - (Optional) Policy resolution
* `secondary_encap_inner` - (Optional) Secondary encap inner for object Domain.
* `switching_mode` - (Optional) Switching mode for object Domain.
