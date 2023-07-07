---
subcategory: - "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_pim_interface_policy"
sidebar_current: "docs-aci-resource-pim_interface_policy"
description: |-
  Manages ACI PIM Interface Policy
---

# aci_piminterface_policy #

Manages ACI PIM Interface Policy

## API Information ##

* `Class` - pimIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/pimifpol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> PIM

## Example Usage ##

```hcl
resource "aci_pim_interface_policy" "example_ip" {
  tenant_dn   = aci_tenant.example.id
  name        = "examplei"
  auth_t      = "none"
  dr_delay    = "3"
  dr_prio     = "1"
  hello_itvl  = "30000"
  jp_interval = "60"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the PIM Interface Policy object.
* `annotation` - (Optional) Annotation of the PIM Interface Policy object.
* `name_alias` - (Optional) Name Alias of the PIM Interface Policy object.
* `auth_key` - (Optional) Authentication Key.
* `auth_t` - (Optional) Authentication type. Allowed values are "ah-md5", "none" and the default value is "none". Type: String.
* `ctrl` - (Optional) Interface controls. Allowed values are "border", "passive" and "strict-rfc-compliant". Type: List.
* `dr_delay` - (Optional) Designated router delay. Allowed range is "1-65535" and the default value is "3".
* `dr_prio` - (Optional) Designated router priority. Allowed range is "1-4294967295" and the default value is "1".
* `hello_itvl` - (Optional) Hello traffic policy. Allowed range is "1-18724286" and the default value is "30000".
* `jp_interval` - (Optional) JP Traffic Policy. Allowed range is "60-65520" and the default value is "60".
* `secure_auth_key` - (Optional) Secure Authentication key.

## Importing ##

An existing PIMInterfacePolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_pim_interface_policy.example <Dn>
```