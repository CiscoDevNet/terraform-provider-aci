---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_ldap_provider"
sidebar_current: "docs-aci-resource-ldap_provider"
description: |-
  Manages ACI LDAP Provider
---

# aci_ldap_provider #
Manages ACI LDAP Provider

## API Information ##
* `Class` - aaaLdapProvider
* `Distinguished Named` - uni/userext/ldapext/ldapprovider-{name} & uni/userext/duoext/ldapprovider-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> LDAP -> Providers & Admin -> AAA -> Authentication -> DUO -> LDAP -> Providers


## Example Usage ##
```hcl
resource "aci_ldap_provider" "example" {
	name = "example"
	type = "duo"
	description = "from terraform"
	annotation = "example_annotation"
	name_alias = "example_name_alias"
	ssl_validation_level = "strict"
	attribute = "CiscoAvPair"
	basedn = "CN=Users,DC=host,DC=com"
	enable_ssl = "yes"
	filter = "sAMAccountName=$userid"
	key = "example_key_value"
	monitor_server = "enabled"
	monitoring_password = "example_monitoring_password"
	monitoring_user = "example_monitoring_user_value"
	port = "389"
	retries = "1"
	rootdn = "CN=admin,CN=Users,DC=host,DC=com"
	timeout = "30"
}
```

## Argument Reference ##
* `name` - (Required) Host name or IP address of object LDAP Provider.
* `type` - (Required) Type of LDAP Provider. Allowed values are "ldap" and "duo".
* `annotation` - (Optional) Annotation of object LDAP Provider.
* `description` - (Optional) Description of object LDAP Provider.
* `name_alias` - (Optional) Name alias of object LDAP Provider.
* `ssl_validation_level` - (Optional) The LDAP Server SSL Certificate validation level. Allowed values are "permissive" and "strict". Default value is "strict". Type: String.
* `attribute` - (Optional) The attribute to be downloaded that contains user role and domain information. Default value is "CiscoAVPair".
* `basedn` - (Optional) The LDAP base DN to be used in a user search. 
* `enable_ssl` - (Optional) A property for enabling an SSL connection with the LDAP provider. Allowed values are "no" and "yes". Default value is "no". Type: String.
* `filter` - (Optional) The LDAP filter to be used in a user search. Default value is "sAMAccountName=$userid".
* `key` - (Optional) A password for the AAA provider database.
* `monitor_server` - (Optional) Periodic Server Monitoring. Allowed values are "disabled" and "enabled". Default value is "disabled". Type: String.
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. 
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. Default value is "default". 
* `port` - (Optional) The service port number for the LDAP service. Allowed range: "1" - "65535". Default value is "389".
* `retries` - (Optional) Retry count of object LDAP Provider. Allowed range: "1" - "5". Default value is "1".
* `rootdn` - (Optional) The root DN or bind DN of the LDAP provider.
* `timeout` - (Optional) The timeout for communication with an LDAP provider server. Allowed range: "5" - "60". Default value is "30". (NOTE: For "duo" LDAP providers, the value of timeout should be greater than or equal to "30".)
* `relation_aaa_rs_prov_to_epp` - (Optional) Represents the relation to a Relation to AProvider Reachability EPP (class fvAREpP). Type: String.
* `relation_aaa_rs_sec_prov_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the provider server is reachable. Type: String.

## Importing ##
An existing LDAP Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ldap_provider.example <Dn>
```