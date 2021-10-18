---
layout: "aci"
page_title: "ACI: aci_saml_provider"
sidebar_current: "docs-aci-resource-saml_provider"
description: |-
  Manages ACI SAML Provider
---

# aci_saml_provider #

Manages ACI SAML Provider

## API Information ##

* `Class` - aaaSamlProvider
* `Distinguished Named` - uni/userext/samlext/samlprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> SAML -> Providers


## Example Usage ##

```hcl
resource "aci_saml_provider" "example" {
  name                      = "example"
  name_alias                = "saml_provider_alias"
  description               = "From Terraform"
  annotation                = "orchestrator:terraform"
  entity_id                 = "entity_id_example" 
  gui_banner_message        = "gui_banner_message_example"
  https_proxy               = "https_proxy_example"
  id_p                      = "adfs"
  key                       = "key_example"
  metadata_url              = "metadata_url_example"
  monitor_server            = "disabled"
  monitoring_password       = "monitoring_password_example"
  monitoring_user           = "default"
  retries                   = "1"
  sig_alg                   = "SIG_RSA_SHA256"
  timeout                   = "5"
  tp                        = "tp_example"
  want_assertions_encrypted = "yes"
  want_assertions_signed    = "yes"
  want_requests_signed      = "yes"
  want_response_signed      = "yes"
}
```

## Argument Reference ##


* `name` - (Required) Name of object SAML Provider.
* `annotation` - (Optional) Annotation of object SAML Provider.
* `name_alias` - (Optional) Name Alias of object SAML Provider. Type: String.
* `description` - (Optional) Description of object SAML Provider. Type: String.
* `entity_id` - (Optional) Entity ID. Type: String.
* `gui_banner_message` - (Optional) Gui Redirect Banner Message. Type: String.
* `https_proxy` - (Optional) Https Proxy to reach IDP's Metadata URL. Type: String. (Note: Value passed for "https_proxy" attribute should be a URL)
* `id_p` - (Optional) Identity Provider. Allowed values are "adfs", "okta", "ping identity". Default value is "adfs". Type: String.
* `key` - (Optional) Key. A password for the AAA provider database. Type: String.
* `metadata_url` - (Optional) Metadata Url provided by IDP. Type: String.
* `monitor_server` - (Optional) Periodic Server Monitoring. Allowed values are "disabled", "enabled". Default value is "disabled". Type: String.
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. Type: String.
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. Default value is "default". Type: String.
* `retries` - (Optional) Retries. Allowed range is "1"-"5" and default value is "1". Type: String.
* `sig_alg` - (Optional) Signature Algorithm. Allowed values are "SIG_RSA_SHA1", "SIG_RSA_SHA224", "SIG_RSA_SHA256", "SIG_RSA_SHA384", "SIG_RSA_SHA512". Default value is "SIG_RSA_SHA256". Type: String.
* `timeout` - (Optional) Timeout in Seconds. The amount of time between authentication attempts. Allowed range is "5"-"60" and default value is "5". Type: String.
* `tp` - (Optional) Certificate Authority. Type: String.
* `want_assertions_encrypted` - (Optional) Want Encrypted SAML Assertions. Allowed values are "no" and "yes". Default value is "yes". Type: String.
* `want_assertions_signed` - (Optional) Want Assertions in SAML Response Signed. Allowed values are "no" and "yes". Default value is "yes". Type: String.
* `want_requests_signed` - (Optional) Want SAML Auth Requests Signed. Allowed values are "no" and "yes". Default value is "yes". Type: String.
* `want_response_signed` - (Optional) Want SAML Response Message Signed. Allowed values are "no" and "yes". Default value is "yes". Type: String.
* `relation_aaa_rs_prov_to_epp` - (Optional) Represents the relation to a Relation to AProvider Reachability EPP (class fvAREpP). Type: String.
* `relation_aaa_rs_sec_prov_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the provider server is reachable. Type: String.



## Importing ##

An existing SAMLProvider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_saml_provider.example <Dn>
```