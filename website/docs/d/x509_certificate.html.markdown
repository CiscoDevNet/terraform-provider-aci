---
layout: "aci"
page_title: "ACI: aci_x509_certificate"
sidebar_current: "docs-aci-data-source-x509_certificate"
description: |-
  Data source for ACI X509 Certificate
---

# aci_x509_certificate

Data source for ACI X509 Certificate

## Example Usage

```hcl
data "aci_x509_certificate" "example" {
  local_user_dn  = aci_local_user.example.id
  name  = "x509_certificate_1"
}
```

## Argument Reference

- `local_user_dn` - (Required) Distinguished name of parent LocalUser object.
- `name` - (Required) Name of Object x509 certificate.

## Attribute Reference

- `id` - Attribute id set to the Dn of the X509 Certificate.
- `description` - (Optional) Description for object x509 certificate.
- `annotation` - (Optional) Annotation for object x509 certificate.
- `data` - (Optional) Data from the user certificate
- `name_alias` - (Optional) Name alias for object x509 certificate.
