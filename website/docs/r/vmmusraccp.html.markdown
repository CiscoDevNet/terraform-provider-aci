---
layout: "aci"
page_title: "ACI: aci_vmm_credential"
sidebar_current: "docs-aci-resource-vmm_credential"
description: |-
  Manages ACI VMM Credential
---

# vmm_credential #

Manages ACI VMM Credential

## API Information ##

* `Class` - vmmUsrAccP
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/usracc-{name}

## GUI Information ##

* `Location` - Virtual Networking -> VMM Domain -> VmmController -> vCenterCredentials


## Example Usage ##

```hcl
resource "aci_vmm_credential" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  name  = "example"
  annotation = "orchestrator:terraform"

  pwd = 
  usr = 
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `name` - (Required) Name of object VMM Credential.
* `annotation` - (Optional) Annotation of object VMM Credential.

* `pwd` - (Optional) Password. Pwd 
* `usr` - (Optional) Username. User 


## Importing ##

An existing VMMCredential can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import vmm_credential.example <Dn>
```