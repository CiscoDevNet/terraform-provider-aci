---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: contract_interface"
sidebar_current: "docs-aci-data-source-contract_interface"
description: |-
  Data source for Contract Interface
---

# aci_contract_interface #

Data source for Contract Interface

## API Information ##

* `Class` - [fvRsConsIf](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/fvRsConsIf/overview)

* `Distinguished Name Formats`
  - `outdefcont-{id}/rsoutDefContToOut-[{tDn}]/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/epp/fv-[{epgPKey}]/ac-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/epp/rtd-[{epgPKey}]/ac-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/epp/sec-[{epgPKey}]/ac-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/ldev-[{priKey}]-ctx-[{ctxDn}]-bd-[{bdDn}]/rsconsIf-{tnVzCPIfName}`
  - `uni/svcdefcont/tensvcdef-{name}/ldevdef-[{lDevDn}]/sepg-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/LDevInst-[{priKey}]-ctx-{ctxName}/G-{graphRn}-N-{nodeRn}-C-{connRn}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/LDevInst-[{priKey}]-ctx-{ctxName}/bd-[{bdDn}]/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/LDevInst-[{priKey}]-ctx-{ctxName}/epgDn-[{shEpgDn}]/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/Tnlepg-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acAnyToEp-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acEpToAny-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acEpToEp-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acEpToExt-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acExtToEp-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/acIpToIp-{name}/acl3outcont/outdef-{name}/instPdef-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/ap-{name}/epg-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/ap-{name}/esg-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/cloudapp-{name}/cloudepg-{name}/rsconsIf-{tnVzCPIfName}`
  - `uni/tn-{name}/cloudapp-{name}/cloudextepg-{name}/rsconsIf-{tnVzCPIfName}`
  - `Too many DN formats to display, see model documentation for all possible parents.`

## GUI Information ##

* `Location` - `Tenants -> ... determine for muliple DN formats ...`

## Example Usage ##

```hcl

data "aci_contract_interface" "example" {
  parent_dn               = aci_application_epg.example.id
  contract_interface_name = "test_tn_vz_cp_if_name"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_application_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/application_epg) ([fvAEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/fvAEPg/overview))
  - The distinquised name (DN) of classes below can be used but currently there is no available resource for it:
    - [cloudEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/cloudEPg/overview)
    - [cloudExtEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/cloudExtEPg/overview)
    - [cloudISvcEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/cloudISvcEPg/overview)
    - [cloudSvcEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/cloudSvcEPg/overview)
    - [dhcpCRelPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/dhcpCRelPg/overview)
    - [dhcpPRelPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/dhcpPRelPg/overview)
    - [fvESg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/fvESg/overview)
    - [fvTnlEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/fvTnlEPg/overview)
    - [infraCEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/infraCEPg/overview)
    - [infraPEPg](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/infraPEPg/overview)
    - [l2extInstP](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/l2extInstP/overview)
    - [l3extInstP](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/l3extInstP/overview)
    - [l3extInstPDef](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/l3extInstPDef/overview)
    - [mgmtInB](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/mgmtInB/overview)
    - [vnsEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/vnsEPpInfo/overview)
    - [vnsREPpInfo](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/vnsREPpInfo/overview)
    - [vnsSDEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/vnsSDEPpInfo/overview)
    - [vnsSHEPpInfo](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/vnsSHEPpInfo/overview)

* `contract_interface_name` - (string) The contract interface name.

### Read-Only

* `id` - (string) The distinguished name (DN) of the Contract Interface object.
* `annotation` - (string) The annotation of the Contract Interface object.
* `priority` - (string) The contract interface priority.

* `annotations` - (list) A list of Annotations objects ([tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/tagAnnotation/overview)).
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.