---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_epg"
sidebar_current: "docs-aci-data-source-cloud_external_epg"
description: |-
  Data source for Cloud Network Controller Cloud External EPg
---

# aci_cloud_external_epg #
Data source for Cloud Network Controller Cloud External EPg  
<b>Note: This resource is supported in Cloud Network Controller only.</b>
## Example Usage ##

```hcl
data "aci_cloud_external_epg" "foo_ext_epg" {

  cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.sample_app.id
  name                           = "dev_ext_epg"
}
```
## Argument Reference ##
* `cloud_applicationcontainer_dn` - (Required) Distinguished name of parent Cloud Application container object.
* `name` - (Required) Name of Object Cloud External EPg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud External EPg.
* `description` - (Optional) Description for object Cloud External EPg.
* `annotation` - (Optional) Annotation for object Cloud External EPg.
* `exception_tag` - (Optional) Exception-tag for object Cloud External EPg.
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
* `match_t` - (Optional) The provider label match criteria. 
* `name_alias` - (Optional) Name alias for object Cloud External EPg.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
* `prio` - (Optional) QOS priority class id.
* `route_reachability` - (Optional) Route reachability for this EPG.
