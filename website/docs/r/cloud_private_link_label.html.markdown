---
subcategory: -
layout: "aci"
page_title: "ACI: aci_private_link_labelfortheservice_epg"
sidebar_current: "docs-aci-resource-private_link_labelfortheservice_epg"
description: |-
  Manages ACI Private Link Label for the service EPg
---

# aci_private_link_labelfortheservice_epg #

Manages ACI Private Link Label for the service EPg

## API Information ##

* `Class` - cloudPrivateLinkLabel
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}/privatelinklabel-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_private_link_labelfortheservice_epg" "example" {
  cloud_service_epg_dn  = aci_cloud_service_epg.example.id
  name  = "example"
  annotation = "orchestrator:terraform"

  name_alias = 
}
```

## Argument Reference ##

* `cloud_service_epg_dn` - (Required) Distinguished name of the parent CloudServiceEPg object.
* `name` - (Required) Name of the Private Link Label for the service EPg object.
* `annotation` - (Optional) Annotation of the Private Link Label for the service EPg object.
* `name_alias` - (Optional) Name Alias of the Private Link Label for the service EPg object.



## Importing ##

An existing PrivateLinkLabelfortheserviceEPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_private_link_labelfortheservice_epg.example <Dn>
```