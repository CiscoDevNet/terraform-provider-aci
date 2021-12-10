---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_saml_provider_group"
sidebar_current: "docs-aci-resource-saml_provider_group"
description: |-
  Manages ACI SAML Provider Group
---

# aci_saml_provider_group #

Manages ACI SAML Provider Group

## API Information ##

* `Class` - aaaSamlProviderGroup
* `Distinguished Named` - uni/userext/samlext/samlprovidergroup-{name}


## Example Usage ##

```hcl
resource "aci_saml_provider_group" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  name_alias  = "saml_provider_group_alias"
  description = "From Terraform"
}
```

## Argument Reference ##


* `name` - (Required) Name of object SAML Provider Group.
* `annotation` - (Optional) Annotation of object SAML Provider Group.
* `name_alias` - (Optional) Name Alias of object SAML Provider Group.
* `description` - (Optional) Description of object SAML Provider Group.


## Importing ##

An existing SAMLProviderGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_saml_provider_group.example <Dn>
```