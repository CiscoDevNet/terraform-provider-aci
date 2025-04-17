---
subcategory: "Guides"
page_title: "Migration to Terraform Plugin Framework"
description: |-
    An overview of resource and data source migration to Terraform Plugin Framework
---

## Introduction

Since its first release in July 2019, the Terraform ACI provider has significantly progressed, adding new resources and data sources to streamline the management of Cisco ACI environments. Over the years, we've transitioned from SDKv1 to [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), and now to the latest [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). This new Plugin Framework is the recommended way for developing Terraform plugins and offers several advantages over SDKv2.

In order to fully leverage the new features of the Terraform Plugin Framework, we decided to undertake a complete rewrite of the Terraform ACI provider's resources and data sources. Although this is a significant effort, it presented an opportunity to review and enhance the existing ACI provider and how the provider is developed. 

This complete rewrite will stretch a series of releases starting with [v2.16.0](https://github.com/CiscoDevNet/terraform-provider-aci/releases/tag/v2.16.0) where all new resources and data sources will be implemented using the new Terraform Plugin Framework and pre-existing resources and data sources will be slowly migrated into the Terraform Plugin Framework. Migrated resources and data sources contain a note that specifies it has been migrated at the top of their documentation. Once all pre-existing resources are migrated to the Terraform Provider Framework, a major release of the provider will be created to remove legacy attributes.

This guide provides detailed steps for migrating from previous versions, outlines key optimization steps to enhance performance, and aims to facilitate a smoother transition for users adapting to these updates. Additionally, it outlines how the migrated resources use the new features introduced by the Terraform Plugin Framework, along with the modifications we implemented in the ACI provider and the rationale behind these changes.

## Upgrading the ACI Terraform Provider

Upgrading the ACI Terraform provider to a new version requires careful planning to ensure a seamless transition and prevent any disruptions to your existing configurations. This guide will lead you through each step of the process, including how to back up your current state.

### Step 1: Backup Your Current State

Before making any changes, it's crucial to back up your current Terraform state. This precaution ensures that you can revert your environment to its previous state if any issues arise during the upgrade process.

1. **Local Backend:** Open a terminal and navigate to the directory where your ACI Terraform configuration files with the state are located. Copy the current state file (terraform.tfstate) to a backup location.

    ```bash
    cd /path/to/your/terraform/project
    cp terraform.tfstate /path/to/your/backup/terraform.tfstate.backup
    ```

2. **Remote Backend:** If you are using a remote backend (e.g., S3, Azure Blob Storage), ensure you have a backup of the state file from the remote location.

### Step 2: Update the Provider Version

The version meta-argument specifies the version constraint for a provider. The version constraint can be specified in the [provider configuration](https://developer.hashicorp.com/terraform/language/providers/configuration) or in the [required_providers block](https://developer.hashicorp.com/terraform/language/providers/requirements). The [version](https://developer.hashicorp.com/terraform/language/providers/configuration#version-deprecated) argument in the provider configuration is deprecated, and will be removed in future Terraform versions.

1. **Open the Terraform Configuration File:** Open the `main.tf` or the relevant Terraform configuration file where the ACI provider version constraint is defined.

2. **Update the Provider Version:** Terraform recommends setting [version constraints](https://developer.hashicorp.com/terraform/language/expressions/version-constraints) for providers. Modify the terraform.required_providers or provider block to specify the new version constraint for the ACI provider. 

    ***Version constraint format***
    ```hcl
    terraform {
      required_providers {
        aci = {
          source = "ciscodevnet/aci"
          version = "~> x.y.z"  # Replace with the new version number
        }
      }
    }
    ```

    ***Deprecated Terraform version constraint format***
    ```hcl
    provider "aci" {
      version = "~> x.y.z"  # Replace with the new version number
      # Other provider configuration options
    }
    ```

3. **Initialize the Configuration:** Run `terraform init -upgrade` to reinitialize your configuration with the updated provider version.

    ```bash
    terraform init -upgrade
    ```

### Step 3: Review the Changes

1. **Plan the Changes:** Run `terraform plan` without modifying the configuration to preview the changes that will be applied with the new provider version.

    ```bash
    terraform plan
    ```

2. **Review the Plan:** Carefully review the output of the `terraform plan` command to ensure that there are no changes.

    ~>  Warnings are expected to be displayed in the plan output for legacy attributes, this topic will be addressed in [Step 4: Migrate deprecated configuration](#step-4-migrate-deprecated-configuration).

    ```bash
    No changes. Your infrastructure matches the configuration.
    ```

The state file does not reflect the changes yet, because a refresh has not taken place. This can be done with the `terraform refresh` or `terraform apply -refresh-only` commands, which do not modify the objects in APIC, but only modifies the state file. For more information see [refresh](https://developer.hashicorp.com/terraform/cli/commands/refresh) and [planning-modes](https://developer.hashicorp.com/terraform/cli/commands/plan#planning-modes) documentation of Terraform.

### Step 4: Migrate deprecated configuration

1. **Identify the deprecated attributes:** The `terraform plan` command only displays a single warning, which can make it challenging to fully analyze deprecated attributes.

    ```bash
    No changes. Your infrastructure matches the configuration.

    Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
    ╷
    │ Warning: Attribute Deprecated
    │ 
    │   with aci_bridge_domain.terraform_bd,
    │   on main.tf line 70, in resource "aci_bridge_domain" "terraform_bd":
    │   70:   tenant_dn = aci_tenant.terraform_tenant.id
    │ 
    │ Attribute 'tenant_dn' is deprecated, please refer to 'parent_dn' instead. The attribute will be removed in the next major version of the provider.
    │ 
    │ (and 3 more similar warnings elsewhere)
    ```

    To display the other warnings the `terraform validate` command can be used.

    ```bash
    terraform validate -json | jq '.diagnostics[] | {detail: .detail, filename: .range.filename, start_line: .range.start.line}'
    {
        "detail": "Attribute 'unk_mcast_act' is deprecated, please refer to 'l3_unknown_multicast_flooding' instead. The attribute will be removed in the next major version of the provider.",
        "filename": "main.tf",
        "start_line": 73
    }
    {
        "detail": "Attribute 'll_addr' is deprecated, please refer to 'link_local_ipv6_address' instead. The attribute will be removed in the next major version of the provider.",
        "filename": "main.tf",
        "start_line": 72
    }
    ```

2. **Change the deprecated attributes:** Replace the deprecated attributes in your configuration file with the redefined attributes and execute the plan again.

    ***Old Configuration***
    ```hcl
    resource "aci_bridge_domain" "terraform_bd" {
      tenant_dn     = aci_tenant.terraform_tenant.id
      name          = "terraform_bd"
      ll_addr       = "::"
      unk_mcast_act = "flood"
    }
    ```

    ***New Configuration***
    ```hcl
    resource "aci_bridge_domain" "terraform_bd" {
      parent_dn                     = aci_tenant.terraform_tenant.id
      name                          = "terraform_bd"
      link_local_ipv6_address       = "::"
      l3_unknown_multicast_flooding = "flood"
    }
    ```

    Continue updating the deprecated attributes in the configuration file until all warnings have been resolved. The `terraform plan` command should not display any warnings anymore.

    ```bash
    No changes. Your infrastructure matches the configuration.

    Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
    ```

### Step 5: Cleanup

1. **Remove Backup (Optional):** If everything is working correctly, you can remove the backup state file.

    ```bash
    rm terraform.tfstate.backup
    ```

By following these steps, you can safely upgrade the ACI Terraform provider to a new version while ensuring that you have a backup of your current state in case anything goes wrong.

## Optimization

1. **Add relationships to parent resource:** Reduces the amount of resources in the configuration by including children inside their parent resources. This will decrease the execution time by decreasing the size of the terraform graph and reducing the amount of REST API calls made towards APIC. See section [Changed Child Objects in Payloads](#changed-child-objects-in-payloads) for more details.

    ***Old Configuration***
    ```hcl
    resource "aci_application_epg" "terraform_epg" {
      parent_dn = aci_application_profile.terraform_ap.id
      name      = "terraform_epg"
    }

    resource "aci_epg_to_domain" "terraform_epg_to_domain" {
      application_epg_dn = aci_application_epg.terraform_epg.id
      tdn                = "uni/phys-phys"
    }
    ```

    ***New Configuration***
    ```
    resource "aci_application_epg" "terraform_epg" {
        parent_dn = aci_application_profile.terraform_ap.id
        name      = "terraform_epg"
        relation_to_domains = [
          {
              target_dn = "uni/phys-phys"
          }
        ]
    }
    ```

    !> Use a child resource or define as nested attribute in the parent resource. Unexpected behaviour will occur when a child resource is used in combination with a nested attribute in the parent resource.

    The child configuration should not be destroyed when moving the configuration from child to parent resource. The child resource should be unmanaged (removed from state) before executing the plan. Terraform will display destroy intent in the `terraform plan` output when executed without unmanaging the child resource first.

    ```bash
    Terraform will perform the following actions:                            

    # aci_epg_to_domain.terraform_epg_to_domain will be destroyed
    # (because aci_epg_to_domain.terraform_epg_to_domain is not in configuration)
    - resource "aci_epg_to_domain" "terraform_epg_to_domain" {
      - application_epg_dn    = "uni/tn-terraform_tenant/ap-terraform_ap/epg-terraform_epg" -> null
      - tdn                   = "uni/phys-phys" -> null
      ...
      - vmm_id                = jsonencode({})
        # (7 unchanged attributes hidden)
    }

    Plan: 0 to add, 0 to change, 1 to destroy.
    ```

    The [terraform state rm](https://developer.hashicorp.com/terraform/cli/commands/state/rm) command `terraform state rm aci_epg_to_domain.terraform_epg_to_domain` should be used to unmanage without destroying a resource. The output should reflect the succesful removal of the resource from state.

    ```bash
    Removed aci_epg_to_domain.terraform_epg_to_domain
    Successfully removed 1 resource instance(s).
    ```

    The `terraform plan` output should reflect no changes, when the resource is unmanaged.
    
    ```bash
    No changes. Your infrastructure matches the configuration.

    Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
    ```

    The state file does not reflect the changes yet, because a refresh has not taken place. This can be done with the `terraform refresh` or `terraform apply -refresh-only` commands, which do not modify the objects in APIC, but only modifies the state file. For more information see [refresh](https://developer.hashicorp.com/terraform/cli/commands/refresh) and [planning-modes](https://developer.hashicorp.com/terraform/cli/commands/plan#planning-modes) documentation of Terraform.

## Changes to the ACI Provider

In this section, we outline the key changes made to the Terraform ACI provider as part of the migration to the Terraform Plugin Framework. These changes aim to enhance the provider's functionality, performance, and usability. The new provider is generated from meta files, ensuring consistency and accuracy across resources and data sources. Below, we detail the specific improvements and modifications implemented during this migration.

## Cleaning and Renaming of Attributes

Over time, some attributes became outdated, were inconsistently named, or were not exposed correctly, making certain resources challenging to understand and use effectively within the ACI provider. By cleaning and renaming these attributes, we aim to:

1. **Improve Clarity:** Ensure that attribute names are clear and descriptive, making it easier for users to understand their purpose.
2. **Enhance Consistency:** Standardize attribute names and structures across resources and data sources.
3. **Simplify Usage:** Remove outdated attributes to streamline the configuration process.

We understand the importance of maintaining backward compatibility to avoid disrupting existing configurations. Therefore, during the transition period, both legacy and redefined attributes are supported. This approach allows users to continue using their existing configurations while gradually adopting the redefined attributes. We have implemented deprecation warnings for legacy attributes, informing users of the upcoming changes and encouraging them to transition to the redefined attribute as soon as possible.

> It is important to note that for the same ACI property, legacy and redefined attributes cannot be used simultaneously. Attempting to do so will result in an error during configuration validation.

A downside to this approach is the increased verbosity in the plan output when a change is detected due to the display of the "known after apply" state for each non-provided legacy attribute. This is a temporary drawback which will be resolved once deprecated attributes are removed in the next major release.

## Changed Behavior for Relations

In an ACI fabric, Managed Objects (MOs) represent the different physical and logical parts within the Management Information Tree (MIT). These MOs can be linked through relationship MOs, which define how different objects are connected. Relationship MOs act as connectors that establish links between two or more objects, helping to organize and structure the hierarchical relationships within the MIT. In ACI, these relationships can be of two types: explicit or named.

1. **Explicit relationships** These relationships require the target Distinguished Name (DN) to be specified in the tDn attribute of the relationship MO, where only one target exists for the relationship. If the target DN is absent, the relationship cannot be established.
2. **Named relationships** These relationships require the name (identifier) attribute to be specified in the relationship MO, triggering a resolving mechanism based on a precedence order. This allows the relationship to form with the first MO that exists in the precedence order.

In non-migrated resources of the Terraform Provider ACI, the relationship types are hidden from the user by allowing a DN (resource ID) or name input. The relationship type determines which attribute (name vs tDn) is added to the payload. In SDKv2, there was no enforcement (only warnings in logs) of the final plan needing to match the applied state, which allowed us to extract the name from a provided DN for named relationships and use that extracted name in the relationship payload. However, the Plugin Framework enforces that the final plan values have to match the applied state.

Consequently, for the ACI Provider, using the name of the provided tDn for the named relationship MOs means that the user provided DN may not match the resolved tDn, potentially causing provider errors. Due to this change, we decided to expose only configurable attributes as input for resources, requiring the name attribute instead of the DN for named relationships.

To ensure the final plan matches the applied state, legacy attributes representing a Distinguished Name (DN) of a named relationship must be resolved into the correct target Distinguished Name (tDn). This requires the object to exist when establishing the relationship; the absence of the object will cause the provider to panic. See provider panic example below.

```bash
╷
│ Warning: Attribute Deprecated
│
│   with aci_bridge_domain.terraform_bd,
│   on main.tf line 59, in resource "aci_bridge_domain" "terraform_bd":
│   59:     relation_fv_rs_bd_to_ep_ret = "uni/tn-terraform_tenant/epRPol-ep_pol"
│
│ Attribute 'relation_fv_rs_bd_to_ep_ret' is deprecated, please refer to 'relation_to_end_point_retention_policy' instead. The attribute will
│ be removed in the next major version of the provider.
╵
╷
│ Error: Provider produced inconsistent result after apply
│
│ When applying changes to aci_bridge_domain.terraform_bd, provider "provider[\"registry.terraform.io/ciscodevnet/aci\"]" produced an
│ unexpected new value: .relation_fv_rs_bd_to_ep_ret: was cty.StringVal("uni/tn-terraform_tenant/epRPol-ep_pol"), but now
│ cty.StringVal("uni/tn-common/epRPol-default").
╵
```

This issue can be avoided by using the redefined resources with a configurable name attribute or by ensuring that the tDn can be resolved during the apply phase.

## Changes to Child Objects in Configuration

In migrated resources, all child Managed Objects (MOs) are represented in the configuration as a map or a list of maps, rather than single attribute types (such as string, boolean, or integer). This approach ensures the configuration closely resembles the model, allowing for the addition of redefined attributes to a child without causing breaking changes.

## Changed Child Objects in Payloads

In non-migrated resources, child objects (including relationship objects) were managed via individual REST API requests, which could lead to a large number of REST API calls for a single resource containing many children. For migrated resources, a single GET/POST request that include all children objects is used to minimize the number of REST API calls made by each resource.

## Changed Annotation Behavior

The annotation attribute can be set for each MO that exposes the attribute, including children inside a resource. By default, the annotation is set to `orchestrator:terraform`, but this default can be overwritten at different levels (provider, resource, child attribute).

```hcl
provider "aci" {
  username   = "admin"
  password   = "password"
  url        = "https://my-cisco-aci.com"
  annotation = "provider-level-annotation-overwrite"
}

resource "aci_application_epg" "example" {
  parent_dn  = aci_application_profile.example.id
  name       = "test_name"
  annotation = "resource-level-annotation-overwrite"
  relation_to_bridge_domain = {
    annotation         = "child-level-annotation-overwrite"
    bridge_domain_name = aci_bridge_domain.example.name
  }
}
```

## Showing ID in Plan for Create

For non-migrated resources, the ID of the resource appears as "known after apply". In the case of migrated resources, the ID is computed during the planning phase and included in the plan output to reduce the cascading attributes displaying a "known after apply".

## Error for Existing MO on Create

Migrated resources support the [allow_existing_on_create](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs#schema) provider-level option to check during the plan if an object already exists. By default, existing MOs will be automatically managed on apply, as `allow_existing_on_create` is set to `true` by default.  This option can be set to `false`, providing a mechanism to prevent the management of the same MO by separate configurations. The drawback of this mechanism is that it requires an additional API call per resource during the plan to verify that the object does not already exist.

## Allowing Empty Input for Attributes

The attributes of migrated resources in the Terraform ACI provider now support null, unknown, and known states which can be used to reset values to empty or to enforce that children objects are not present. Terraform supports three states for any value: null (missing), unknown ("known after apply"), and known. Previous SDKs did not support null and unknown states, which prevented differentiating between an empty value and a non-provided value. Because of this limitation, updating a string value to an empty string, updating a map to an empty map or a list of maps to an empty list was not possible before the Plugin Framework but is now supported in the migrated resources.

## Include Annotations and Tags

All migrated resources and data sources expose [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview) and [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview) MOs as children and children of children inside the resource when the model allows for its configuration. See the [ACI documentation](https://www.cisco.com/c/en/us/td/docs/dcn/aci/apic/5x/system-management-configuration/cisco-apic-system-management-configuration-guide-52x/m-alias-annotations-and-tags.html) for more explanation about these MOs.

## Documentation Enhancements

The documentation for migrated resources and data sources is automatically generated to enhance relevancy and consistency. The following changes have been made:

1. Class names have references to the [model documentation](https://developer.cisco.com/site/apic-mim-ref-api/?version=latest).
2. Parents of the resource are provided with references to their corresponding resources (including links to the documentation pages) at the `parent_dn` attribute.
3. Relational children of the resource include references to both their corresponding child resource and the target resource they are associated with.
4. Non-relational children of the resource are provided with references to their corresponding resources at the bottom of the documentation.
5. The default behaviour of each attribute is documented individually.
6. The valid inputs for each attribute are documented individually.
7. Applicable versions of classes and attributes are documented.
8. UI location information is provided for each configuration option.
9. Resource DN format (resource ID) is documented.
10. Both minimal and full examples are provided, featuring up to two parent resources when applicable.
11. Examples for importing resources using both CLI commands and HCL configuration blocks are provided.
12. Migrated resources contain a note that specifies it has been migrated at the top of the resource documentation.
