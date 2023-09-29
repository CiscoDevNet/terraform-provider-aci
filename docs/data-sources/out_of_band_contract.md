---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: out_of_band_contract"
sidebar_current: "docs-aci-data-source-out_of_band_contract"
description: |-
  Data source for Out Of Band Contract
---

# aci_out_of_band_contract #

Data source for Out Of Band Contract

## API Information ##

* `Class` - `vzOOBBrCP`

* `Distinguished Name Formats`
  - `uni/tn-mgmt/oobbrc-{name}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> Contracts -> Out-Of-Band Contracts`

## Example Usage ##

```hcl
data "aci_out_of_band_contract" "example" {
  name = "test_name"
}
```

## Schema

### Required

* `name` - (string) The name of the Out Of Band Contract object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the Out Of Band Contract object.
* `annotation` - (string) The annotation of the Out Of Band Contract object.
* `description` - (string) The description of the Out Of Band Contract object.
* `intent` - (string) Install Rules or Estimate Nummber of Rules.
* `name_alias` - (string) The name alias of the Out Of Band Contract object.
* `owner_key` - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `priority` - (string) null.
* `scope` - (string) Represents the scope of this contract. If the scope is set as application-profile, the epg can only communicate with epgs in the same application-profile.
* `target_dscp` - (string) contract level dscp value.

* `annotations` - (list) A list of Annotations objects `tagAnnotation`.
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.