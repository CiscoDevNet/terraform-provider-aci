---
# Documentation generated by "gen/generator.go"; DO NOT EDIT.
# In order to regenerate this file execute `go generate` from the repository root.
# More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_mld_snooping_policy"
sidebar_current: "docs-aci-data-source-aci_mld_snooping_policy"
description: |-
  Data source for ACI MLD Snooping Policy
---

# aci_mld_snooping_policy #

Data source for ACI MLD Snooping Policy

## API Information ##

* Class: [mldSnoopPol](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mldSnoopPol/overview)

* Supported in ACI versions: 4.1(1i) and later.

* Distinguished Name Formats:
  - `uni/fabric/mldsnoopPol-{name}`
  - `uni/tn-{name}/mldsnoopPol-{name}`

## GUI Information ##

* Location: `Tenants -> Policies -> Protocol -> MLD Snoop`

## Example Usage ##

```hcl

data "aci_mld_snooping_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}

```

## Schema ##

### Required ###

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tenant) ([fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview))
* `name` (name) - (string) The name of the MLD Snooping Policy object.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the MLD Snooping Policy object.
* `admin_state` (adminSt) - (string) The administrative state of the MLD Snooping Policy object.
* `annotation` (annotation) - (string) The annotation of the MLD Snooping Policy object.
* `control` (ctrl) - (list) The controls for the MLD Snooping Policy object.
* `description` (descr) - (string) The description of the MLD Snooping Policy object.
* `last_member_interval` (lastMbrIntvl) - (string) The last member interval (seconds) of the MLD Snooping Policy object. The group state is removed when no host responds before the timeout.
* `name_alias` (nameAlias) - (string) The name alias of the MLD Snooping Policy object.
* `owner_key` (ownerKey) - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` (ownerTag) - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `query_interval` (queryIntvl) - (string) The query interval (seconds) of the MLD Snooping Policy object.
* `response_interval` (rspIntvl) - (string) The response interval (seconds) of the MLD Snooping Policy object.
* `start_query_count` (startQueryCnt) - (string) The start query count of the MLD Snooping Policy object.
* `start_query_interval` (startQueryIntvl) - (string) The query interval (seconds) of the MLD Snooping Policy object at start-up.
* `version` (ver) - (string) The MLD version.
* `annotations` - (list) A list of Annotations (ACI object [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)). This attribute is supported in ACI versions: 3.2(1l) and later.
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.
* `tags` - (list) A list of Tags (ACI object [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview)). This attribute is supported in ACI versions: 3.2(1l) and later.
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.