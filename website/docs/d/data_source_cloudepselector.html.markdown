---
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selector"
sidebar_current: "docs-aci-data-source-cloud_endpoint_selector"
description: |-
  Data source for ACI Cloud Endpoint Selector
---

# aci_cloud_endpoint_selector #
Data source for ACI Cloud Endpoint Selector
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_endpoint_selector" "example" {

  cloud_e_pg_dn  = "${aci_cloud_e_pg.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `cloud_e_pg_dn` - (Required) Distinguished name of parent CloudEPg object.
* `name` - (Required) name of Object cloud_endpoint_selector.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Endpoint Selector.
* `annotation` - (Optional) annotation for object cloud_endpoint_selector.
* `match_expression` - (Optional) match_expression for object cloud_endpoint_selector.
* `name_alias` - (Optional) name_alias for object cloud_endpoint_selector.
