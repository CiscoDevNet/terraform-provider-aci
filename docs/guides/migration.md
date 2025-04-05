---
subcategory: "Guides"
page_title: "Migration to Terraform Plugin Framework"
description: |-
    An overview of resource and data source migration to Terraform Plugin Framework
---

## Introduction

Since its first release in September 2020, the Terraform ACI provider has significantly progressed, adding new resources and data sources to streamline the management of Cisco ACI environments. Over the years, we've transitioned from SDKv1 to [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), and now to the latest [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). This new Plugin Framework is the recommended way for developing Terraform plugins and offers several advantages over SDKv2.

In order to fully leverage the new features of the Terraform Plugin Framework, we undertook a complete rewrite of the Terraform ACI provider's resources and data sources. Although this was a significant effort, it presented an opportunity to review and enhance the existing ACI provider. This guide outlines the new features introduced by the Terraform Plugin Framework, along with the modifications we implemented in the ACI provider and the rationale behind these changes.

## Upgrading the ACI Terraform Provider

Upgrading the ACI Terraform provider to a new version requires careful planning to ensure a seamless transition and prevent any disruptions to your existing configurations. This guide will lead you through each step of the process, including how to back up your current state.

### Step 1: Backup Your Current State

Before making any changes, it's crucial to back up your current Terraform state. This precaution ensures that you can revert your environment to its previous state if any issues arise during the upgrade process.

1. **Local Backend**: Open a terminal and navigate to the directory where your ACI Terraform configuration files with the state are located. Copy the current state file (terraform.tfstate) to a backup location..
   
   ```bash
   cd /path/to/your/terraform/project
   cp terraform.tfstate terraform.tfstate.backup
   ```

2. **Remote Backend**: If you are using a remote backend (e.g., S3, Azure Blob Storage), ensure you have a backup of the state file from the remote location.

### Step 2: Update the Provider Version

1. **Open the Terraform Configuration File**: Open the `main.tf` or the relevant Terraform configuration file where the ACI provider is defined.

2. **Update the Provider Version**: Modify the provider block to specify the new version of the ACI provider.

   ```hcl
   provider "aci" {
     version = "x.y.z"  # Replace with the new version number
     # Other provider configuration options
   }
   ```

3. **Initialize the Configuration**: Run `terraform init` to reinitialize your configuration with the updated provider version.

   ```bash
   terraform init
   ```

### Step 3: Review the Changes

1. **Plan the Changes**: Run `terraform plan` without modifying the configuration to preview the changes that will be applied with the new provider version.

   ```bash
   terraform plan
   ```

2. **Review the Plan**: Carefully review the output of the `terraform plan` command to ensure that there are no changes.

    ```bash
    No changes. Your infrastructure matches the configuration.
    ```

### Step 4: Apply the Changes

1. **Apply the Changes**: If the plan output is satisfactory, apply the changes to upgrade to the new provider version.

   ```bash
   terraform apply
   ```

2. **Verify the Changes**: Once the apply step is complete, verify that the changes have been applied correctly and ensure your environment is functioning as expected. The state file should now reflect the redefined attributes.

### Step 5: Migrate deprecated configuration

1. **Identify the deprecated attributes**: The `terraform plan` command only displays a single warning, which can make it challenging to fully analyze deprecated attributes.

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

2. **Change the deprecated attributes**: Replace the deprecated attributes in your configuration file with the redefined attributes and execute the plan again.

    ***Old Configuration***
    ```hcl
    resource "aci_bridge_domain" "terraform_bd" {
        tenant_dn = aci_tenant.terraform_tenant.id
        name = "terraform_bd"
        ll_addr = "::"
        unk_mcast_act = "flood"
    }
    ```

    ***New Configuration***
    ```hcl
    resource "aci_bridge_domain" "terraform_bd" {
        parent_dn = aci_tenant.terraform_tenant.id
        name = "terraform_bd"
        link_local_ipv6_address = "::"
        l3_unknown_multicast_flooding = "flood"
    }
    ```

    The `terraform plan` command should not display any warnings anymore.

    ```bash
    No changes. Your infrastructure matches the configuration.

    Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
    ```

### Step 6: Cleanup

1. **Remove Backup (Optional)**: If everything is working correctly, you can remove the backup state file.

   ```bash
   rm terraform.tfstate.backup
   ```

By following these steps, you can safely upgrade the ACI Terraform provider to a new version while ensuring that you have a backup of your current state in case anything goes wrong.

## Changes to the ACI Provider

In this section, we outline the key changes made to the Terraform ACI provider as part of the migration to the Terraform Plugin Framework. These changes aim to enhance the provider's functionality, performance, and usability. The new provider is generated from meta files, ensuring consistency and accuracy across resources and data sources. Below, we detail the specific improvements and modifications implemented during this migration.

## Cleaning and Renaming of Attributes

Over time, some attributes became outdated, were inconsistently named, or were not exposed correctly, making certain resources challenging to understand and use effectively within the ACI provider. By cleaning and renaming these attributes, we aim to:

1. **Improve Clarity**: Ensure that attribute names are clear and descriptive, making it easier for users to understand their purpose.
2. **Enhance Consistency**: Standardize attribute names and structures across resources and data sources.
3. **Simplify Usage**: Remove outdated attributes to streamline the configuration process.

