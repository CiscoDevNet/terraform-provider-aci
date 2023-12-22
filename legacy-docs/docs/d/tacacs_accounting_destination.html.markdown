---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_accounting_destination"
sidebar_current: "docs-aci-data-source-aci_tacacs_accounting_destination"
description: |-
  Data source for ACI TACACS Accounting Destination
---

# aci_tacacs_accounting_destination #
Data source for ACI TACACS Accounting Destination


## API Information ##
* `Class` - tacacsTacacsDest
* `Distinguished Name` - uni/fabric/tacacsgroup-{name}/tacacsdest-{host}-port-{port}

## GUI Information ##
* `Location` - Admin -> External Data Collectors -> Monitoring Destinations -> TACACS -> TACACS Destinations

## Example Usage ##
```hcl
data "aci_tacacs_accounting_destination" "example" {
  tacacs_accounting_dn = aci_tacacs_accounting.example.id
  host  = "cisco.com"
  port  = "49"
}
```

## Argument Reference ##
* `tacacs_accounting_dn` - (Required) Distinguished name of parent TACACS Accounting object.
* `host` - (Required) Host or IP address of object TACACS Accounting Destination.
* `port` - (Optional) Port of object TACACS Accounting Destination. Allowed Range: "1" - "65535". Default value: "49". 

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the TACACS Destination.
* `annotation` - (Optional) Annotation of object TACACS Accounting Destination.
* `name_alias` - (Optional) Name Alias of object TACACS Accounting Destination.
* `name` - (Optional) Name of object TACACS Accounting Destination.
* `auth_protocol` - (Optional) Authentication Protocol of object TACACS Accounting Destination. 
* `description` - (Optional) Description of object TACACS Accounting Destination.

