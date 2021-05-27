---
layout: "aci"
page_title: "ACI: aci_vmm_domain"
sidebar_current: "docs-aci-data-source-vmm_domain"
description: |-
  Data source for ACI VMM Domain
---

# aci_vmm_domain

Data source for ACI VMM Domain

## Example Usage

```hcl
data "aci_vmm_domain" "dev_vmmdom" {
  provider_profile_dn  = "uni/vmmp-VMware"
  name                 = "demo_vmmdomp"
}
```
## Argument Reference ##
* `provider_profile_dn` - (Required) Distinguished name of parent Provider Profile object.
  * Here is a map of vendor and provider_profile_dn for reference.

        | Vendor         | provider_profile_dn     |
        | -----------    | -----------             |
        | Microsoft      |  uni/vmmp-Microsoft     |
        | CloudFoundry   |  uni/vmmp-CloudFoundry  |
        | OpenShift      |  uni/vmmp-OpenShift     |
        | OpenStack      |  uni/vmmp-OpenStack     |
        | VMware         |  uni/vmmp-VMware        |
        | Kubernetes     |  uni/vmmp-Kubernetes    |
        | Redhat         |  uni/vmmp-Redhat        |

* `name` - (Required) name of Object vmm_domain.

## Argument Reference

- `provider_profile_dn` - (Required) Distinguished name of parent ProviderProfile object.
- `name` - (Required) name of Object vmm_domain.

## Attribute Reference

- `id` - Attribute id set to the Dn of the VMM Domain.
- `access_mode` - (Optional) Access mode for object vmm domain.
- `annotation` - (Optional) Annotation for object vmm domain.
- `arp_learning` - (Optional) Enable/Disable arp learning for AVS Domain.
- `ave_time_out` - (Optional) ACI Virtual Edge time-out for object vmm domain.
- `config_infra_pg` - (Optional) Flag to enable configure infra port groups for object vmm domain.
- `ctrl_knob` - (Optional) Type pf control knob to use. Allowed values are "none" and "epDpVerify".
- `delimiter` - (Optional) Delimiter for object vmm domain.
- `enable_ave` - (Optional) Flag to enable ACI Virtual Edge for object vmm domain.
- `enable_tag` - (Optional) Flag enable tagging for object vmm domain.
- `encap_mode` - (Optional) The layer 2 encapsulation protocol to use with the virtual switch.
- `enf_pref` - (Optional) The switching enforcement preference. This determines whether switches can be done within the virtual switch (Local Switching) or whether all switched traffic must go through the fabric (No Local Switching).
- `ep_inventory_type` - (Optional) Determines which end point inventory type to use for object VMM domain.
- `ep_ret_time` - (Optional) End point retention time for object vmm domain. Allowed value range is "0" - "600". Default value is "0".
- `hv_avail_monitor` - (Optional) Flag to enable host availability monitor for object VMM domain.
- `mcast_addr` - (Optional) The multicast address of the VMM domain profile.
- `mode` - (Optional) The switch to be used for the domain profile.
- `name_alias` - (Optional) Name alias for object VMM domain.
- `pref_encap_mode` - (Optional) The preferred encapsulation mode for object VMM domain.
