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
  name        = "example"
  annotation  = "tag_aaa"
  name_alias  = "alias_aaa"
}
```

## Argument Reference ##

* `name` - (Required) name of Object aaa domain.
* `annotation` - (Optional) annotation for object aaa domain.
* `name_alias` - (Optional) name_alias for object aaa domain.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the Dn of the aaa domain.

## Importing ##

An existing aaa domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_aaa_domain.example <Dn>
```
