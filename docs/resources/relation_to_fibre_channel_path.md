---
# Documentation generated by "gen/generator.go"; DO NOT EDIT.
# In order to regenerate this file execute `go generate` from the repository root.
# More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_relation_to_fibre_channel_path"
sidebar_current: "docs-aci-resource-aci_relation_to_fibre_channel_path"
description: |-
  Manages ACI Relation To Fibre Channel Path
---

# aci_relation_to_fibre_channel_path #

Manages ACI Relation To Fibre Channel Path



## API Information ##

* Class: [fvRsFcPathAtt](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvRsFcPathAtt/overview)

* Supported in ACI versions: 2.0(1m) and later.

* Distinguished Name Format: `uni/tn-{name}/ap-{name}/epg-{name}/rsfcPathAtt-[{tDn}]`

## GUI Information ##

* Location: `Tenants -> Application Profiles -> Application EPGs -> Fibre Channel (Paths)`

## Example Usage ##

The configuration snippet below creates a Relation To Fibre Channel Path with only required attributes.

```hcl

resource "aci_relation_to_fibre_channel_path" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
}

```
The configuration snippet below shows all possible attributes of the Relation To Fibre Channel Path.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl

resource "aci_relation_to_fibre_channel_path" "full_example_application_epg" {
  parent_dn   = aci_application_epg.example.id
  annotation  = "annotation"
  description = "description_1"
  target_dn   = "topology/pod-1/paths-101/pathep-[eth1/1]"
  vsan        = "vsan-10"
  vsan_mode   = "native"
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

All examples for the Relation To Fibre Channel Path resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/resources/aci_relation_to_fibre_channel_path) folder.

## Schema ##

### Required ###

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_application_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/application_epg) ([fvAEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvAEPg/overview))
* `target_dn` (tDn) - (string) The distinguished name of the target.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the Relation To Fibre Channel Path object.

### Optional ###
  
* `annotation` (annotation) - (string) The annotation of the Relation To Fibre Channel Path object.
  - Default: `orchestrator:terraform`
* `description` (descr) - (string) The description of the Relation To Fibre Channel Path object.
* `vsan` (vsan) - (string) The virtual storage area network (VSAN) of the Relation To Fibre Channel Path object.
  - Default: `unknown`
* `vsan_mode` (vsanMode) - (string) The virtual storage area network (VSAN) mode of the Relation To Fibre Channel Path object.
  - Default: `regular`
  - Valid Values: `native`, `regular`.

* `annotations` - (list) A list of Annotations (ACI object [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)). Annotations can also be configured using a separate [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  
  #### Required ####
  
  * `key` (key) - (string) The key used to uniquely identify this configuration object.
  * `value` (value) - (string) The value of the property.

* `tags` - (list) A list of Tags (ACI object [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview)). Tags can also be configured using a separate [aci_tag](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tag) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  
  #### Required ####
  
  * `key` (key) - (string) The key used to uniquely identify this configuration object.
  * `value` (value) - (string) The value of the property.

## Importing

An existing Relation To Fibre Channel Path can be [imported](https://www.terraform.io/docs/import/index.html) into this resource with its distinguished name (DN), via the following command:

```
terraform import aci_relation_to_fibre_channel_path.example_application_epg uni/tn-{name}/ap-{name}/epg-{name}/rsfcPathAtt-[{tDn}]
```

Starting in Terraform version 1.5, an existing Relation To Fibre Channel Path can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/ap-{name}/epg-{name}/rsfcPathAtt-[{tDn}]"
  to = aci_relation_to_fibre_channel_path.example_application_epg
}
```