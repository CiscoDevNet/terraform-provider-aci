---
subcategory: "Virtual Networking"
layout: "aci"
page_title: "ACI: aci_vmm_credential"
sidebar_current: "docs-aci-data-source-vmm_credential"
description: |-
  Data source for ACI VMM Credential
---

# aci_vmm_credential #

Data source for ACI VMM Credential


## API Information ##

* `Class` - vmmUsrAccP
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/usracc-{name}

## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> vCenter Credentials

## Example Usage ##

```hcl
data "aci_vmm_credential" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  name  = "example"
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
* `name` - (Required) Name of object VMM Credential.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the VMM Credential.
* `annotation` - (Optional) Annotation of object VMM Credential.
* `description` - (Optional) Description of object VMM Credential.
* `name_alias` - (Optional) Name Alias of object VMM Credential.
* `pwd` - (Optional) Password.
* `usr` - (Optional) Username.
