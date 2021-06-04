---
layout: "aci"
page_title: "ACI: aci_x509_certificate"
sidebar_current: "docs-aci-resource-x509_certificate"
description: |-
  Manages ACI X509 Certificate
---

# aci_x509_certificate

Manages ACI X509 Certificate

## Example Usage

```hcl
resource "aci_x509_certificate" "example" {
  local_user_dn  = aci_local_user.example.id
  name  = "x509_certificate_1"
  annotation  = "x509_certificate_tag"
  data  = "example"
  name_alias  = "alias_name"
}
```

## Argument Reference

- `local_user_dn` - (Required) Distinguished name of parent LocalUser object.
- `name` - (Required) Name of Object x509 certificate.
- `annotation` - (Optional) Annotation for object x509 certificate.
- `description` - (Optional) Description for object x509 certificate.
- `data` - (Optional) Data from the user certificate
- `name_alias` - (Optional) Name alias for object x509 certificate.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the X509 Certificate.

## Importing

An existing X509 Certificate can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_x509_certificate.example <Dn>
```
