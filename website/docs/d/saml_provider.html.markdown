---
layout: "aci"
page_title: "ACI: aci_saml_provider"
sidebar_current: "docs-aci-data-source-saml_provider"
description: |-
  Data source for ACI SAML Provider
---

# aci_saml_provider #

Data source for ACI SAML Provider


## API Information ##

* `Class` - aaaSamlProvider
* `Distinguished Named` - uni/userext/samlext/samlprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> SAML -> Providers



## Example Usage ##

```hcl
data "aci_saml_provider" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object SAML Provider.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the SAML Provider.
* `annotation` - (Optional) Annotation of object SAML Provider.
* `name_alias` - (Optional) Name Alias of object SAML Provider.
* `description` - (Optional) Description of object SAML Provider.
* `entity_id` - (Optional) Entity ID. 
* `gui_banner_message` - (Optional) Gui Redirect Banner Message. 
* `https_proxy` - (Optional) Https Proxy to reach IDP's Metadata URL. 
* `id_p` - (Optional) Identity Provider. 
* `key` - (Optional) Key. A password for the AAA provider database.
* `metadata_url` - (Optional) Metadata Url provided by IDP. 
* `monitor_server` - (Optional) Periodic Server Monitoring. 
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. 
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. 
* `retries` - (Optional) Retries. null
* `sig_alg` - (Optional) Signature Algorithm. 
* `timeout` - (Optional) Timeout in Seconds. The amount of time between authentication attempts.
* `tp` - (Optional) Certificate Authority. 
* `want_assertions_encrypted` - (Optional) Want Encrypted SAML Assertions. 
* `want_assertions_signed` - (Optional) Want Assertions in SAML Response Signed. 
* `want_requests_signed` - (Optional) Want SAML Auth Requests Signed. 
* `want_response_signed` - (Optional) Want SAML Response Message Signed. 
