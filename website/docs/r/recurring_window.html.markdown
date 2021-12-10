---
subcategory: "Scheduler"
layout: "aci"
page_title: "ACI: aci_recurring_window"
sidebar_current: "docs-aci-resource-recurring_window"
description: |-
  Manages ACI Recurring Window
---

# aci_recurring_window #

Manages ACI Recurring Window

## API Information ##

* `Class` - trigRecurrWindowP
* `Distinguished Named` - uni/fabric/schedp-{name}/recurrwinp-{name}

## GUI Information ##

* `Location` - Admin -> Schedulers -> Fabric -> Select Scheduler -> Create Recurring Window


## Example Usage ##

```hcl
resource "aci_recurring_window" "example" {
  scheduler_dn  = aci_trigger_scheduler.example.id
  name  = "example"
  concur_cap = "unlimited"
  day = "every-day"
  hour = "0"
  minute = "0"
  node_upg_interval = "0"
  proc_break = "none"
  proc_cap = "unlimited"
  time_cap = "unlimited"
  annotation = "Example"
}
```

## Argument Reference ##

* `scheduler_dn` - (Required) Distinguished name of parent Scheduler object.
* `name` - (Required) Name of object Recurring Window.
* `annotation` - (Optional) Annotation of object Recurring Window.
* `concur_cap` - (Optional) Maximum Concurrent Tasks. The concurrency capacity limit. This is the maximum number of tasks that can be processed concurrently. Range: "1" - "65535". Default value is "unlimited"(If user sets "0" as a value, provider will accept it but it'll set it as "unlimited"). Type: String.
* `day` - (Optional) Recurring Window Schedule Day. The day of the week that the recurring window begins. Allowed values are "Friday", "Monday", "Saturday", "Sunday", "Thursday", "Tuesday", "Wednesday", "even-day", "every-day", "odd-day", and default value is "every-day". Type: String.
* `hour` - (Optional) Schedule Hour. The hour that the recurring window begins. Range: "0" - "23". Default value is "0". 
* `minute` - (Optional) Schedule Minute. The minute that the recurring window begins. Range: "0" - "59". Default value is "0".
* `node_upg_interval` - (Optional) Delay between node upgrades. Delay between node upgrades in seconds. Range: "0" - "18000". Default value is "0".
* `proc_break` - (Optional) procBreak. A period of time taken between processing of items within the concurrency cap. Allowed Min Value: "00:00:00:00.001"(Format is DD:HH:MM:SS.Milliseconds).  Default value is "none" (If user sets "00:00:00:00.000" as a value, provider will accept it but it'll set it as "none"). Type: String.
* `proc_cap` - (Optional) procCap. Processing size capacity limitation specification. Indicates the limit of items to be processed within this window. Range: "1" - "65535". Default value is "unlimited" (If user sets "0" as a value, provider will accept it but it'll set it as "unlimited"). Type: String.
* `time_cap` - (Optional) Maximum Running Time. The processing time capacity limit. This is the maximum duration of the window. Allowed Range: "00:00:00:00.001" to "00:23:59:59.000"(Max input should be less than 24 Hours, Format is DD:HH:MM:SS.Milliseconds). Default value is "unlimited" (If user sets "00:00:00:00.000" as a value, provider will accept it but it'll set it as "unlimited"). Type: String.


## Importing ##

An existing RecurringWindow can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_recurring_window.example <Dn>
```