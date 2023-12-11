---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_hsrp_group_policy"
sidebar_current: "docs-aci-data-source-aci_hsrp_group_policy"
description: |-
  Data source for ACI HSRP Group Policy
---

# aci_hsrp_group_policy

Data source for ACI HSRP Group Policy

## Example Usage

```hcl
data "aci_hsrp_group_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object hsrp group policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the HSRP Group Policy.
- `annotation` - (Optional) Annotation for object hsrp group policy.
- `description` - (Optional) Description for object hsrp group policy.
- `ctrl` - (Optional) The control state.
- `hello_intvl` - (Optional) The hello interval.
- `hold_intvl` - (Optional) The period of time before declaring that the neighbor is down.
- `key` - (Optional) The key or password used to uniquely identify this configuration object.
- `name_alias` - (Optional) Name alias for object hsrp group policy.
- `preempt_delay_min` - (Optional) HSRP Group's Minimum Preemption delay.
- `preempt_delay_reload` - (Optional) Preemption delay after switch reboot.
- `preempt_delay_sync` - (Optional) Maximum number of seconds to allow IPredundancy clients to prevent preemption.
- `prio` - (Optional) The QoS priority class ID.
- `timeout` - (Optional) Amount of time between authentication attempts.
- `hsrp_group_policy_type` - (Optional) The specific type of the object or component.
