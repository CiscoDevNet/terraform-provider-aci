---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_filter"
sidebar_current: "docs-aci-resource-aci_filter"
description: |-
  Manages ACI Filter
---

# aci_filter #
Manages ACI Filter

## Example Usage ##

```hcl

resource "aci_filter" "example" {
  tenant_dn   = aci_tenant.dev_tenant.id
  description = "From Terraform"
  name        = "demo_filter"
  annotation  = "tag_filter"
  name_alias  = "alias_filter"
}

```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object filter.
* `description` - (Optional) Description for object filter.
* `annotation` - (Optional) Annotation for object filter.
* `name_alias` - (Optional) Name alias for object filter.
* `relation_vz_rs_filt_graph_att` - (Optional) Relation to class vnsInTerm. Type: String.
                
* `relation_vz_rs_fwd_r_flt_p_att` - (Optional) **Deprecated** Relation to class vzAFilterableUnit. Type: String.      
* `relation_vz_rs_rev_r_flt_p_att` - (Optional) **Deprecated** Relation to class vzAFilterableUnit. Type: String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the Filter.

## Importing ##

An existing Filter can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_filter.example <Dn>
```

Starting in Terraform version 1.5, an existing Filter can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

import {
  id = "<Dn>"
  to = aci_filter.example
}
