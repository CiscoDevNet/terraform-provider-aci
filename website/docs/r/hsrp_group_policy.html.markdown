---
layout: "aci"
page_title: "ACI: aci_hsrp_group_policy"
sidebar_current: "docs-aci-resource-hsrp_group_policy"
description: |-
  Manages ACI HSRP Group Policy
---

# aci_hsrp_group_policy

Manages ACI HSRP Group Policy

## Example Usage

```hcl
resource "aci_hsrp_group_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  ctrl = "preempt"
  hello_intvl  = "3000"
  hold_intvl  = "10000"
  key  = "cisco"
  name_alias  = "example"
  preempt_delay_min  = "60"
  preempt_delay_reload  = "60"
  preempt_delay_sync  = "60"
  prio  = "100"
  timeout  = "60"
  hsrp_group_policy_type = "md5"

}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object hsrp group policy.
- `annotation` - (Optional) Annotation for object hsrp group policy.
- `ctrl` - (Optional) The control state.  
  Allowed values: "preempt", "0". Default value: "0".
- `hello_intvl` - (Optional) The hello interval. Default value: "3000".

- `hold_intvl` - (Optional) The period of time before declaring that the neighbor is down. Default value: "10000".

- `key` - (Optional) The key or password used to uniquely identify this configuration object. If `key` is set, the object key will reset when terraform configuration is applied. Default value: "cisco".

- `name_alias` - (Optional) Name alias for object hsrp group policy.

- `preempt_delay_min` - (Optional) HSRP Group's Minimum Preemption delay. Default value: "0".

- `preempt_delay_reload` - (Optional) Preemption delay after switch reboot. Default value: "0".

- `preempt_delay_sync` - (Optional) Maximum number of seconds to allow IPredundancy clients to prevent preemption. Default value: "0".

- `prio` - (Optional) The QoS priority class ID. Default value: "100".

- `timeout` - (Optional) Amount of time between authentication attempts. Default value: "0".

- `hsrp_group_policy_type` - (Optional) The specific type of the object or component.  
  Allowed values: "md5", "simple". Default value: "simple".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the HSRP Group Policy.

## Importing

An existing HSRP Group Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_hsrp_group_policy.example <Dn>
```
