---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_private_link_label"
sidebar_current: "docs-aci-data-source-cloud_private_link_label"
description: |-
  Data source for ACI Private Link Label
---

# aci_cloud_private_link_label #

Data source for ACI Private Link Label


## API Information ##

* `Class` - cloudPrivateLinkLabel
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/privatelinklabel-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
data "aci_private_link_labelfortheservice_epg" "example" {
  parent_dn  = aci_cloud_service_epg.example.id
  name                  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent Cloud Service EPG or Cloud Subnet objects.
* `name` - (Required) Name of the Private Link Label for the service EPG object.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Private Link Label.
* `annotation` - (Read-Only) Annotation of the Private Link Label.
* `name_alias` - (Read-Only) Name Alias of the Private Link Label.
