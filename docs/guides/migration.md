---
subcategory: "Guides"
page_title: "Migration"
description: |-
    An overview of resource and datasource migration changes
---

# Introduction

Since the first release in September 2020, the Terraform ACI provider has come a long way, adding new resources and datasources to make managing Cisco ACI environments easier. Over the years, we've moved from SDKv1 to [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2), and now to the latest [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), which is the recommended way to develop Terraform plugins and offers many advantages over SDKv2.

To take full advantage of the new features in the Terraform Plugin Framework, we had to completely rewrite the Terraform ACI provider resources and datasources. While this was a big task, it also gave us a chance to review and improve the current ACI provider. This guide explains the new features that the Terraform Plugin Framework provides and the changes we have made to the ACI provider and why we made them.

# Changes to the ACI Provider

In this section, we outline the key changes made to the Terraform ACI provider as part of the migration to the Terraform Plugin Framework. These changes aim to enhance the provider's functionality, performance, and usability. The new provider is generated from meta files, ensuring consistency and accuracy across resources and datasources. Below, we detail the specific improvements and modifications implemented during this migration.

## Cleaning and Renaming of Attributes

Over time, some attributes became outdated, were inconsistently named, or were not exposed correctly. This made the resources sometimes challenging to understand and thus use the ACI provider effectively. By cleaning and renaming these attributes, we aim to:

1. **Improve Clarity**: Ensure that attribute names are clear and descriptive, making it easier for users to understand their purpose.
2. **Enhance Consistency**: Standardize attribute names and structures across resources and datasources.
3. **Simplify Usage**: Remove outdated attributes to streamline the configuration process.

We understand the importance of maintaining backward compatibility to avoid disrupting existing configurations. This is why during a transition period, both old and new style attributes can be used. It allows users to continue using their existing configurations while gradually adopting the new style attributes. We have implemented deprecation warnings for old style attributes, informing users of the upcoming changes and encouraging them to transition to the new attribute style as soon as possible.

> For the same ACI property, the old and new style attributes cannot be used together. Doing so will cause the provider to error during configuration validation.

A downside to this approach is that the plan output will be more verbose due to known after applies for each old style attribute that is not provided when a change is detected. This is a temporary drawback which will be resolved when the deprecated attributes are removed in the next major release.

## Changed Behavior for Relations

In an ACI setup, Managed Objects (MOs) represent the different physical and logical parts within the Management Information Tree (MIT). These MOs can be linked through relationship MOs, which define how different MOs are connected. Relationship MOs act as connectors that establish links between two or more MOs, helping to organize and structure the hierarchical relationships within the MIT. In ACI, these relationships can be of two types: explicit or named.

1. **Explicit relations** need the target DN to be filled in the tDn attribute of the relationship MO, where there is only one target for the relationship. If the target DN doesn't exist, the relationship can't be formed.
2. **Named relations** need the name (identifier) attribute to be filled in the relationship MO, which triggers a resolving mechanism based on precedence order. This means the relationship could be formed with any MO in the precedence order, usually ending with a common tenant default MO.

In non-migrated resources of the Terraform Provider ACI, the relationship types are hidden from the user by allowing a DN (resource ID) or name input. The relationship type determines which attribute (name vs tDn) is added to the payload. In SDKv2, there was no enforcement (only warnings in logs) of the final plan needing to match the applied state, which allowed us to strip the name from a provided DN for named relations and add that to the payload for named relations. However, the migration to the Plugin Framework enforces the final plan to match the applied state.

This means that for the ACI Provider, when the name of the provided tDn is used for named relation MOs, we can't guarantee the DN input will match the resolved tDn, potentially causing provider panics. Because of this change, we decided to only expose configurable attributes as input for resources, so for a named relationship instead of DN we require the name attribute to be provided.

Finally, because the final plan needs to match the applied state, old style attributes that represent a DN of a named relation require to be resolved into the correct tDn. This means that the object must exist when creating the relationship; if you fail to do so, the provider will panic. This can be prevented by leveraging new style resources where the configurable attribute is used.

## Changed Child Objects in Configuration

In migrated resources, all children MOs are exposed in configuration as a map or a list of maps, instead of single attribute types (e.g., string, boolean, integer). The decision is made to have the configuration resemble the model more to provide flexibility of adding new attributes for a child without introducing breaking changes.

## Changed Child Objects in Payloads

In non-migrated resources, the children were managed via individual REST API requests, which could lead to a large number of REST API calls for a single resource when a lot of children were defined. For migrated resources, the decision is made to include all children inside a single GET/POST request to limit the number of REST API calls.

## Changed Annotation Behavior

The annotation attribute can be set for each MO that exposes the attribute, including children inside a resource. By default, the annotation is set to `orchestrator:terraform`, but this default can be overwritten at the provider level.

## Showing ID in Plan for Create

In non-migrated resources, the ID of the resource would show as "known after apply". For migrated resources, the decision is made to calculate the ID during the plan and provide this in the plan output.

## Error for Existing MO on Create

In migrated resources, we provide the provider level option (`allow_existing_on_create`) to check during the plan if an object already exists. By default, we allow existing MOs to be managed, but this option can be set to `false` to allow for a safeguard mechanism to prevent managing the same MO by different configuration parts. The drawback of this safeguard mechanism is that during a plan, an additional API call is made per resource to verify that the object does not already exist.

## Allowing Empty Input for Attributes

Terraform supports three states for any value: null (missing), unknown ("known after apply"), and known. The previous SDKs did not support null and unknown value states, which meant that before the Plugin Framework, we couldn't differentiate between an empty value and a value that wasn't provided. Because of this limitation, we were unable to update a string value to an empty string.

## Include Annotations and Tags

All migrated resources and datasources expose [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview) and [tagTag](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagTag/overview) MOs as children inside the resource when the model allows for its configuration.

## Documentation Enhancements

The documentation for migrated resources and datasources is auto-generated to provide better accuracy and consistency. The following changes have been made:

1. Class names have references to the [model documentation](https://developer.cisco.com/site/apic-mim-ref-api/?version=latest).
2. Parents of the resource are provided with references to their corresponding resources at the `parent_dn` attribute.
3. Relational children of the resource are provided with references to their corresponding child resource and the resource the child is pointing towards.
4. Non-relational children of the resource are provided with references to their corresponding resources at the bottom of the documentation.
5. Attribute default behavior is documented per attribute.
6. Attribute valid inputs are documented per attribute.
7. Applicable versions of classes and attributes are documented.
8. UI location information is provided for every configuration option.
9. Resource DN format (resource ID) is documented.
10. A minimal and full example is provided (max 2 parents).
11. Importing examples for both CLI and blocks in HCL configuration.
