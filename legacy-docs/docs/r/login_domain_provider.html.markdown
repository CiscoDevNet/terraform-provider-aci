---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_login_domain_provider"
sidebar_current: "docs-aci-resource-aci_login_domain_provider"
description: |-
  Manages ACI Login Domain Provider
---

# aci_login_domain_provider #
Manages ACI Login Domain Provider

## API Information ##
* `Class` - aaaProviderRef
* `Supported Distinguished Names`<br>
[1] uni/userext/duoext/duoprovidergroup-{name}/providerref-{name}<br>
[2] uni/userext/rsaext/rsaprovidergroup-{name}/providerref-{name}<br>
[3] uni/userext/samlext/samlprovidergroup-{name}/providerref-{name}<br>
[4] uni/userext/tacacsext/tacacsplusprovidergroup-{name}/providerref-{name}<br>
[5] uni/userext/radiusext/radiusprovidergroup-{name}/providerref-{name}<br>
[6] uni/userext/ldapext/ldapprovidergroup-{name}/providerref-{name}<br>


## Example Usage ##

```hcl
resource "aci_login_domain_provider" "example" {
  parent_dn  = aci_duo_provider_group.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  order = "0"
  name_alias = "example_name_alias"
  description = "from terraform"
}
```

## Argument Reference ##
* `parent_dn` - (Required) Distinguished name of parent.
* `name` - (Required) Name of object Login Domain Provider.
* `annotation` - (Optional) Annotation of object Login Domain Provider.
* `name_alias` - (Optional) Name Alias of object Login Domain Provider.
* `description` - (Optional) Description of object Login Domain Provider.
* `order` - (Optional) Order in which Providers are Tried. The relative priority in which the AAA provider will be contacted within the provider group. Allowed Range: "0" - "16". Allowed value "lowest-available". Default value is "0". (NOTE: "lowest-available" will set lowest value of order and will be translated by postConfig code to the numeric order value of 0.)

## Importing ##

An existing Login Domain Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_login_domain_provider.example <Dn>
```
