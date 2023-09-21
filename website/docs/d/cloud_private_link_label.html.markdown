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
* `Distinguished Name` - uni/tn-{tenant_name}/cloudapp-{application_name}/cloudsvcepg-{cloud_service_epg_name}/privatelinklabel-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
data "aci_cloud_private_link_label" "example" {
  parent_dn  = aci_cloud_service_epg.example.id
  name       = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent Cloud Service EPG or Cloud Subnet objects. Type: String.
* `name` - (Required) Name of the Private Link Label for the Cloud Service EPG object. Type: String.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Private Link Label. Type: String.
* `annotation` - (Read-Only) Annotation of the Private Link Label. Type: String.
* `name_alias` - (Read-Only) Name Alias of the Private Link Label. Type: String.
