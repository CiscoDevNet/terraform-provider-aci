---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_service_endpoint_selector"
sidebar_current: "docs-aci-resource-cloud_service_endpoint_selector"
description: |-
  Manages ACI Cloud Service Endpoint Selector
---

# aci_cloud_service_endpoint_selector #

Manages ACI Cloud Service Endpoint Selector.
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudSvcEPSelector
* `Distinguished Name` - uni/tn-{tenant_name}/cloudapp-{application_name}/cloudsvcepg-{cloud_service_epg_name}/svcepselector-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs -> Actions -> Create EPG


## Example Usage ##

```hcl
resource "aci_cloud_service_endpoint_selector" "example" {
  cloud_service_epg_dn  = aci_cloud_service_epg.example.id
  name                  = "example"
  match_expression      = "IP=='11.11.11.0/24'"
}
```

## Argument Reference ##

* `cloud_service_epg_dn` - (Required) Distinguished name of the parent Cloud Service EPG object. Type: String.
* `name` - (Required) Name of the Cloud Service Endpoint Selector object. Type: String.
* `annotation` - (Optional) Annotation of the Cloud Service Endpoint Selector object. Type: String.
* `name_alias` - (Optional) Name Alias of the Cloud Service Endpoint Selector object. Type: String.
* `match_expression` - (Optional) Expression used to define matching tag. Type: String.



## Importing ##

An existing Cloud Service Endpoint Selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_service_endpoint_selector.example <Dn>
```

Starting in Terraform version 1.5, an existing Cloud Service Endpoint Selector can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "<Dn>"
  to = aci_cloud_service_endpoint_selector.example
}
```