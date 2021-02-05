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
* `tenant_dn` - (Required) distinguished name of parent Tenant object.
* `name` - (Required) name for l4 l7 service graph template object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L4-L7 Service Graph Template.
* `annotation` - (Optional) annotation for l4 l7 service graph template object.
* `name_alias` - (Optional) name_alias for l4 l7 service graph template object.
* `l4_l7_service_graph_template_type` - (Optional) component type for l4 l7 service graph template object.
* `ui_template_type` - (Optional) UI template type for l4 l7 service graph template object.
