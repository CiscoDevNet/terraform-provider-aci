---
layout: "aci"
page_title: "ACI: aci_rest_managed"
sidebar_current: "docs-aci-resource-rest-managed"
description: |-
  Manages ACI Model Objects via REST API calls. Any Model Object that is not supported by provider can be created/managed using this resource. Compared to the aci_rest resource, this resource can only manage a single API object (no children), but is able to read the state and therefore reconcile configuration drift.
---

# aci_rest_managed #

Manages ACI Model Objects via REST API calls. Any Model Object that is not supported by provider can be created/managed using this resource. Compared to the `aci_rest` resource, this resource can only manage a single API object (no children), but is able to read the state and therefore reconcile configuration drift.

## Example Usage ##

```hcl
resource "aci_rest_managed" "fvTenant" {
  dn         = "uni/tn-EXAMPLE_TENANT"
  class_name = "fvTenant"
  content = {
    name  = "EXAMPLE_TENANT"
    descr = "Example description"
  }
}

resource "aci_rest_managed" "mgmtConnectivityPrefs" {
  dn         = "uni/fabric/connectivityPrefs"
  class_name = "mgmtConnectivityPrefs"
  content = {
    interfacePref = "ooband"
  }
}
```

## Argument Reference ##

* `dn` - (Required) Distinguished name of object being managed including its relative name, e.g. uni/tn-EXAMPLE_TENANT.
* `class_name` - (Required) Which class object is being created. (Make sure there is no colon in the classname)
* `content` - (Optional) Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the Dn of the object created by it.

## Importing ##

This resource does not support import.
