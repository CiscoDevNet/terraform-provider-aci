---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_mgmt_preference"
sidebar_current: "docs-aci-resource-aci_mgmt_preference"
description: |-
  Manages ACI Mgmt Preference
---

# aci_mgmt_preference #
Manages ACI Mgmt Preference

## API Information ##
* `Class` - mgmtConnectivityPrefs
* `Distinguished Name` - uni/fabric/connectivityPrefs

## GUI Information ##
* `Location` - System -> System Settings -> APIC Connectivity Profile -> Policy

## Example Usage ##
```hcl
resource "aci_mgmt_preference" "example" {
  interface_pref = "inband"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
}
```

## NOTE ##
Users can use the resource of type `aci_mgmt_preference` to change the configuration of the object Mgmt Preference. Users cannot create more than one instance of object Mgmt Preference.

## Argument Reference ##
* `interface_pref` - (Optional) Management interface that has to be used. Allowed values are "inband" and "ooband". 
* `annotation` - (Optional) Annotation of object Mgmt Preference.
* `description` - (Optional) Description of object Mgmt Preference.
* `name_alias` - (Optional) Name Alias of object Mgmt Preference.

## Importing ##

An existing Mgmt Preference can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_mgmt_preference.example <Dn>
```