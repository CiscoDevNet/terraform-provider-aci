---
layout: "aci"
page_title: "ACI: aci_aaa_domain"
sidebar_current: "docs-aci-data-source-aaa_domain"
description: |-
  Data source for ACI aaa Domain
---

# aci_aaa_domain #
Data source for ACI aaa Domain

## Example Usage ##

```hcl

data "aci_aaa_domain" "example" {
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object aaa domain.



## Attribute Reference

* `id` - Attribute id set to the Dn of the aaa domain.
* `annotation` - (Optional) annotation for object aaa domain.
* `name_alias` - (Optional) name_alias for object aaa domain.
