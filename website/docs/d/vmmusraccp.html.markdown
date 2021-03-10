---
layout: "aci"
page_title: "ACI: aci_vmm_credential"
sidebar_current: "docs-aci-data-source-vmm_credential"
description: |-
  Data source for ACI VMM Credential
---

# aci_vmm_credential #
Data source for ACI VMM Credential

## Example Usage ##

```hcl
data "aci_vmm_credential" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `name` - (Required) name of Object vmm_credential.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VMM Credential.
* `annotation` - (Optional) annotation for object vmm_credential.
* `name_alias` - (Optional) name_alias for object vmm_credential.
* `pwd` - (Optional) user account profile password
* `usr` - (Optional) user name
