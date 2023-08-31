---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_service_endpoint_selector"
sidebar_current: "docs-aci-resource-cloud_service_endpoint_selector"
description: |-
  Manages ACI Cloud Service Endpoint Selector
---

# aci_cloud_service_endpoint_selector #

Manages ACI Cloud Service Endpoint Selector

## API Information ##

* `Class` - cloudSvcEPSelector
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/svcepselector-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
resource "aci_cloud_service_endpoint_selector" "example" {
  cloud_service_epg_dn  = aci_cloud_service_epg.example.id
  name                  = "example"
  annotation            = "orchestrator:terraform"
  match_expression      = "IP=='7.1.0.0/24'"
}
```

## Argument Reference ##

* `cloud_service_epg_dn` - (Required) Distinguished name of the parent CloudServiceEPg object.
* `name` - (Required) Name of the Cloud Service Endpoint Selector object.
* `annotation` - (Optional) Annotation of the Cloud Service Endpoint Selector object.
* `name_alias` - (Optional) Name Alias of the Cloud Service Endpoint Selector object.
* `match_expression` - (Optional) Expression used to define matching tag.



## Importing ##

An existing CloudServiceEndpointSelector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_service_endpoint_selector.example <Dn>
```