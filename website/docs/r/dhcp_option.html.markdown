---
layout: "aci"
page_title: "ACI: aci_dhcp_option"
sidebar_current: "docs-aci-resource-dhcp_option"
description: |-
  Manages ACI DHCP Option
---

# aci_dhcp_option

Manages ACI DHCP Option

## Example Usage

```hcl
resource "aci_dhcp_option" "example" {

  dhcp_option_policy_dn  = "${aci_dhcp_option_policy.example.id}"

  name  = "example"
  annotation  = "example"
  data  = "example"
  dhcp_option_id  = "1"
  name_alias  = "example"
}
```

## Argument Reference

- `dhcp_option_policy_dn` - (Required) Distinguished name of parent DHCPOptionPolicy object.
- `name` - (Required) Name of Object dhcp_option.
- `annotation` - (Optional) Annotation for object dhcp_option.
- `data` - (Optional) DHCP option data
- `dhcp_option_id` - (Optional) DHCP option id (Unsigned Integer)
- `name_alias` - (Optional) name_alias for object dhcp_option.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the DHCP Option.

## Importing

An existing DHCP Option can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_dhcp_option.example <Dn>
```
