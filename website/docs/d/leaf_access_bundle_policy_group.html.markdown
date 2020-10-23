---
layout: "aci"
page_title: "ACI: aci_leaf_access_bundle_policy_group"
sidebar_current: "docs-aci-data-source-leaf_access_bundle_policy_group"
description: |-
  Data source for ACI leaf access bundle policy group
---

# aci_leaf_access_bundle_policy_group #
Data source for ACI leaf access bundle policy group

## Example Usage ##

```hcl

data "aci_leaf_access_bundle_policy_group" "dev_pol_grp" {
  name  = "foo_pol_grp"
}

```


## Argument Reference ##
* `name` - (Required) name of Object leaf_access_bundle_policy_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the ACI leaf access bundle policy group.
* `annotation` - (Optional) annotation for object aci_leaf_access_bundle_policy_group.
* `lag_t` - (Optional) The bundled ports group link aggregation type: port channel vs virtual port channel.
* `name_alias` - (Optional) name_alias for object aci_leaf_access_bundle_policy_group.
