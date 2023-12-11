---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_fc_domain"
sidebar_current: "docs-aci-resource-aci_fc_domain"
description: |-
  Manages ACI FC Domain
---

# aci_fc_domain

Manages ACI FC Domain

## Example Usage

```hcl
resource "aci_fc_domain" "example" {
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of Object fibre channel domain.
- `annotation` - (Optional) Annotation for object fibre channel domain.
- `name_alias` - (Optional) Name alias for object fibre channel domain.

- `relation_infra_rs_vlan_ns` - (Optional) Relation to class fvnsVlanInstP. Cardinality - N_TO_ONE. Type - String.
- `relation_fc_rs_vsan_ns` - (Optional) Relation to class fvnsVsanInstP. Cardinality - N_TO_ONE. Type - String.
- `relation_fc_rs_vsan_attr` - (Optional) Relation to class fcVsanAttrP. Cardinality - N_TO_ONE. Type - String.
- `relation_infra_rs_vlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
- `relation_infra_rs_vip_addr_ns` - (Optional) Relation to class fvnsAddrInst. Cardinality - N_TO_ONE. Type - String.
- `relation_infra_rs_dom_vxlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
- `relation_fc_rs_vsan_attr_def` - (Optional) Relation to class fcVsanAttrP. Cardinality - N_TO_ONE. Type - String.
- `relation_fc_rs_vsan_ns_def` - (Optional) Relation to class fvnsAVsanInstP. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the FC Domain.

## Importing

An existing FC Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_fc_domain.example <Dn>
```
