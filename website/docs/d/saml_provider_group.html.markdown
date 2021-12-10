---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_saml_provider_group"
sidebar_current: "docs-aci-data-source-saml_provider_group"
description: |-
  Data source for ACI SAML Provider Group
---

# aci_saml_provider_group #

Data source for ACI SAML Provider Group


## API Information ##

* `Class` - aaaSamlProviderGroup
* `Distinguished Named` - uni/userext/samlext/samlprovidergroup-{name}


## Example Usage ##

```hcl
data "aci_saml_provider_group" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object SAML Provider Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the SAML Provider Group.
* `annotation` - (Optional) Annotation of object SAML Provider Group.
* `name_alias` - (Optional) Name Alias of object SAML Provider Group.
* `description` - (Optional) Description of object SAML Provider Group.