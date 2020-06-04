---
layout: "aci"
page_title: "ACI: aci_contract_provider_consumer"
sidebar_current: "docs-aci-data-source-contract_provider_consumer"
description: |-
  Data source for ACI Contract Provider and Consumer.
---

# aci_contract_provider_consumer #

Data source for ACI Contract Provider and Consumer.

## Example Usage ##

```hcl
data "aci_contract_consumer_provider" "example" {
  application_epg_dn  = "${aci_application_epg.inherit_epg.id}"
  contract_name  = "example"
  contract_type = "provider"
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `contract_name` - (Required) Name of Contract to be created.
* `contract_type` - (Required) Type of contract to be created. Allowed values are "consumer" and "provider'.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Contract Provider.
* `annotation` - (Optional) annotation for object contract_provider_consumer.
* `match_t` - (Optional) match type. Allowed values are "All", "AtleastOne", "AtmostOne" and "None".
* `prio` - (Optional) service priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6" and "unspecified".
