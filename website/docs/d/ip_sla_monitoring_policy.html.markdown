---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_ip_sla_monitoring_policy"
sidebar_current: "docs-aci-data-source-ip_sla_monitoring_policy"
description: |-
  Data source for ACI IP SLA Monitoring Policy
---

# aci_ipsla_monitoring_policy #

Data source for ACI IP SLA Monitoring Policy


## API Information ##

* `Class` - fvIPSLAMonitoringPol
* `Distinguished Name` - uni/tn-{name}/ipslaMonitoringPol-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> IP SLA

## Example Usage ##

```hcl
data "aci_ip_sla_monitoring_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the IP SLA Monitoring Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the IP SLA Monitoring Policy.
* `annotation` - (Optional) Annotation of the IP SLA Monitoring Policy object.
* `name_alias` - (Optional) Name Alias of the IP SLA Monitoring Policy object. 
* `http_version` - (Optional) HTTP Version used for probing. Type: String.
* `type_of_service` - (Optional) Type of Service value for IPv4 packets which provides an indication of the desired Quality of Service (QoS). Allowed range is 0-255 and default value is "0". Type: String.
* `traffic_class_value` - (Optional) Traffic Class Value indicates class or priority of IPv6 packet. Type: String.
* `request_data_size` - (Optional) Minimum size of the IP SLA packet. Type: String.
* `sla_detect_multiplier` - (Optional) Detect Multiplier value for number of missed probes. Type: String.
* `sla_frequency` - (Optional) The SLA frequency value for forwarding packets. Type: String.
* `sla_port` - (Optional) The SLA destination port number. Type: String.
* `sla_type` - (Optional) The IP SLA protocol type. Type: String.
* `threshold` - (Optional) The threshold value at which the SLA is considered as failed. Type: String.
* `timeout` - (Optional) The amount of time between authentication attempts. Type: String.