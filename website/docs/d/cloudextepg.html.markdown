---
layout: "aci"
page_title: "ACI: aci_cloud_external_epg"
sidebar_current: "docs-aci-data-source-cloud_external_e_pg"
description: |-
  Data source for ACI Cloud External EPg
---

# aci_cloud_external_epg #
Data source for ACI Cloud External EPg  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
data "aci_cloud_external_epg" "foo_ext_epg" {

  cloud_applicationcontainer_dn  = "${aci_cloud_applicationcontainer.sample_app.id}" 
  name                           = "dev_ext_epg"
}
```
## Argument Reference ##
* `cloud_applicationcontainer_dn` - (Required) Distinguished name of parent CloudApplicationcontainer object.
* `name` - (Required) name of Object cloud_external_e_pg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud External EPg.
* `annotation` - (Optional) annotation for object cloud_external_e_pg.
* `exception_tag` - (Optional) exception_tag for object cloud_external_e_pg.
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
* `match_t` - (Optional) The provider label match criteria. 
* `name_alias` - (Optional) name_alias for object cloud_external_e_pg.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
* `prio` - (Optional) qos priority class id.
* `route_reachability` - (Optional) Route reachability for this EPG.
