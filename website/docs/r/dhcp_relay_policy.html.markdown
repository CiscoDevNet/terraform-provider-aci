---
layout: "aci"
page_title: "ACI: aci_dhcp_relay_policy"
sidebar_current: "docs-aci-resource-dhcp_relay_policy"
description: |-
  Manages ACI DHCP Relay Policy
---

# aci_dhcp_relay_policy #
Manages ACI DHCP Relay Policy.

## Example Usage ##

```hcl
resource "aci_dhcp_relay_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

  name  = "name_example"
  annotation  = "annotation_example"
  mode  = "visible"
  name_alias  = "alias_example"
  owner  = "infra"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object dhcp_relay_policy.
* `annotation` - (Optional) Annotation for object dhcp_relay_policy.
* `mode` - (Optional) dhcp relay policy mode. Allowed Values are "visible" and "not-visible". Default Value is "visible".
* `name_alias` - (Optional) name_alias for object dhcp_relay_policy.
* `owner` - (Optional) Owner of the target relay servers. Allowed values are "infra" and "tenant". Default value is "infra".

* `relation_dhcp_rs_prov` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the DHCP Relay Policy.

## Importing ##

An existing DHCP Relay Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_dhcp_relay_policy.example <Dn>
```