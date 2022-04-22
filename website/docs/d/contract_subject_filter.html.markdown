---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract_subject_filter"
sidebar_current: "docs-aci-data-source-contract_subject_filter"
description: |-
  Data source for ACI Contract Subject Filter
---

# aci_contract_subject_filter #

Data source for ACI Contract Subject Filter


## API Information ##

* `Class` - vzRsSubjFiltAtt
* `Distinguished Name` - uni/tn-{name}/brc-{name}/subj-{name}/rssubjFiltAtt-{tnVzFilterName}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Contracts -> Standard -> {contract_name} -> Contract Subject -> Filters



## Example Usage ##

```hcl
data "aci_contract_subject_filter" "example" {
  contract_subject_dn  = aci_contract_subject.example.id
  tn_vz_filter_name  = "example"
}
```

## Argument Reference ##

* `contract_subject_dn` - (Required) Distinguished name of parent Contract Subject object.
* `tn_vz_filter_name` - (Required) Name of Contract Subject Filter object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Subject Filter.
* `annotation` - (Optional) Annotation of the Contract Subject Filter object.
* `action` - (Optional) The action required when the condition is met.
* `directives` - (Optional) Directives of the Contract Subject Filter object. 
* `priority_override` - (Optional) Priority Override of the Contract Subject Filter object. 
