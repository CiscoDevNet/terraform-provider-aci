---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject_one_way_filter"
sidebar_current: "docs-aci-contract-subject-one-way-filter"
description: |-
  Manages ACI One Way Filter
---

# aci_contract_subject_one_way_filter #

Manages ACI One Way Filter

NOTE: The aci_contract_subject_one_way_filter should only be used with datasources as it will create conflicts with the aci_contract_subject resource

## API Information ##

* `Class` - vzRsFiltAtt
* ` Supported Distinguished Name` - <br>
[1] uni/tn-{tenant_name}/brc-{contract_name}/subj-{contract_subject_name}/intmnl/rsfiltAtt-{tnVzFilterName}<br>
[2] uni/tn-{tenant_name}/brc-{contract_name}/subj-{contract_subject_name}/outmnl/rsfiltAtt-{tnVzFilterName}}<br>

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Contracts -> Standard -> {contract_name} -> {subject_name} -> Filter Chain For Consumer to Provider / Filter Chain For Provider to Consumer


## Example Usage ##

```hcl
resource "aci_contract_subject_one_way_filter" "example" {
  contract_subject_dn  = aci_contract_subject.example.id
  filter_dn  = data.aci_filter.test_filter.id
  action = "permit"
  directives = ["none"]
  priority_override = "default"

}
```

## Argument Reference ##

* `contract_subject_dn` - (Required) Distinguished name of the parent Contract Subject object.
* `filter_dn` - (Required) Distinguished name of the Filter object.
* `annotation` - (Optional) Annotation of the Filter object.
* `action` - (Optional) The action required when the condition is met. Allowed values are "deny", "permit", and the default value is "permit". Type: String.
* `directives` - (Optional) Directives of the Contract Subject Filter object. Allowed values are "log", "no_stats", "none", and the default value is "none". Type: List.
* `priority_override` - (Optional) Priority override of the Filter object. It is only used when action is set to deny. Allowed values are "default", "level1", "level2", "level3", and the default value is "default". Type: String.


## Importing ##

An existing Filter can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract_subject_one_way_filter.example <Dn>
```