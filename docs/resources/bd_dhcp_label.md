---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_bd_dhcp_label"
sidebar_current: "docs-aci-resource-aci_bd_dhcp_label"
description: |-
  Manages ACI BD DHCP Label
---

# aci_bd_dhcp_label

Manages ACI BD DHCP Label

## Example Usage

```hcl
resource "aci_bd_dhcp_label" "foo_bd_dhcp_label" {
  bridge_domain_dn  = aci_bridge_domain.foo_bridge_domain.id
  name  = "example"
  annotation  = "example"
  description = "from terraform"
  name_alias  = "example"
  owner  = "tenant"
  tag  = "aqua"
}
```

## Argument Reference

- `bridge_domain_dn` - (Required) Distinguished name of parent Bridge Domain object.
- `name` - (Required) The Bridge Domain DHCP label name. This name can be up to 64 alphanumeric characters.
- `annotation` - (Optional) Annotation for object BD DHCP Label.
- `description` - (Optional) Description for object BD DHCP Label.
- `name_alias` - (Optional) Name alias for object BD DHCP Label.
- `owner` - (Optional) Owner of the target relay servers.  
  Allowed values: "infra", "tenant". Default value: "infra".
- `tag` - (Optional) Label color.

- `relation_dhcp_rs_dhcp_option_pol` - (Optional) Relation to class dhcpOptionPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BD DHCP Label.

## Importing

An existing BD DHCP Label can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bd_dhcp_label.example <Dn>
```
