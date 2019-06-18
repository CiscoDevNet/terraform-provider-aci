---
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selectorfor_external_e_pgs"
sidebar_current: "docs-aci-data-source-cloud_endpoint_selectorfor_external_e_pgs"
description: |-
  Data source for ACI Cloud Endpoint Selector for External EPgs
---

# aci_cloud_endpoint_selectorfor_external_e_pgs #
Data source for ACI Cloud Endpoint Selector for External EPgs
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_endpoint_selectorfor_external_e_pgs" "example" {

  cloud_external_e_pg_dn  = "${aci_cloud_external_e_pg.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `cloud_external_e_pg_dn` - (Required) Distinguished name of parent CloudExternalEPg object.
* `name` - (Required) name of Object cloud_endpoint_selectorfor_external_e_pgs.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Endpoint Selector for External EPgs.
* `annotation` - (Optional) annotation for object cloud_endpoint_selectorfor_external_e_pgs.
* `is_shared` - (Optional) is_shared for object cloud_endpoint_selectorfor_external_e_pgs.
* `match_expression` - (Optional) match_expression for object cloud_endpoint_selectorfor_external_e_pgs.
* `name_alias` - (Optional) name_alias for object cloud_endpoint_selectorfor_external_e_pgs.
* `subnet` - (Optional) subnet for object cloud_endpoint_selectorfor_external_e_pgs.
