---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_redistribute_policy"
sidebar_current: "docs-aci-data-source-l3out_redistribute_policy"
description: |-
  Data source for L3out Redistribute Policy
---

# aci_l3out_redistribute_policy #

Data source for L3out Redistribute Policy

## API Information ##

* `Class` - l3extRsRedistributePol

* `Distinguished Name Formats`
  - uni/tn-{name}/out-{name}/rsredistributePol-[{tnRtctrlProfileName}]-{src}

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Redistribute Policies

## Example Usage ##

```hcl

data "aci_l3out_redistribute_policy" "example" {
  parent_dn          = aci_l3_outside.example.id
  src                = "direct"
  route_profile_name = "test_tn_rtctrl_profile_name"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - `aci_l3_outside` (class: l3extOut)
* `src` - (string) The source IP address.
* `route_profile_name` - (string) The name of the route profile associated with this object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the L3out Redistribute Policy object.
* `annotation` - (string) The annotation of the L3out Redistribute Policy object.

* `annotations` - (list) A list of Annotation objects (tagAnnotation).
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.