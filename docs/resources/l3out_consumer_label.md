---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_consumer_label"
sidebar_current: "docs-aci-resource-l3out_consumer_label"
description: |-
  Manages ACI L3out Consumer Label
---

# aci_l3out_consumer_label #

Manages ACI L3out Consumer Label

## API Information ##

* `Class` - `l3extConsLbl`

* `Distinguished Name Formats`
  - `uni/tn-{name}/out-{name}/conslbl-{name}`

## GUI Information ##

* `Location` - `Tenants -> Networking -> L3Outs`

## Example Usage ##

```hcl

resource "aci_l3out_consumer_label" "example" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

```

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_l3_outside](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside) (`l3extOut`)
* `name` - (string) The name of the L3out Consumer Label object.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Consumer Label object.

### Optional
  
* `annotation` - (string) The annotation of the L3out Consumer Label object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the L3out Consumer Label object.
* `name_alias` - (string) The name alias of the L3out Consumer Label object.
* `owner` - (string) The owner of the target relay. The DHCP relay is any host that forwards DHCP packets between clients and servers. The relays are used to forward requests and replies between clients and servers when they are not on the same physical subnet.
  - Default: `infra`
  - Valid Values: `infra`, `tenant`.
* `owner_key` - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `tag` - (string) Specifies the color of a policy label.
  - Valid Values: `alice-blue`, `antique-white`, `aqua`, `aquamarine`, `azure`, `beige`, `bisque`, `black`, `blanched-almond`, `blue`, `blue-violet`, `brown`, `burlywood`, `cadet-blue`, `chartreuse`, `chocolate`, `coral`, `cornflower-blue`, `cornsilk`, `crimson`, `cyan`, `dark-blue`, `dark-cyan`, `dark-goldenrod`, `dark-gray`, `dark-green`, `dark-khaki`, `dark-magenta`, `dark-olive-green`, `dark-orange`, `dark-orchid`, `dark-red`, `dark-salmon`, `dark-sea-green`, `dark-slate-blue`, `dark-slate-gray`, `dark-turquoise`, `dark-violet`, `deep-pink`, `deep-sky-blue`, `dim-gray`, `dodger-blue`, `fire-brick`, `floral-white`, `forest-green`, `fuchsia`, `gainsboro`, `ghost-white`, `gold`, `goldenrod`, `gray`, `green`, `green-yellow`, `honeydew`, `hot-pink`, `indian-red`, `indigo`, `ivory`, `khaki`, `lavender`, `lavender-blush`, `lawn-green`, `lemon-chiffon`, `light-blue`, `light-coral`, `light-cyan`, `light-goldenrod-yellow`, `light-gray`, `light-green`, `light-pink`, `light-salmon`, `light-sea-green`, `light-sky-blue`, `light-slate-gray`, `light-steel-blue`, `light-yellow`, `lime`, `lime-green`, `linen`, `magenta`, `maroon`, `medium-aquamarine`, `medium-blue`, `medium-orchid`, `medium-purple`, `medium-sea-green`, `medium-slate-blue`, `medium-spring-green`, `medium-turquoise`, `medium-violet-red`, `midnight-blue`, `mint-cream`, `misty-rose`, `moccasin`, `navajo-white`, `navy`, `old-lace`, `olive`, `olive-drab`, `orange`, `orange-red`, `orchid`, `pale-goldenrod`, `pale-green`, `pale-turquoise`, `pale-violet-red`, `papaya-whip`, `peachpuff`, `peru`, `pink`, `plum`, `powder-blue`, `purple`, `red`, `rosy-brown`, `royal-blue`, `saddle-brown`, `salmon`, `sandy-brown`, `sea-green`, `seashell`, `sienna`, `silver`, `sky-blue`, `slate-blue`, `slate-gray`, `snow`, `spring-green`, `steel-blue`, `tan`, `teal`, `thistle`, `tomato`, `turquoise`, `violet`, `wheat`, `white`, `white-smoke`, `yellow`, `yellow-green`.

* `annotations` - (list) A list of Annotations objects `tagAnnotation` which can be configured using the [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource.
  
  #### Required
  
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.

## Importing

An existing L3out Consumer Label can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_l3out_consumer_label.example uni/tn-{name}/out-{name}/conslbl-{name}
```

Starting in Terraform version 1.5, an existing L3out Consumer Label can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/out-{name}/conslbl-{name}"
  to = aci_l3out_consumer_label.example
}
```

## Child Resources
  
  - [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation)