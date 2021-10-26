---
layout: "aci"
page_title: "ACI: aci_isis_domain_policy"
sidebar_current: "docs-aci-resource-isis_domain_policy"
description: |-
  Manages ACI ISIS Domain Policy and ISIS Level
---

# aci_isis_domain_policy #
Manages ACI ISIS Domain Policy amd ISIS Level

## API Information ##
* `Class` - isisDomPol & isisLvlComp
* `Distinguished Named` - uni/fabric/isisDomP-{name} & uni/fabric/isisDomP-{name}/lvl-{type}

## GUI Information ##
* `Location` - System -> System Settings -> ISIS Policy


## Example Usage ##
```hcl
resource "aci_isis_domain_policy" "example" {
  annotation = "orchestrator:terraform"
  mtu = "1492"
  redistrib_metric = "63"
  description = "from terraform"
  name_alias = "example_name_alias"
  lsp_fast_flood = "disabled"
  lsp_gen_init_intvl = "50"
  lsp_gen_max_intvl = "8000"
  lsp_gen_sec_intvl = "50"
  spf_comp_init_intvl = "50"
  spf_comp_max_intvl = "8000"
  spf_comp_sec_intvl = "50"
  isis_level_name = "example"
  isis_level_type = "l1"
}
```

## NOTE ##
Users can use the resource of type `aci_isis_domain_policy` to change the configuration of the object ISIS Domain Policy and ISIS Level. Users cannot create more than one instance of object ISIS Domain Policy and ISIS Level.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object ISIS Domain Policy.
* `mtu` - (Optional) Maximum Transmission Unit of object ISIS Domain Policy. Allowed range: "256" - "4352".
* `redistrib_metric` - (Optional) Metric used for redistributed routes. Allowed range: "1" - "63".
* `description` - (Optional) Description of object ISIS Domain Policy.
* `name_alias` - (Optional) Name alias of object ISIS Domain Policy.
* `lsp_fast_flood` - (Optional) The IS-IS Fast-Flooding of LSPs improves Intermediate System-to-Intermediate System (IS-IS) convergence time when new link-state packets (LSPs) are generated in the network and shortest path first (SPF) is triggered by the new LSPs. Allowed values are "disabled" and "enabled".
* `lsp_gen_init_intvl` - (Optional) The LSP generation initial wait interval. Allowed range: "50" - "120000".
* `lsp_gen_max_intvl` - (Optional) The LSP generation maximum wait interval. Allowed range: "50" - "120000".
* `lsp_gen_sec_intvl` - (Optional) The LSP generation second wait interval. Allowed range: "50" - "120000".
* `spf_comp_init_intvl` - (Optional) The SPF computation frequency initial wait interval.  Allowed range: "50" - "120000".
* `spf_comp_max_intvl` - (Optional) The SPF computation frequency maximum wait interval.  Allowed range: "50" - "120000".
* `spf_comp_sec_intvl` - (Optional) The SPF computation frequency second wait interval. Allowed range: "50" - "120000".
* `isis_level_name` - (Optional) (Optional) The name of ISIS Level object.
* `isis_level_type` - (Optional) The type of ISIS Level object. Allowed values are "l1" and "l2". Default value is "l1".

## Importing ##
An existing ISIS DomainPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_isis_domain_policy.example <Dn>
```