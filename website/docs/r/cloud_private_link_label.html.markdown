---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_private_link_label"
sidebar_current: "docs-aci-resource-cloud_private_link_label"
description: |-
  Manages ACI Private Link Label 
---

# aci_cloud_private_link_label #

Manages ACI Private Link Label

## API Information ##

* `Class` - cloudPrivateLinkLabel
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/privatelinklabel-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
resource "aci_cloud_private_link_label" "example" {
  parent_dn  = aci_cloud_service_epg.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent Cloud Service EPG or Cloud Subnet object.
* `name` - (Required) Name of the Private Link Label.
* `annotation` - (Optional) Annotation of the Private Link Label.
* `name_alias` - (Optional) Name Alias of the Private Link Label.



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