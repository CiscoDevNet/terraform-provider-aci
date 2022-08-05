---
subcategory: "Virtual Networking"
layout: "aci"
page_title: "ACI: aci_vmm_credential"
sidebar_current: "docs-aci-resource-vmm_credential"
description: |-
  Manages ACI VMM Credential
---

# aci_vmm_credential #

Manages ACI VMM Credential

## API Information ##

* `Class` - vmmUsrAccP
* `Distinguished Name` - uni/vmmp-{vendor}/dom-{name}/usracc-{name}

## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> vCenter Credentials
  
## Example Usage ##

```hcl
resource "aci_vmm_credential" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  name  = "vmm_credential_1"
  annotation = "vmm_credential_tag"
  description = "from terraform"
  name_alias = "vmm_credential_alias"
  pwd = "password"
  usr = "username"
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
* `name` - (Required) Name of object VMM Credential.
* `annotation` - (Optional) Annotation of object VMM Credential.
* `description` - (Optional) Description of object VMM Credential.
* `name_alias` - (Optional) Name alias of object VMM Credential.
* `pwd` - (Optional) Password.
* `usr` - (Optional) Username. Min length is "0". Max length is "128". If any value is assigned to a username then it cannot be updated to an empty value. 

## Importing ##

An existing VMMCredential can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html
  
```
terraform import vmm_credential.example <Dn>
```