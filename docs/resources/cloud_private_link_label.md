---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_private_link_label"
sidebar_current: "docs-aci-resource-cloud_private_link_label"
description: |-
  Manages ACI Cloud Private Link Label 
---

# aci_cloud_private_link_label #

Manages ACI Cloud Private Link Label
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudPrivateLinkLabel
* `Distinguished Name` - uni/tn-{tenant_name}/cloudapp-{application_name}/cloudsvcepg-{cloud_service_epg_name}/privatelinklabel-{name}

## GUI Information ##

* `Location` 
  - Application Management -> EPGs -> Create EPG
  - Application Management -> Application Profiles -> Create Application Profile


## Example Usage ##

```hcl
resource "aci_cloud_private_link_label" "example" {
  parent_dn  = aci_cloud_service_epg.example.id
  name       = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent Cloud Service EPG or Cloud Subnet object. Type: String.
* `name` - (Required) Name of the Cloud Private Link Label. Type: String.
* `annotation` - (Optional) Annotation of the Cloud Private Link Label. Type: String.
* `name_alias` - (Optional) Name Alias of the Cloud Private Link Label. Type: String.



## Importing ##

An existing Private Link Label can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_private_link_label.example <Dn>
```

Starting in Terraform version 1.5, an existing Private Link Label can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "<Dn>"
  to = aci_cloud_private_link_label.example
}
```