---
layout: "aci"
page_title: "ACI: aci_contract_provider_consumer"
sidebar_current: "docs-aci-resource-contract_provider_consumer"
description: |-
  Manages ACI Contract Provider and Consumer.
---

# aci_contract_provider_consumer #

Manages ACI Contract Provider and Consumer.

## Example Usage ##

```hcl
resource "aci_contract_consumer_provider" "example" {
  application_epg_dn  = "${aci_application_epg.inherit_epg.id}"
  contract_name  = "contract1"
  contract_type = "consumer"
  annotation  = "con"
  prio  = "level3"
  match_t = "All"
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `contract_name` - (Required) Name of Contract to be created.
* `contract_type` - (Required) Type of contract to be created. Allowed values are "consumer" and "provider'.
* `annotation` - (Optional) annotation for object contract_provider_consumer.
* `match_t` - (Optional) match type. Allowed values are "All", "AtleastOne", "AtmostOne" and "None".
* `prio` - (Optional) service priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6" and "unspecified".

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract Provider.

## Importing ##

An existing Contract Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_contract_provider.example <Dn>
```