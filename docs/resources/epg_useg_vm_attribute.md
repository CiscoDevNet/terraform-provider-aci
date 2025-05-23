---
# Documentation generated by "gen/generator.go"; DO NOT EDIT.
# In order to regenerate this file execute `go generate` from the repository root.
# More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_useg_vm_attribute"
sidebar_current: "docs-aci-resource-aci_epg_useg_vm_attribute"
description: |-
  Manages ACI EPG uSeg VM Attribute
---

# aci_epg_useg_vm_attribute #

Manages ACI EPG uSeg VM Attribute



## API Information ##

* Class: [fvVmAttr](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvVmAttr/overview)

* Supported in ACI versions: 1.1(1j) and later.

* Distinguished Name Formats:
  - `uni/tn-{name}/ap-{name}/epg-{name}/crtrn/crtrn-{name}/vmattr-{name}`
  - `uni/tn-{name}/ap-{name}/epg-{name}/crtrn/vmattr-{name}`

## GUI Information ##

* Location: `Tenants -> Application Profiles -> uSeg EPGs -> uSeg Attributes`

## Example Usage ##

The configuration snippet below creates a EPG uSeg VM Attribute with only required attributes.

```hcl

resource "aci_epg_useg_vm_attribute" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  name      = "vm_attribute"
  value     = "default_value"
}

resource "aci_epg_useg_vm_attribute" "example_epg_useg_sub_block_statement" {
  parent_dn = aci_epg_useg_sub_block_statement.example.id
  name      = "vm_attribute"
  value     = "default_value"
}

```
The configuration snippet below shows all possible attributes of the EPG uSeg VM Attribute.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl

resource "aci_epg_useg_vm_attribute" "full_example_epg_useg_block_statement" {
  parent_dn   = aci_epg_useg_block_statement.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description_1"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias_1"
  operator    = "contains"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  type        = "domain"
  value       = "default_value"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}

resource "aci_epg_useg_vm_attribute" "full_example_epg_useg_sub_block_statement" {
  parent_dn   = aci_epg_useg_sub_block_statement.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description_1"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias_1"
  operator    = "contains"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  type        = "domain"
  value       = "default_value"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}

```

All examples for the EPG uSeg VM Attribute resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/resources/aci_epg_useg_vm_attribute) folder.

## Schema ##

### Required ###

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_epg_useg_block_statement](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/epg_useg_block_statement) ([fvCrtrn](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvCrtrn/overview))
  - [aci_epg_useg_sub_block_statement](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/epg_useg_sub_block_statement) ([fvSCrtrn](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvSCrtrn/overview))
* `name` (name) - (string) The name of the EPG uSeg VM Attribute object.
* `value` (value) - (string) The value of the EPG uSeg VM Attribute object.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the EPG uSeg VM Attribute object.

### Optional ###

* `annotation` (annotation) - (string) The annotation of the EPG uSeg VM Attribute object. This attribute is supported in ACI versions: 3.2(1l) and later.
  - Default: `orchestrator:terraform`
* `category` (category) - (string) The category of the EPG uSeg VM Attribute object. This attribute is supported in ACI versions: 2.3(1e) and later.
* `description` (descr) - (string) The description of the EPG uSeg VM Attribute object.
* `label_name` (labelName) - (string) The label name of the EPG uSeg VM Attribute object.
* `name_alias` (nameAlias) - (string) The name alias of the EPG uSeg VM Attribute object. This attribute is supported in ACI versions: 2.2(1k) and later.
* `operator` (operator) - (string) The operator of the EPG uSeg VM Attribute object.
  - Default: `equals`
  - Valid Values: `contains`, `endsWith`, `equals`, `notEquals`, `startsWith`.
* `owner_key` (ownerKey) - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` (ownerTag) - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `type` (type) - (string) The type of the EPG uSeg VM Attribute object.
  - Default: `vm-name`
  - Valid Values: `custom-label`, `domain`, `guest-os`, `hv`, `rootContName`, `tag`, `vm`, `vm-folder`, `vm-name`, `vmfolder-path`, `vnic`.
* `annotations` - (list) A list of Annotations (ACI object [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)). Annotations can also be configured using a separate [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  #### Required ####
  
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.
* `tags` - (list) A list of Tags (ACI object [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview)). Tags can also be configured using a separate [aci_tag](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tag) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  #### Required ####
  
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.

## Importing

An existing EPG uSeg VM Attribute can be [imported](https://www.terraform.io/docs/import/index.html) into this resource with its distinguished name (DN), via the following command:

```
terraform import aci_epg_useg_vm_attribute.example_epg_useg_block_statement uni/tn-{name}/ap-{name}/epg-{name}/crtrn/crtrn-{name}/vmattr-{name}
```

Starting in Terraform version 1.5, an existing EPG uSeg VM Attribute can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/ap-{name}/epg-{name}/crtrn/crtrn-{name}/vmattr-{name}"
  to = aci_epg_useg_vm_attribute.example_epg_useg_block_statement
}
```