We understand the importance of maintaining backward compatibility to avoid disrupting existing configurations. Therefore, during the transition period, both legacy and redefined attributes are supported. This approach allows users to continue using their existing configurations while gradually adopting the redefined attributes. We have implemented deprecation warnings for legacy attributes, informing users of the upcoming changes and encouraging them to transition to the redefined attribute as soon as possible.

> It is important to note that for the same ACI property, legacy and redefined attributes cannot be used simultaneously. Attempting to do so will result in an error during configuration validation.

A downside to this approach is increased verbosity in the plan output due to known after applies for each legacy attribute not provided when a change is detected. This is a temporary drawback which will be resolved once deprecated attributes are removed in the next major release.

## Changed Behavior for Relations

In an ACI setup, Managed Objects (MOs) represent the different physical and logical parts within the Management Information Tree (MIT). These MOs can be linked through relationship MOs, which define how different MOs are connected. Relationship MOs act as connectors that establish links between two or more MOs, helping to organize and structure the hierarchical relationships within the MIT. In ACI, these relationships can be of two types: explicit or named.

1. **Explicit relations** These require the target Distinguished Name (DN) to be specified in the tDn attribute of the relationship MO, where only one target exists for the relationship. If the target DN is absent, the relationship cannot be established.
2. **Named relations** These require the name (identifier) attribute to be specified in the relationship MO, triggering a resolving mechanism based on precedence order. This allows the relationship to form with any MO in the precedence order, typically ending with a default MO in the common tenant.

In non-migrated resources of the Terraform Provider ACI, the relationship types are hidden from the user by allowing a DN (resource ID) or name input. The relationship type determines which attribute (name vs tDn) is added to the payload. In SDKv2, there was no enforcement (only warnings in logs) of the final plan needing to match the applied state, which allowed us to strip the name from a provided DN for named relations and add that to the payload for named relations. However, the migration to the Plugin Framework enforces the final plan to match the applied state.

Consequently, for the ACI Provider, using the name of the provided tDn for named relation MOs means the DN input may not match the resolved tDn, potentially causing provider errors. Due to this change, we decided to expose only configurable attributes as input for resources, requiring the name attribute instead of the DN for named relationships.

To ensure the final plan matches the applied state, legacy attributes representing a Distinguished Name (DN) of a named relation must be resolved into the correct target Distinguished Name (tDn). This requires the object to exist when establishing the relationship; failure to do so will cause the provider to panic. This issue can be avoided by using redefined resources where a configurable attribute is utilized.

## Changes to Child Objects in Configuration

In migrated resources, all child Managed Objects (MOs) are represented in the configuration as a map or a list of maps, rather than single attribute types (such as string, boolean, or integer). This approach ensures the configuration closely resembles the model, allowing for the addition of redefined attributes to a child without causing breaking changes.

## Changed Child Objects in Payloads

In non-migrated resources, child objects were managed via individual REST API requests, which could lead to a large number of REST API calls for a single resource when many children were defined. For migrated resources, the decision is made to include all children inside a single GET/POST request to limit the number of REST API calls.

## Changed Annotation Behavior

The annotation attribute can be set for each MO that exposes the attribute, including children inside a resource. By default, the annotation is set to `orchestrator:terraform`, but this default can be overwritten at the provider level.

```hcl
provider "aci" {
  username   = "admin"
  password   = "password"
  url        = "https://my-cisco-aci.com"
  annotation = "provider-level-annotation-overwrite"
}
```

## Showing ID in Plan for Create

For non-migrated resources, the ID of the resource appears as "known after apply." In the case of migrated resources, the ID is calculated during the planning phase and included in the plan output.

## Error for Existing MO on Create

Migrated resources offer a provider-level option (`allow_existing_on_create`) to check during the plan if an object already exists. By default, existing MOs can be managed, but this option can be set to `false`, providing a safeguard mechanism to prevent managing the same MO by different configuration parts. The drawback of this safeguard mechanism is an additional API call per resource during the plan to verify that the object does not already exist.

## Allowing Empty Input for Attributes

Terraform supports three states for any value: null (missing), unknown ("known after apply"), and known. Previous SDKs did not support null and unknown states, which prevented differentiating between an empty value and a non-provided value. Because of this limitation, updating a string value to an empty string was not possible before the Plugin Framework.

## Include Annotations and Tags

All migrated resources and data sources expose [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview) and [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview) MOs as children inside the resource when the model allows for its configuration.

## Documentation Enhancements

The documentation for migrated resources and data sources is automatically generated to enhance relevancy and consistency. The following changes have been made:

1. Class names have references to the [model documentation](https://developer.cisco.com/site/apic-mim-ref-api/?version=latest).
2. Parents of the resource are provided with references to their corresponding resources at the `parent_dn` attribute.
3. Relational children of the resource include references to both their corresponding child resource and the target resource they are associated with.
4. Non-relational children of the resource are provided with references to their corresponding resources at the bottom of the documentation.
5. The default behaviour of each attribute is documented individually.
6. The valid inputs for each attribute are documented individually.
7. Applicable versions of classes and attributes are documented.
8. UI location information is provided for each configuration option.
9. Resource DN format (resource ID) is documented.
10. Both minimal and full examples are provided, featuring up to two parent resources when applicable.
11. Examples for importing resources using both CLI commands and HCL configuration blocks are provided.
