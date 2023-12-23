---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spanning_tree_interface_policy"
sidebar_current: "docs-aci-resource-aci_spanning_tree_interface_policy"
description: |-
  Manages ACI Spanning Tree Interface Policy
---

# aci_spanning_tree_interface_policy

Manages ACI Spanning Tree Interface Policy

## API Information

- `Class` - stpIfPol
- `Distinguished Name` - uni/infra/ifPol-{name}

## GUI Information

- `Location` - Fabric > Access Policies > Policies > Interface > Spanning Tree Interface

## Example Usage

```hcl
resource "aci_spanning_tree_interface_policy" "example" {
  name        = "spanning_tree_interface_policy"
  annotation  = "spanning_tree_interface_policy_tag"
  description = "from terraform"
  name_alias = "spanning_tree_interface_policy_alias"
  ctrl        = ["unspecified"]
}
```

## Argument Reference

- `name` - (Required) Name of object Spanning Tree Interface Policy.
- `annotation` - (Optional) Annotation of object Spanning Tree Interface Policy.
- `name_alias` - (Optional) Name alias of object Spanning Tree Interface Policy.
- `description` - (Optional) Description of object Spanning Tree Interface Policy.
- `ctrl` - (Optional) Interface controls. Allowed values are "bpdu-filter", "bpdu-guard", "unspecified". The default value is "unspecified". Unspecified value should not be added along with other allowed values. Type: List.

## Importing

An existing SpanningTreeInterfacePolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import spanning_tree_interface_policy.example <Dn>
```
