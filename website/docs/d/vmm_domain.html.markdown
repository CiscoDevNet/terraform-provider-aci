---
layout: "aci"
page_title: "ACI: aci_vmm_domain"
sidebar_current: "docs-aci-data-source-vmm_domain"
description: |-
  Data source for ACI VMM Domain
---

# aci_vmm_domain #
Data source for ACI VMM Domain

## Example Usage ##

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



## Attribute Reference

* `id` - Attribute id set to the Dn of the VMM Domain.
* `access_mode` - (Optional) access_mode for object vmm_domain.
* `annotation` - (Optional) annotation for object vmm_domain.
* `arp_learning` - (Optional) Enable/Disable arp learning for AVS Domain.
* `ave_time_out` - (Optional) ave_time_out for object vmm_domain.
* `config_infra_pg` - (Optional) Flag to enable config_infra_pg for object vmm_domain.
* `ctrl_knob` - (Optional) Type pf control knob to use.
* `delimiter` - (Optional) delimiter for object vmm_domain.
* `enable_ave` - (Optional) Flag to enable ave for object vmm_domain.
* `enable_tag` - (Optional) Flag enable tagging for object vmm_domain.
* `encap_mode` - (Optional) The layer 2 encapsulation protocol to use with the virtual switch.
* `enf_pref` - (Optional) The switching enforcement preference. This determines whether switches can be done within the virtual switch (Local Switching) or whether all switched traffic must go through the fabric (No Local Switching).
* `ep_inventory_type` - (Optional) Determines which end point inventory_type to use for object vmm_domain. 
* `ep_ret_time` - (Optional) end point retention time for object vmm_domain.
* `hv_avail_monitor` - (Optional) Flag to enable hv_avail_monitor for object vmm_domain.
* `mcast_addr` - (Optional) The multicast address of the VMM domain profile.
* `mode` - (Optional) The switch to be used for the domain profile.
* `name_alias` - (Optional) name_alias for object vmm_domain.
* `pref_encap_mode` - (Optional) The preferred encapsulation mode for object vmm_domain.
