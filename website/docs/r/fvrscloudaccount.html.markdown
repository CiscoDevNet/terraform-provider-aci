---
layout: "aci"
page_title: "ACI: aci_tenanttoaccountassociation"
sidebar_current: "docs-aci-resource-tenanttoaccountassociation"
description: |-
  Manages ACI Tenant to account association
---

# aci_tenanttoaccountassociation #

Manages ACI Tenant to account association

## API Information ##

* `Class` - fvRsCloudAccount
* `Distinguished Name` - uni/tn-{name}/rsCloudAccount

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_tenanttoaccountassociation" "example" {
  tenant_dn  = aci_tenant.example.id
  annotation = "orchestrator:terraform"
  t_dn = 
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.

* `annotation` - (Optional) Annotation of the object Tenant to account association.

* `t_dn` - (Optional) Target-dn.The distinguished name of the target.


## Importing ##

An existing Tenanttoaccountassociation can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tenanttoaccountassociation.example <Dn>
```