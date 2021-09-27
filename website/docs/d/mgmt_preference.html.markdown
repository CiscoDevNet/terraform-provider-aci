---
layout: "aci"
page_title: "ACI: aci_mgmt_preference"
sidebar_current: "docs-aci-data-source-mgmt_preference"
description: |-
  Data source for ACI Mgmt Preference
---

# aci_mgmt_preference #
Data source for ACI Mgmt Preference

## API Information ##
* `Class` - mgmtConnectivityPrefs
* `Distinguished Named` - uni/fabric/connectivityPrefs

## GUI Information ##
* `Location` - System -> System Settings -> APIC Connectivity Profile -> Policy

## Example Usage ##

```hcl
data "aci_mgmt_preference" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Mgmt preference.
* `interface_pref` - (Optional) Management interface that has to be used.
* `annotation` - (Optional) Annotation of object Mgmt preference.
* `name_alias` - (Optional) Name Alias of object Mgmt preference.
* `description` - (Optional) Description of object Mgmt Preference.
