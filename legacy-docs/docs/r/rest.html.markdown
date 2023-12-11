---
layout: "aci"
page_title: "ACI: aci_rest"
sidebar_current: "docs-aci-resource-aci_rest"
description: |-
  Manages ACI Model Objects via REST API calls. Any Model Object that is not supported by provider can be created/managed using this resource.
---

# aci_rest

Manages ACI Model Objects via REST API calls. Any Model Object that is not supported by provider can be created/managed using this resource.

## Example Usage

```hcl
resource "aci_tenant" "tenant_for_rest_example" {
  name        = "tenant_for_rest"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_rest" "rest_l3_ext_out" {
  path       = "/api/node/mo/${aci_tenant.tenant_for_rest_example.id}/out-test_ext.json"
  class_name = "l3extOut"
  content = {
    "name" = "test_ext"
  }
}

resource "aci_rest" "madebyrest_yaml" {
  path       = "/api/mo/uni/tn-sales_main.json"
  payload = <<EOF
{
        "fvTenant": {
          "attributes": {
            "name": "sales_main",
            "descr": "Sales department json"
          }
        }
}
  EOF
}

resource "aci_rest" "rest_yaml" {
  path       = "/api/mo/uni/tn-sales.json"
  payload = <<EOF
  fvTenant:
        attributes:
          name: sale
          descr: Sales department
  EOF
}
```

## Argument Reference

- `path` - (Required) ACI path where object should be created. Starting with api/node/mo/{parent-dn}(if applicable)/{rn of object}.json
- `class_name` - (Optional) Which class object is being created. (Make sure there is no colon in the classname )
- `content` - (Optional) Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.
- `payload` - (Optional) Freestyle JSON or YAML payload which can directly be passed to the REST endpoint added in path. Either of content or payload is required.
- `dn` - (Optional) Distinguished name of object being managed.

**NOTE:** We don't set the Status field explicitly, as it creates an issue with the relation objects. If you have requirement to pass the status field, pass it in the content.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the object created by it.

## Importing

This resource does not support import.
