---
layout: "aci"
page_title: "ACI: aci_l2_outside"
sidebar_current: "docs-aci-resource-l2_outside"
description: |-
  Manages ACI L2 Outside
---

# aci_l2_outside #
Manages ACI L2 Outside

## Example Usage ##

```hcl
resource "aci_l2_outside" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  target_dscp  = "example"
}
```


## Argument Reference ##
* `tenant_dn` - (Required) distinguished name of parent Tenant object.
* `name` - (Required) name of Object l2_outside.
* `annotation` - (Optional) annotation for object l2_outside.
* `name_alias` - (Optional) name_alias for object l2_outside.
* `target_dscp` - (Optional) target dscp

* `relation_l2ext_rs_e_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
                
* `relation_l2ext_rs_l2_dom_att` - (Optional) Relation to class l2extDomP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L2 Outside.

## Importing ##

An existing L2 Outside can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l2_outside.example <Dn>
```