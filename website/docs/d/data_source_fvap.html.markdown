---
layout: "aci"
page_title: "ACI: aci_application_profile"
sidebar_current: "docs-aci-data-source-application_profile"
description: |-
  Data source for ACI Application Profile
---

# aci_application_profile #
Data source for ACI Application Profile

## Example Usage ##

```hcl
data "aci_application_profile" "dev_apps" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "foo_app"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object application_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Application Profile.
* `annotation` - (Optional) annotation for object application_profile.
* `name_alias` - (Optional) name_alias for object application_profile.
* `prio` - (Optional) priority class id
