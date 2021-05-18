---
layout: "aci"
page_title: "ACI: aci_l4_l7_service_graph_template"
sidebar_current: "docs-aci-data-source-l4_l7_service_graph_template"
description: |-
  Data source for ACI L4-L7 Service Graph Template
---

# aci_l4_l7_service_graph_template #
Data source for ACI L4-L7 Service Graph Template

## Example Usage ##

```hcl

data "aci_l4_l7_service_graph_template" "check" {
  tenant_dn = "${aci_tenant.tenentcheck.id}"
  name      = "first"
}

```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name for L4-L7 service graph template object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L4-L7 Service Graph Template.
* `annotation` - (Optional) Annotation for L4-L7 service graph template object.
* `name_alias` - (Optional) name_alias for L4-L7 service graph template object.
* `l4_l7_service_graph_template_type` - (Optional) Component type for L4-L7 service graph template object.
* `ui_template_type` - (Optional) UI template type for L4-L7 service graph template object.
