---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject_filter"
sidebar_current: "docs-aci-resource-contract_subject_filter"
description: |-
  Manages ACI Contract Subject Filter
---

# aci_contract_subject_filter #

Manages ACI Contract Subject Filter

## API Information ##

* `Class` - vzRsSubjFiltAtt
* `Distinguished Name` - uni/tn-{name}/brc-{name}/subj-{name}/rssubjFiltAtt-{tnVzFilterName}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Contracts -> Standard -> {contract_name} -> Contract Subject -> Filters


## Example Usage ##

```hcl
resource "aci_contract_subject_filter" "example" {
  contract_subject_dn  = aci_contract_subject.example.id
  tnVzFilterName  = "example"
  action = "permit"
  directives = ["none"]
  priority_override = "default"

}
```

## Argument Reference ##

* `contract_subject_dn` - (Required) Distinguished name of the parent Contract Subject object.
* `tnVzFilterName` - (Required) TnVzFilterName of the Contract Subject Filter object.
* `annotation` - (Optional) Annotation of the Contract Subject Filter object.
* `action` - (Optional) The action required when the condition is met. Allowed values are "deny", "permit", and default value is "permit". Type: String.
* `directives` - (Optional) Directives of the Contract Subject Filter object. Allowed values are "log", "no_stats", "none", and default value is "none". Type: List.
* `priority_override` - (Optional) Priority override of the Contract Subject Filter object. Allowed values are "default", "level1", "level2", "level3", and default value is "default". Type: String.


## Importing ##

An existing Contract Subject Filter can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract_subject_filter.example <Dn>
```