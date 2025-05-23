---
# Documentation generated by "gen/generator.go"; DO NOT EDIT.
# In order to regenerate this file execute `go generate` from the repository root.
# More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_certificate_authority"
sidebar_current: "docs-aci-resource-aci_certificate_authority"
description: |-
  Manages ACI Certificate Authority
---

# aci_certificate_authority #

Manages ACI Certificate Authority



## API Information ##

* Class: [pkiTP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/pkiTP/overview)

* Supported in ACI versions: 1.0(1e) and later.

* Distinguished Name Formats:
  - `uni/tn-{name}/certstore/tp-{name}`
  - `uni/userext/pkiext/tp-{name}`

## GUI Information ##

* Locations:
  - `Admin -> AAA -> Security -> Certificate Authorities`
  - `Cloud Network Controller -> Administrative -> Security -> Certificate Authorities`

## Example Usage ##

The configuration snippet below creates a Certificate Authority with only required attributes.

```hcl

resource "aci_certificate_authority" "example" {
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  name              = "test_name"
}

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_certificate_authority" "example_tenant" {
  parent_dn         = aci_tenant.example.id
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  name              = "test_name"
}

```
The configuration snippet below shows all possible attributes of the Certificate Authority.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl

resource "aci_certificate_authority" "full_example" {
  annotation        = "annotation"
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  certificate_usage = ["WebSvcOrAuth"]
  description       = "description_1"
  name              = "test_name"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
}

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_certificate_authority" "full_example_tenant" {
  parent_dn         = aci_tenant.example.id
  annotation        = "annotation"
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  certificate_usage = ["WebSvcOrAuth"]
  description       = "description_1"
  name              = "test_name"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
}

```

All examples for the Certificate Authority resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/tree/master/examples/resources/aci_certificate_authority) folder.

## Schema ##

### Required ###

* `certificate_chain` (certChain) - (string) The PEM-encoded chain of trust from the trustpoint to a trusted root authority.
* `name` (name) - (string) The name of the Certificate Authority object.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the Certificate Authority object.

### Optional ###

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/fvTenant) ([fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview))
  - Default: `uni/userext/pkiext`
* `annotation` (annotation) - (string) The annotation of the Certificate Authority object. This attribute is supported in ACI versions: 3.2(1l) and later.
  - Default: `orchestrator:terraform`
* `certificate_usage` (certUsage) - (list) The usage of the Certificate Authority object which can be Web Service/Authentication (WebSvcOrAuth), Platform Experience Grid (pxGrid), or both.
  - Default: `WebSvcOrAuth`.
  - Valid Values: `WebSvcOrAuth`, `pxGrid`.
* `description` (descr) - (string) The description of the Certificate Authority object.
* `name_alias` (nameAlias) - (string) The name alias of the Certificate Authority object. This attribute is supported in ACI versions: 2.2(1k) and later.
* `owner_key` (ownerKey) - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` (ownerTag) - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.
* `annotations` - (list) A list of Annotations (ACI object [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)). Annotations can also be configured using a separate [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  #### Required ####
  
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.
* `tags` - (list) A list of Tags (ACI object [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview)). Tags can also be configured using a separate [aci_tag](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tag) resource. This attribute is supported in ACI versions: 3.2(1l) and later.
  #### Required ####
  
    * `key` (key) - (string) The key used to uniquely identify this configuration object.
    * `value` (value) - (string) The value of the property.

## Importing

An existing Certificate Authority can be [imported](https://www.terraform.io/docs/import/index.html) into this resource with its distinguished name (DN), via the following command:

```
terraform import aci_certificate_authority.example uni/tn-{name}/certstore/tp-{name}
```

Starting in Terraform version 1.5, an existing Certificate Authority can be imported
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/certstore/tp-{name}"
  to = aci_certificate_authority.example
}
```
