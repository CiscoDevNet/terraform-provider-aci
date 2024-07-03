---
# Documentation generated by "gen/generator.go"; DO NOT EDIT.
# In order to regenerate this file execute `go generate` from the repository root.
# More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_relation_to_contract_master"
sidebar_current: "docs-aci-data-source-aci_relation_to_contract_master"
description: |-
  Data source for Relation To Contract Master
---

# aci_relation_to_contract_master #

Data source for Relation To Contract Master

## API Information ##

* Class: [fvRsSecInherited](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvRsSecInherited/overview)

* Supported in ACI versions: 2.3(1e) and later.

* Distinguished Name Formats:
  - Too many DN formats to display, see model documentation for all possible parents of [fvRsSecInherited](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvRsSecInherited/overview).
  - `uni/tn-{name}/ap-{name}/epg-{name}/rssecInherited-[{tDn}]`
  - `uni/tn-{name}/ap-{name}/esg-{name}/rssecInherited-[{tDn}]`
  - `uni/tn-{name}/l2out-{name}/instP-{name}/rssecInherited-[{tDn}]`
  - `uni/tn-{name}/out-{name}/instP-{name}/rssecInherited-[{tDn}]`

## GUI Information ##

* Locations:
  - `Tenants -> Application Profiles -> Application EPGs`
  - `Tenants -> Application Profiles -> Endpoint Security Groups`
  - `Tenants -> Networking -> L3Outs -> External EPGs`
  - `Tenants -> Networking -> L2Outs -> External EPGs`

## Example Usage ##

```hcl

data "aci_relation_to_contract_master" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  target_dn = aci_application_epg.example_2.id
}

data "aci_relation_to_contract_master" "example_endpoint_security_group" {
  parent_dn = aci_endpoint_security_group.example.id
  target_dn = aci_endpoint_security_group.example_2.id
}

```

## Schema ##

### Required ###

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_cloud_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/cloud_epg) ([cloudEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/cloudEPg/overview))
  - [aci_cloud_external_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/cloud_external_epg) ([cloudExtEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/cloudExtEPg/overview))
  - [aci_cloud_service_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/cloud_service_epg) ([cloudSvcEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/cloudSvcEPg/overview))
  - [aci_application_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/application_epg) ([fvAEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvAEPg/overview))
  - [aci_endpoint_security_group](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/endpoint_security_group) ([fvESg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvESg/overview))
  - [aci_l2out_extepg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l2out_extepg) ([l2extInstP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l2extInstP/overview))
  - [aci_external_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/external_network_instance_profile) ([l3extInstP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extInstP/overview))
  - [aci_node_mgmt_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/node_mgmt_epg) ([mgmtInB](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtInB/overview))
  - The distinguished name (DN) of classes below can be used but currently there is no available resource for it:
    - [cloudISvcEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/cloudISvcEPg/overview)
    - [dhcpCRelPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/dhcpCRelPg/overview)
    - [dhcpPRelPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/dhcpPRelPg/overview)
    - [fvIntEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvIntEPg/overview)
    - [fvTnlEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTnlEPg/overview)
    - [infraCEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/infraCEPg/overview)
    - [infraPEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/infraPEPg/overview)
    - [l3extInstPDef](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extInstPDef/overview)
    - [vnsEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vnsEPpInfo/overview)
    - [vnsREPpInfo](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vnsREPpInfo/overview)
    - [vnsSDEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vnsSDEPpInfo/overview)
    - [vnsSHEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vnsSHEPpInfo/overview)

* `target_dn` (tDn) - (string) The distinguished name of the target.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the Relation To Contract Master object.
* `annotation` (annotation) - (string) The annotation of the Relation To Contract Master object.

* `annotations` - (list) A list of Annotations (ACI object [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)). This attribute is supported in ACI versions: 3.2(1l) and later.
  * `key` (key) - (string) The key used to uniquely identify this configuration object.
  * `value` (value) - (string) The value of the property.

* `tags` - (list) A list of Tags (ACI object [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview)). This attribute is supported in ACI versions: 3.2(1l) and later.
  * `key` (key) - (string) The key used to uniquely identify this configuration object.
  * `value` (value) - (string) The value of the property.