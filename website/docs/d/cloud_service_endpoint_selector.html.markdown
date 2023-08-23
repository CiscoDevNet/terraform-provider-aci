---
subcategory: -
layout: "aci"
page_title: "ACI: aci_cloud_service_endpoint_selector"
sidebar_current: "docs-aci-data-source-cloud_service_endpoint_selector"
description: |-
  Data source for ACI Cloud Service Endpoint Selector
---

# aci_cloud_service_endpoint_selector #

Data source for ACI Cloud Service Endpoint Selector
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudSvcEPSelector
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/svcepselector-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
data "aci_cloud_service_endpoint_selector" "example" {
  cloud_service_epg_dn  = aci_cloud_service_epg.example.id
  name  = "example"
}
```

## Argument Reference ##

* `cloud_service_epg_dn` - (Required) Distinguished name of the parent CloudServiceEPg object.
* `name` - (Required) Name of the Cloud Service Endpoint Selector object.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud Service Endpoint Selector.
* `annotation` - (Read-Only) Annotation of the Cloud Service Endpoint Selector object.
* `name_alias` - (Read-Only) Name Alias of the Cloud Service Endpoint Selector object.
* `match_expression` - (Read-Only) Expression used to define matching tagTags. 
