---
layout: "aci"
page_title: "ACI: aci_vmm_credential"
sidebar_current: "docs-aci-resource-vmm_credential"
description: |-
  Manages ACI VMM Credential
---

# aci_vmm_credential #
Manages ACI VMM Credential

## Example Usage ##

```hcl
resource "aci_vmm_credential" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  pwd  = "example"
  usr  = "example"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `name` - (Required) name of Object vmm_credential.
* `annotation` - (Optional) annotation for object vmm_credential.
* `name_alias` - (Optional) name_alias for object vmm_credential.
* `pwd` - (Optional) user account profile password
* `usr` - (Optional) user name
