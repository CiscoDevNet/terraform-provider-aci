---
layout: "aci"
page_title: "ACI: aci_dhcp_option_policy"
sidebar_current: "docs-aci-resource-dhcp_option_policy"
description: |-
  Manages ACI DHCP Option Policy
---

# aci_dhcp_option_policy

Manages ACI DHCP Option Policy

## Example Usage

```hcl
resource "aci_dhcp_option_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"

  dhcp_option {
      name  = "example_one"
      annotation  = "annotation_one"
      data  = "data_one"
      dhcp_option_id  = "1"
      name_alias  = "one"
  }
  dhcp_option {
      name  = "example_two"
      annotation  = "annotation_two"
      data  = "data_two"
      dhcp_option_id  = "2"
      name_alias  = "two"
  }

}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object DHCP Option Policy.
- `annotation` - (Optional) Annotation for object DHCP Option Policy.
- `name_alias` - (Optional) Name alias for object DHCP Option Policy.
- `dhcp_option` - (Optional) to manage DHCP Option from the DHCP Option Policy resource. It has the attributes like name, annotation,data,dhcp_option_id and name_alias.
- `dhcp_option.name` - (Required) Name of Object DHCP Option.
- `dhcp_option.annotation` - (Optional) Annotation for object DHCP Option.
- `dhcp_option.data` - (Optional) DHCP Option data.
- `dhcp_option.dhcp_option_id` - (Optional) DHCP Option id (Unsigned Integer).
- `dhcp_option.name_alias` - (Optional) Name alias for object DHCP Option.

## Attribute Reference

- `id` - Dn of the DHCP Option Policy.
- `dhcp_option.id` - exports this attribute for DHCP Option object. Set to the Dn for the DHCP Option managed by the DHCP Option policy.

## Importing

An existing DHCP Option Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_dhcp_option_policy.example <Dn>
```
