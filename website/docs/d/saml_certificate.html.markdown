---
layout: "aci"
page_title: "ACI: aci_saml_certificate"
sidebar_current: "docs-aci-data-source-saml_certificate"
description: |-
  Data source for ACI SAML Encryption Certificate
---

# aci_saml_certificate #
Data source for ACI SAML Encryption Certificate


## API Information ##
* `Class` - aaaSamlEncCert
* `Distinguished Named` - uni/userext/samlext/samlenccert-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> Management -> Public Key for SAML Encryption

## Example Usage ##
```hcl
data "aci_saml_certificate" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Key pair for SAML Encryption Certificate.
* `annotation` - (Optional) Annotation of object SAML Encryption Certificate.
* `name_alias` - (Optional) Name Alias of object SAML Encryption Certificate.
* `description` - (Optional) Description of object SAML Encryption Certificate.
* `regenerate` - (Optional) Regenerate Encryption Key Pair. 
* `certificate` - (Optional) Certificate of SAML Encryption Key.
* `certificate_validty` - (Optional) Certificate validity of SAML Encryption Certificate.
* `expiry_status` - (Optional) Expiry status of SAML Encryption Certificate.

