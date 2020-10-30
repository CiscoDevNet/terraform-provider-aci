---
layout: "aci"
page_title: "ACI: aci_client_end_point"
sidebar_current: "docs-aci-data-source-client_end_point"
description: |-
  Data source for ACI Client End Point
---

# aci_client_end_point #
Data source for ACI Client End Point

## Example Usage ##

```hcl

data "aci_client_end_point" "check" {
  application_epg_dn  = "${aci_application_epg.epg.id}"
  mac                 = "25:56:68:78:98:74"
  ip                  = "1.2.3.4"
  vlan                = "5"
}

```


## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `name` - (Optional) name of Object client end point.
* `mac` - (Optional) Mac address of the object client end point.
* `ip` - (Optional) ip address of the object client end point.
* `vlan` - (Optional) vlan for the object client end point.



## Attribute Reference

* `id` - Attribute id set as all Dns for matching the Client End Point.
* `object_dns` - list of all Dns which matched to the given filter attributes.
