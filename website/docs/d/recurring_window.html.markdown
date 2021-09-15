---
layout: "aci"
page_title: "ACI: aci_recurring_window"
sidebar_current: "docs-aci-data-source-recurring_window"
description: |-
  Data source for ACI Recurring Window
---

# aci_recurring_window #

Data source for ACI Recurring Window


## API Information ##

* `Class` - trigRecurrWindowP
* `Distinguished Named` - uni/fabric/schedp-{name}/recurrwinp-{name}

## GUI Information ##

* `Location` - Admin -> Schedulers -> Fabric -> Select Scheduler -> Create Recurring Window 



## Example Usage ##

```hcl
data "aci_recurring_window" "example" {
  scheduler_dn  = aci_scheduler.example.id
  name  = "example"
}
```

## Argument Reference ##

* `scheduler_dn` - (Required) Distinguished name of parent Scheduler object.
* `name` - (Required) name of object Recurring Window.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Recurring Window.
* `annotation` - (Optional) Annotation of object Recurring Window.
* `name_alias` - (Optional) Name Alias of object Recurring Window.
* `concur_cap` - (Optional) Maximum Concurrent Tasks. The concurrency capacity limit. This is the maximum number of tasks that can be processed concurrently.
* `day` - (Optional) Recurring Window Schedule Day. The day of the week that the recurring window begins.
* `hour` - (Optional) Schedule Hour. The hour that the recurring window begins.
* `minute` - (Optional) Schedule Minute. The minute that the recurring window begins.
* `node_upg_interval` - (Optional) Delay between node upgrades. Delay between node upgrades in seconds.
* `proc_break` - (Optional) procBreak. A period of time taken between processing of items within the concurrency cap.
* `proc_cap` - (Optional) procCap. Processing size capacity limitation specification. Indicates the limit of items to be processed within this window.
* `time_cap` - (Optional) Maximum Running Time. The processing time capacity limit. This is the maximum duration of the window.
