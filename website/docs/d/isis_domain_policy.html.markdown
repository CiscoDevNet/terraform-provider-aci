---
layout: "aci"
page_title: "ACI: aci_isis_domain_policy"
sidebar_current: "docs-aci-data-source-isis_domain_policy"
description: |-
  Data source for ACI ISIS Domain Policy and ISIS Level
---

# aci_isis_domain_policy #
Data source for ACI ISIS Domain Policy and ISIS Level


## API Information ##
* `Class` - isisDomPol & isisLvlComp
* `Distinguished Named` - uni/fabric/isisDomP-{name} & uni/fabric/isisDomP-{name}/lvl-{type}

## GUI Information ##
* `Location` - System -> System Settings -> ISIS Policy



## Example Usage ##
```hcl
data "aci_isis_domain_policy" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the ISIS Domain Policy.
* `annotation` - (Optional) Annotation of object ISIS Domain Policy and ISIS Level.
* `name_alias` - (Optional) Name Alias of object ISIS Domain Policy and ISIS Level.
* `mtu` - (Optional) Maximum Transmission Unit. The IS-IS Domain policy LSP MTU.
* `redistrib_metric` - (Optional) Metric. Metric used for redistributed routes.
* `description` - (Optional) Description of object ISIS Domain Policy and ISIS Level.
* `lsp_fast_flood` - (Optional) The IS-IS Fast-Flooding of LSPs improves Intermediate System-to-Intermediate System (IS-IS) convergence time when new link-state packets (LSPs) are generated in the network and shortest path first (SPF) is triggered by the new LSPs. Allowed values are "disabled" and "enabled".
* `lsp_gen_init_intvl` - (Optional) The LSP generation initial wait interval. 
* `lsp_gen_max_intvl` - (Optional) The LSP generation maximum wait interval. 
* `lsp_gen_sec_intvl` - (Optional) The LSP generation second wait interval. 
* `spf_comp_init_intvl` - (Optional) The SPF computation frequency initial wait interval. 
* `spf_comp_max_intvl` - (Optional) The SPF computation frequency maximum wait interval.  
* `spf_comp_sec_intvl` - (Optional) The SPF computation frequency second wait interval.
* `isis_level_name` - (Optional) The name of ISIS Level object. 
* `isis_level_type` - (Optional) The type of ISIS Level object.

