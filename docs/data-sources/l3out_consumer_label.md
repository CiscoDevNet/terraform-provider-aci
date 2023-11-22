---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_consumer_label"
sidebar_current: "docs-aci-data-source-l3out_consumer_label"
description: |-
  Data source for L3out Consumer Label
---

# aci_l3out_consumer_label #

Data source for L3out Consumer Label

## API Information ##

* `Class` - l3extConsLbl

* `Distinguished Name Formats`
  - uni/tn-{name}/out-{name}/conslbl-{name}

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs

## Example Usage ##

```hcl

data "aci_l3out_consumer_label" "example" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_l3out_consumer_label"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - `aci_l3_outside` (class: l3extOut)
* `name` - (string) The name of the L3out Consumer Label object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the L3out Consumer Label object.
* `annotation` - (string) The annotation of the L3out Consumer Label object.
* `description` - (string) The description of the L3out Consumer Label object.
* `name_alias` - (string) The name alias of the L3out Consumer Label object.
* `owner` - (string) The owner of the target relay. The DHCP relay is any host that forwards DHCP packets between clients and servers. The relays are used to forward requests and replies between clients and servers when they are not on the same physical subnet.
* `owner_key` - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `tag` - (string) Specifies the color of a policy label.

* `annotations` - (list) A list of Annotation objects (tagAnnotation).
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.