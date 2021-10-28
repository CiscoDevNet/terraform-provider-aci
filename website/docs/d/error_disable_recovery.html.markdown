---
layout: "aci"
page_title: "ACI: aci_error_disable_recovery"
sidebar_current: "docs-aci-data-source-error_disable_recovery"
description: |-
  Data source for ACI Error Disable Recovery
---

# aci_error_disable_recovery #

Data source for ACI Error Disable Recovery


## API Information ##

* `Class` - edrErrDisRecoverPol and edrEventP
* `Distinguished Named` - uni/infra/edrErrDisRecoverPol-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> Error Disabled Recovery Policy



## Example Usage ##

```hcl
data "aci_error_disable_recovery" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Error Disable Recovery.
* `annotation` - (Optional) Annotation of object Error Disable Recovery.
* `name_alias` - (Optional) Name Alias of object Error Disable Recovery.
* `err_dis_recov_intvl` - (Optional) Error Disable Recovery Interval. Sets the error disable recovery interval, which specifies the time to recover from an error-disabled state.
* `description` - (Optional) Description of object Error Disable Recovery.
