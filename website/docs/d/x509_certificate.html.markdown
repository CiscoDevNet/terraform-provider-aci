---
layout: "aci"
page_title: "ACI: aci_x509_certificate"
sidebar_current: "docs-aci-data-source-x509_certificate"
description: |-
  Data source for ACI X509 Certificate
---

# aci_x509_certificate #
Data source for ACI X509 Certificate

## Example Usage ##

```hcl
data "aci_x509_certificate" "example" {

  local_user_dn  = "${aci_local_user.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `local_user_dn` - (Required) Distinguished name of parent LocalUser object.
* `name` - (Required) name of Object x509_certificate.



## Attribute Reference

* `id` - Attribute id set to the Dn of the X509 Certificate.
* `annotation` - (Optional) annotation for object x509_certificate.
* `data` - (Optional) data from the user certificate
* `name_alias` - (Optional) name_alias for object x509_certificate.
