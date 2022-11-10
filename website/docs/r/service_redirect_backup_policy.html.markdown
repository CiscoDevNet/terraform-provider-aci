---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_service_redirect_backup_policy"
sidebar_current: "docs-aci-resource-service_redirect_backup_policy"
description: |-
  Manages ACI PBR Backup Policy
---

# aci_service_redirect_backup_policy #

Manages ACI PBR Backup Policy

## API Information ##

* `Class` - vnsBackupPol
* `Distinguished Name` - uni/tn-{tenant_name}/svcCont/backupPol-{backup_pol_name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> L4-L7 Policy-Based Redirect Backup


## Example Usage ##

```hcl
resource "aci_service_redirect_backup_policy" "pbr_backup_policy" {
  tenant_dn = aci_tenant.example.id
  name      = "pbr_backup_policy"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the PBR Backup Policy object.
* `annotation` - (Optional) Annotation of the PBR Backup Policy object.
* `name_alias` - (Optional) Name Alias of the PBR Backup Policy object.


## Importing ##

An existing PBRBackupPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_service_redirect_backup_policy.example <Dn>
```