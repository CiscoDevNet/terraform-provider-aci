---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_aaa_domain"
sidebar_current: "docs-aci-data-source-aci_aaa_domain"
description: |-
  Data source for ACI aaa Domain
---

# aci_aaa_domain #
Data source for ACI aaa Domain

## Example Usage ##

```hcl

data "aci_aaa_domain" "example" {
  name  = "aaa_domain_1"
}

```


## Argument Reference ##
* `name` - (Required) Name of object aaa domain.

## Attribute Reference

* `id` - Attribute id set to the Dn of the aaa domain.
* `description` - (Optional) Description for object aaa domain.
* `annotation` - (Optional) Annotation for object aaa domain.
* `name_alias` - (Optional) Name alias for object aaa domain.
