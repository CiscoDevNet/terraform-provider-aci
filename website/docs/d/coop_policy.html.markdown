---
layout: "aci"
page_title: "ACI: aci_coop_policy"
sidebar_current: "docs-aci-data-source-coop_policy"
description: |-
  Data source for ACI COOP Policy
---

# aci_coop_policy #

Data source for ACI COOP Policy


## API Information ##

* `Class` - coopPol
* `Distinguished Named` - uni/fabric/pol-{name}

## GUI Information ##

* `Location` - System -> System Settings -> COOP Group -> Policy



## Example Usage ##

```hcl
data "aci_coop_policy" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the COOP Group Policy.
* `annotation` - (Optional) Annotation of object COOP Group Policy.
* `name_alias` - (Optional) Name Alias of object COOP Group Policy.
* `type` - (Optional) Authentication type. The specific type of the object or component.
* `description` - (Optional) Description of object COOP Group Policy.