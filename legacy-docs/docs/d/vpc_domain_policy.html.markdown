---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_vpc_domain_policy"
sidebar_current: "docs-aci-data-source-vpc_domain_policy"
description: |-
  Data source for ACI VPC Domain Policy
---

# aci_vpc_domain_policy #
Data source for ACI VPC Domain Policy


## API Information ##
* `Class` - vpcInstPol
* `Distinguished Name` - uni/fabric/vpcInst-{name}

## GUI Information ##
* `Location` - Fabric -> Access Policies -> Policies -> Switch -> VPC Domain -> Create VPC Domain Policy

## Example Usage ##

```hcl
data "aci_vpc_domain_policy" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object VPC Domain Policy.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the VPC Domain Policy.
* `annotation` - (Optional) Annotation of object VPC Domain Policy.
* `name_alias` - (Optional) Name Alias of object VPC Domain Policy.
* `dead_intvl` - (Optional) The VPC peer dead interval time of object VPC Domain Policy
* `description` - (Optional) Description of object VPC Domain Policy.
