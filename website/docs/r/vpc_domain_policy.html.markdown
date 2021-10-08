---
layout: "aci"
page_title: "ACI: aci_vpc_domain_policy"
sidebar_current: "docs-aci-resource-vpc_domain_policy"
description: |-
  Manages ACI VPC Domain Policy
---

# aci_vpc_domain_policy #
Manages ACI VPC Domain Policy

## API Information ##
* `Class` - vpcInstPol
* `Distinguished Named` - uni/fabric/vpcInst-{name}

## GUI Information ##
* `Location` - Fabric -> Access Policies -> Policies -> Switch -> VPC Domain -> Create VPC Domain Policy


## Example Usage ##

```hcl
resource "aci_vpc_domain_policy" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  dead_intvl = "200"
  name_alias = "example"
  description = "from terraform"
}
```

## Argument Reference ##
* `name` - (Required) Name of object VPC Domain Policy.
* `annotation` - (Optional) Annotation of object VPC Domain Policy.
* `dead_intvl` - (Optional) The VPC peer dead interval time of object VPC Domain Policy. Range: "5" - "600". Default value is "200".
* `name_alias` - (Optional) Name Alias of object VPC Domain Policy.
* `description` - (Optional) Description of object VPC Domain Policy.



## Importing ##
An existing VPC Domain Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vpc_domain_policy.example <Dn>
```