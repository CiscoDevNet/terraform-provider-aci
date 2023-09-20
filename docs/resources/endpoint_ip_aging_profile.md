---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_endpoint_ip_aging_profile"
sidebar_current: "docs-aci-resource-endpoint_ip_aging_profile"
description: |-
  Manages ACI Endpoint IP Aging Profile
---

# aci_endpoint_ip_aging_profile #
Manages ACI Endpoint IP Aging Profile

## API Information ##
* `Class` - epIpAgingP
* `Distinguished Name` - uni/infra/ipAgingP-{name}

## GUI Information ##
* `Location` - System -> System Settings -> Endpoint Controls -> IP Aging -> Policy


## Example Usage ##

```hcl
resource "aci_endpoint_ip_aging_profile" "example" {
  admin_st = "disabled"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
}
```

## NOTE ##
Users can use the resource of type `aci_endpoint_ip_aging_profile` to change the configuration of the object Endpoint IP Aging Profile. Users cannot create more than one instance of object Endpoint IP Aging Profile.

## Argument Reference ##
* `admin_st` - (Optional) The administrative state of the object Endpoint IP Aging Profile. Allowed values are "disabled" and "enabled".
* `annotation` - (Optional) Annotation of object Endpoint IP Aging Profile.
* `description` - (Optional) Description of object Endpoint IP Aging Profile.
* `name_alias` - (Optional) Name Alias of object Endpoint IP Aging Profile.

## Importing ##
An existing Endpoint IP Aging Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_ip_aging_profile.example <Dn>
```

