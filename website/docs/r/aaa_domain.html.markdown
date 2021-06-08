---
layout: "aci"
page_title: "ACI: aci_aaa_domain"
sidebar_current: "docs-aci-resource-aaa_domain"
description: |-
  Manages ACI aaa Domain
---

# aci_aaa_domain #

Manages ACI aaa Domain

## Example Usage ##

```hcl

resource "aci_aaa_domain" "example" {
  name        = "aaa_domain_1"
  description = "from terraform"
  annotation  = "tag_aaa"
  name_alias  = "alias_aaa"
}

```

## Argument Reference ##

* `name` - (Required) Name of object aaa domain.
* `description` - (Optional) Description for object aaa domain.
* `annotation` - (Optional) Annotation for object aaa domain.
* `name_alias` - (Optional) Name alias for object aaa domain.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the Dn of the aaa domain.

## Importing ##

An existing aaa domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_aaa_domain.example <Dn>
```
