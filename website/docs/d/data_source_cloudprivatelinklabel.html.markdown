---
subcategory: -
layout: "aci"
page_title: "ACI: aci_private_link_labelfortheservice_epg"
sidebar_current: "docs-aci-data-source-private_link_labelfortheservice_epg"
description: |-
  Data source for ACI Private Link Label for the service EPg
---

# aci_private_link_labelfortheservice_epg #

Data source for ACI Private Link Label for the service EPg


## API Information ##

* `Class` - cloudPrivateLinkLabel
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/privatelinklabel-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
data "aci_private_link_labelfortheservice_epg" "example" {
  cloud_service_epg_dn  = aci_cloud_service_epg.example.id
  name  = "example"
}
```

## Argument Reference ##

* `cloud_service_epg_dn` - (Required) Distinguished name of the parent CloudServiceEPg object.
* `name` - (Required) Name of the Private Link Label for the service EPg object.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Private Link Label for the service EPg.
* `annotation` - (Read-Only) Annotation of the Private Link Label for the service EPg object.
* `name_alias` - (Read-Only) Name Alias of the Private Link Label for the service EPg object.
