---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_ldap_provider"
sidebar_current: "docs-aci-data-source-ldap_provider"
description: |-
  Data source for ACI LDAP Provider
---

# aci_ldap_provider #
Data source for ACI LDAP Provider


## API Information ##
* `Class` - aaaLdapProvider
* `Distinguished Name` - uni/userext/ldapext/ldapprovider-{name} & uni/userext/duoext/ldapprovider-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> LDAP -> Providers & Admin -> AAA -> Authentication -> DUO -> LDAP -> Providers

## Example Usage ##
```hcl
data "aci_ldap_provider" "example" {
  name  = "example"
  type = "duo"
}
```

## Argument Reference ##
* `name` - (Required) Host name or IP address of object LDAP Provider.
* `type` - (Required) Type of LDAP Provider. Allowed values are "ldap" and "duo".

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the LDAP Provider.
* `annotation` - (Optional) Annotation of object LDAP Provider.
* `description` - (Optional) Description of object LDAP Provider.
* `name_alias` - (Optional) Name alias of object LDAP Provider.
* `ssl_validation_level` - (Optional) The LDAP Server SSL Certificate validation level. 
* `attribute` - (Optional) The attribute to be downloaded that contains user role and domain information.
* `basedn` - (Optional) The LDAP base DN to be used in a user search. 
* `enable_ssl` - (Optional) A property for enabling an SSL connection with the LDAP provider. 
* `filter` - (Optional) The LDAP filter to be used in a user search. 
* `monitor_server` - (Optional) Periodic Server Monitoring. 
* `monitoring_user` - (Optional) Periodic Server Monitoring Username
* `port` - (Optional) The service port number for the LDAP service. 
* `retries` - (Optional) Retry count of object LDAP Provider. 
* `rootdn` - (Optional) The root DN or bind DN of the LDAP provider.
* `timeout` - (Optional) The timeout for communication with an LDAP provider server. 

