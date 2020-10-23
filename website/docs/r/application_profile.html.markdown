---
layout: "aci"
page_title: "ACI: aci_application_profile"
sidebar_current: "docs-aci-resource-application_profile"
description: |-
  Manages ACI Application Profile
---

# aci_application_profile #
Manages ACI Application Profile

## Example Usage ##

```hcl
resource "aci_application_profile" "test_ap" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "demo_ap"
  annotation = "tag1,tag2"
  name_alias = "test_ap"
  prio       = "level1"
}

```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object application_profile.
* `annotation` - (Optional) annotation for object application_profile.
* `name_alias` - (Optional) name_alias for object application_profile.
* `prio` - (Optional) priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.

* `relation_fv_rs_ap_mon_pol` - (Optional) Relation to class monEPGPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Application Profile.

## Importing ##

An existing Application Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_application_profile.example <Dn>
```