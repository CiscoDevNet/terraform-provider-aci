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

data "aci_client_end_point" "example" {
  application_epg_dn  = "${aci_application_epg.example.id}"
  name                = "example"
}

```


## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `name` - (Required) name of Object client_end_point.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Client End Point.
* `annotation` - (Optional) annotation for object client_end_point.
* `client_end_point_id` - (Optional) object identifier
* `name_alias` - (Optional) name_alias for object client_end_point.
