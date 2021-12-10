---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_endpoint_ip_aging_profile"
sidebar_current: "docs-aci-data-source-endpoint_ip_aging_profile"
description: |-
  Data source for ACI Endpoint IP Aging Profile
---

# aci_endpoint_ip_aging_profile #
Data source for ACI IP Aging Profile


## API Information ##
* `Class` - epIpAgingP
* `Distinguished Named` - uni/infra/ipAgingP-{name}

## GUI Information ##
* `Location` - System -> System Settings -> Endpoint Controls -> IP Aging -> Policy

## Example Usage ##
```hcl
data "aci_endpoint_ip_aging_profile" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint IP Aging Profile.
* `annotation` - (Optional) Annotation of object Endpoint IP Aging Profile.
* `name_alias` - (Optional) Name Alias of object Endpoint IP Aging Profile.
* `admin_st` - (Optional) The administrative state of the object Endpoint IP Aging Profile.
* `description` - (Optional) Description of object Endpoint IP Aging Profile.
