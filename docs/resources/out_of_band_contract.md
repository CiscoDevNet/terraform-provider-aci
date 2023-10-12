---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: out_of_band_contract"
sidebar_current: "docs-aci-resource-out_of_band_contract"
description: |-
  Manages ACI Out Of Band Contract
---

# aci_out_of_band_contract #

Manages ACI Out Of Band Contract

## API Information ##

* `Class` - [vzOOBBrCP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vzOOBBrCP/overview)

* `Distinguished Name Formats`
  - `uni/tn-mgmt/oobbrc-{name}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> Contracts -> Out-Of-Band Contracts`

## Example Usage ##

The configuration snippet below creates a Out Of Band Contract with only required attributes.

```hcl
resource "aci_out_of_band_contract" "example" {
  name = "test_name"
}
```

The configuration snippet below below shows all possible attributes of the Out Of Band Contract.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl
resource "aci_out_of_band_contract" "example" {
  annotation  = "annotation"
  description = "description"
  intent      = "estimate_add"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  priority    = "level1"
  scope       = "application-profile"
  target_dscp = "AF11"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
```

All examples for the Out Of Band Contract resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/examples/resources/aci_out_of_band_contract) folder.

## Schema

### Required

* `name` - (string) The name of the Out Of Band Contract object.

### Read-Only

* `id` - (string) The distinguished name (DN) of the Out Of Band Contract object.

### Optional
  
* `annotation` - (string) The annotation of the Out Of Band Contract object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the Out Of Band Contract object.
* `intent` - (string) Install Rules or Estimate Nummber of Rules.
  - Default: `install`
  - Valid Values: `estimate_add`, `estimate_delete`, `install`.
* `name_alias` - (string) The name alias of the Out Of Band Contract object.
* `owner_key` - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `priority` - (string) null.
  - Default: `unspecified`
  - Valid Values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.
* `scope` - (string) Represents the scope of this contract. If the scope is set as application-profile, the epg can only communicate with epgs in the same application-profile.
  - Default: `context`
  - Valid Values: `application-profile`, `context`, `global`, `tenant`.
* `target_dscp` - (string) contract level dscp value.
  - Default: `unspecified`
  - Valid Values: `AF11`, `AF12`, `AF13`, `AF21`, `AF22`, `AF23`, `AF31`, `AF32`, `AF33`, `AF41`, `AF42`, `AF43`, `CS0`, `CS1`, `CS2`, `CS3`, `CS4`, `CS5`, `CS6`, `CS7`, `EF`, `VA`, `unspecified`.

* `annotations` - (list) A list of Annotations objects ([tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)) which can be configured using the [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource.
  
  #### Required
  
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.

## Importing

An existing Out Of Band Contract can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_out_of_band_contract.example uni/tn-mgmt/oobbrc-{name}
```

Starting in Terraform version 1.5, an existing Out Of Band Contract can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-mgmt/oobbrc-{name}"
  to = aci_out_of_band_contract.example
}
```
