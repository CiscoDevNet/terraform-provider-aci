---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_ip_sla_monitoring_policy"
sidebar_current: "docs-aci-resource-aci_ip_sla_monitoring_policy"
description: |-
  Manages ACI IP SLA Monitoring Policy
---

# aci_ip_sla_monitoring_policy #

Manages ACI IP SLA Monitoring Policy

## API Information ##

* `Class` - fvIPSLAMonitoringPol
* `Distinguished Name` - uni/tn-{name}/ipslaMonitoringPol-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> IP SLA

## Example Usage ##

```hcl
resource "aci_ip_sla_monitoring_policy" "example" {
  tenant_dn             = aci_tenant.tf_tenant.id
  name                  = "example"
  http_uri              = "/index.html"
  http_version          = "HTTP/1.1"
  type_of_service       = "0"
  traffic_class_value   = "0"
  req_data_size         = "28"
  sla_detect_multiplier = "3"
  sla_frequency         = "60"
  sla_type              = "http"
  threshold             = "900"
  timeout               = "900"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the IP SLA Monitoring Policy object.
* `annotation` - (Optional) Annotation of the IP SLA Monitoring Policy object.
* `http_uri` - (Optional) URI used for HTTP probing. This is required when `sla_type` is set to "http". Type: String.
* `http_version` - (Optional) HTTP Version used for probing. Allowed values are "HTTP/1.0", "HTTP/1.1", and default value is "HTTP/1.0". Type: String.
* `type_of_service` - (Optional) Type of Service value for IPv4 packets which provides an indication of the desired Quality of Service (QoS). Allowed range is 0-255 and default value is "0". Type: String.
* `traffic_class_value` - (Optional) Traffic Class Value indicates class or priority of IPv6 packet. Allowed range is 0-255 and default value is "0". Type: String.
* `request_data_size` - (Optional) Request Data Size. Allowed range is 0-17512 and default value is "28". Type: String.
* `sla_detect_multiplier` - (Optional) Detect Multiplier value for number of missed probes. Allowed range is 1-100 and default value is "3". Type: String.
* `sla_frequency` - (Optional) The SLA frequency value for forwarding packets. Allowed range is 1-300 and default value is "60". Type: String.
* `sla_port` - (Optional) The SLA destination port number. This is required when `sla_type` is set to "tcp". Type: String.
* `sla_type` - (Optional) The SLA type. Allowed values are "http", "icmp", "l2ping", "tcp", and default value is "icmp". Type: String.
* `threshold` - (Optional) The threshold value at which the SLA is considered as failed. Allowed range is 0-604800000 and default value is "900". Type: String.
* `timeout` - (Optional) The amount of time between authentication attempts. Allowed range is 0-604800000 and default value is "900". Type: String.


## Importing ##

An existing IP SLA Monitoring Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ip_sla_monitoring_policy.example <Dn>
```