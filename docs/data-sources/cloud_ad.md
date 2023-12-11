---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_ad"
sidebar_current: "docs-aci-data-source-aci_cloud_ad"
description: |-
  Data source for Cloud Network Controller Cloud Active Directory
---

# aci_cloud_ad #

Data source for Cloud Network Controller Cloud Active Directory
Note: This data source is supported in Cloud Network Controller only.

## API Information ##

* `Class` - cloudAD
* `Distinguished Name` - uni/tn-{tenant_name}/ad-{id}

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> Tenants  -> {tenant_name}



## Example Usage ##

```hcl
data "aci_cloud_ad" "example" {
  tenant_dn  = aci_tenant.example.id
  active_directory_id  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Active Directory object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Active Directory.
* `annotation` - (Optional) Annotation of the Active Directory object.
* `name_alias` - (Optional) Name Alias of the Active Directory object.
* `active_directory_id` - (Optional) Id of the Azure Active Directory.
