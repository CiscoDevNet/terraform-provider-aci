---
layout: "aci"
page_title: "ACI: aci_cloud_e_pg"
sidebar_current: "docs-aci-data-source-cloud_e_pg"
description: |-
  Data source for ACI Cloud EPg
---

# aci_cloud_e_pg #
Data source for ACI Cloud EPg
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_e_pg" "example" {

  cloud_applicationcontainer_dn  = "${aci_cloud_applicationcontainer.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `cloud_applicationcontainer_dn` - (Required) Distinguished name of parent CloudApplicationcontainer object.
* `name` - (Required) name of Object cloud_e_pg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud EPg.
* `annotation` - (Optional) annotation for object cloud_e_pg.
* `exception_tag` - (Optional) exception_tag for object cloud_e_pg.
* `flood_on_encap` - (Optional) flood_on_encap for object cloud_e_pg.
* `match_t` - (Optional) match criteria
* `name_alias` - (Optional) name_alias for object cloud_e_pg.
* `pref_gr_memb` - (Optional) pref_gr_memb for object cloud_e_pg.
* `prio` - (Optional) qos priority class id
