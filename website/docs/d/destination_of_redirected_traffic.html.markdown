---
layout: "aci"
page_title: "ACI: aci_destination_of_redirected_traffic"
sidebar_current: "docs-aci-data-source-destination_of_redirected_traffic"
description: |-
  Data source for ACI Destination of redirected traffic
---

# aci_destination_of_redirected_traffic

Data source for ACI Destination of redirected traffic

## Example Usage

```hcl

data "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn  = aci_service_redirect_policy.example.id
  ip                          = "1.2.3.4"
}

```

## Argument Reference

- `service_redirect_policy_dn` - (Required) Distinguished name of parent Service Redirect Policy object.
- `ip` - (Required) The IP address.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Destination of redirected traffic.
- `annotation` - (Optional) Annotation for object destination of redirected traffic.
- `destination` - (Optional) Specifies the description of a policy component.
- `dest_name` - (Optional) The destination name to which data was exported. This utility creates a summary report containing configuration information, logs and diagnostic data that will help TAC in troubleshooting and resolving a technical issue.
- `ip2` - (Optional) IP2 for object destination of redirected traffic.
- `mac` - (Optional) The MAC address.
- `name_alias` - (Optional) Name alias for object destination of redirected traffic.
- `pod_id` - (Optional) The Pod identifier.
