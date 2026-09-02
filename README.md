
# Cisco ACI Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 
  - v0.12 and higher (ACI Provider v1.0.0 or higher)
  - v0.11.x or below (ACI Provider v0.7.1 or below)

- [Go](https://golang.org/doc/install) Latest Version

## Building The Provider

Clone this repository to: `$GOPATH/src/github.com/CiscoDevNet/terraform-provider-cisco-aci`.

```sh
$ mkdir -p $GOPATH/src/github.com/CiscoDevNet; cd $GOPATH/src/github.com/CiscoDevNet
$ git clone https://github.com/CiscoDevNet/terraform-provider-aci.git
```

Enter the provider directory and run dep ensure to install all the dependancies. After, that run make build to build the provider binary.

```sh
$ cd $GOPATH/src/github.com/CiscoDevNet/terraform-provider-aci
$ dep ensure
$ make build

```

## Using The Provider

If you are building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/cli/plugins/index.html) After placing it into your plugins directory, run `terraform init` to initialize it.

ex.
```hcl
terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  # cisco-aci user name
  username = "admin"
  # cisco-aci password
  password = "password"
  # cisco-aci url
  url      = "https://my-cisco-aci.com"
  insecure = true
  proxy_url = "https://proxy_server:proxy_port"
}

resource "aci_tenant" "test-tenant" {
  name        = "test-tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "test-app" {
  tenant_dn   = aci_tenant.test-tenant.id
  name        = "test-app"
  description = "This application profile is created by terraform"
}
```
Note : If you are facing the issue of `invalid character '<' looking for beginning of value` while running `terraform apply`, use signature based authentication in that case, or else use `-parallelism=1` with `terraform plan` and `terraform apply` to limit the concurrency to one thread.

```
terraform plan -parallelism=1
terraform apply -parallelism=1
```  


```hcl
  provider "aci" {
      # cisco-aci user name
      username = "admin"
      # private key path
      private_key = "path to private key"
      # Certificate Name
      cert_name = "user-cert"
      # cisco-aci url
      url      = "https://my-cisco-aci.com"
      insecure = true
  }
```

Note: The value of "cert_name" argument must match the name of the certificate object attached to the APIC user (aaaUserCert) used for signature-based authentication

## Developing The Provider

Currently the ACI provider is a [muxed provider](https://developer.hashicorp.com/terraform/plugin/mux) which allows us to simultaneously serve [terraform-plugin-sdk/v2](https://developer.hashicorp.com/terraform/plugin/sdkv2) and [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework) provider SDK implementations. This adds some complexity to the development process, but allows us to leverage the new capabilities that terraform-plugin-framework provides, while working on a migration strategy for existing resources.

### Pre-Requirements

1. Install latest version of [Go](http://www.golang.org)

### Existing resources and data-sources developed with terraform-plugin-sdk/v2

* Existing resources and data-sources are located in the [aci](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/aci) directory.

* Changes are made directly in the `provider.go`, `resource_*.go`, `data_source_*.go` and `utils` files. The [aci-go-client](https://github.com/ciscoecosystem/aci-go-client) is leveraged to construct payload constructs and handle REST communication towards the APIC.

* Documentation is manually maintained in the [legacy-docs](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/legacy-docs) directory and are copied automatically up on execution of `go generate` command.

* Examples are manually maintained in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/examples) directory, where each resource has it's own directory.

#### Manual Development Process

The below steps should be followed for developing `terraform-plugin-sdk/v2` resources and data-sources:

1. Create issue (if not created yet) and comment that you will be working on the issue.

2. Fork the terraform-provider-aci repository.

3. Clone the forked code to your local machine.

4. Make changes to the files manually.
    * Code changes
    * Examples changes
    * Documentation changes 

5. Run `go generate` in the root of the local repository where the `main.go` is located. Assure that the documentation is copied to the docs folder and no changes are made to files in the [internal/provider](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider) directory.

5. Test the code.

6. Create PR for the code and request review from active maintainers.

7. Review process

### Resources and data-sources developed with terraform-plugin-framework

The terraform-plugin-framework implementation is located in
[internal/provider](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider).
Reusable generated ACI class models are located in
[`internal/provider/models`](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider/models).

The generator is being migrated incrementally. A generated marker does not by
itself mean that the current generator owns a file: legacy generated files keep
their marker until their templates have been replaced. The exact files owned by
normal generation are recorded in [`gen/managed_files.json`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/main/gen/managed_files.json).
Currently, normal generation owns the class model files. Resource and
data-source implementations, tests, documentation, and examples remain static
or legacy outputs until their replacement templates are introduced.

Generated files should not be edited manually. Fix their normalized metadata,
class definition, generator logic, or active template and regenerate the full
output set. Provider bootstrap and registry files, `docs/index.md`, and
`examples/provider/provider.tf` are intentionally static. The complete current
ownership contract is documented in
[`gen/docs/templates.md`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/main/gen/docs/templates.md),
and the generated class-model contract is documented in
[`gen/docs/model.md`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/main/gen/docs/model.md).

#### Generator process

The generator entry point is
[`gen/generator.go`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/main/gen/generator.go).
It uses the following inputs:

- `gen/meta/<className>.json` contains the APIC class metadata.
- `gen/definitions/<className>.yaml` contains per-class overrides, while
  `gen/definitions/global.yaml` contains shared configuration.
- `gen/templates/` contains active templates.
- `gen/templates/legacy_templates/` contains migration references and is not
  loaded by the generator.

Run normal generation from the repository root:

```sh
go generate
```

This regenerates every manifest-managed output. To retrieve selected metadata
before generation, set `GEN_ACI_TF_META_CLASSES` to a comma-separated class
list. `GEN_ACI_TF_META_HOST` optionally selects the metadata host; when it is
unset, the DevNet metadata host is used.

```sh
GEN_ACI_TF_META_CLASSES=fvTenant,fvBD go generate
```

`internal/provider/annotation_unsupported.go` is a refresh-only generated file.
Refresh it explicitly from the complete ACI metadata:

```sh
GEN_ANNOTATION_UNSUPPORTED=1 go generate
```

When changing a definition produced by
`gen/scripts/migrate_class_definitions.go`, update the migration source and
regenerate the definition instead of changing only its emitted YAML. After
generation, review all generated changes and run the relevant unit and
acceptance tests before opening a pull request.

### Troubleshooting 

#### Go Generate Fail

If you encounter an error message while generating the resources using go generate, a few steps can be followed:

1. Make sure that you're running the latest version of Go

2. Update dependencies and populate the vendor directory: If you're using Go modules, you can update your dependencies and populate your vendor directory by running the following commands in your terminal:

```sh
go mod tidy
go mod vendor
```

The go mod tidy command will clean up unused dependencies and add missing ones. The go mod vendor command will copy all dependencies into the vendor directory.

3. Disable vendor mode (if necessary): If the error still persists, consider disabling vendor mode by setting GOFLAGS="-mod=mod", and then run go get again:

```sh
export GOFLAGS="-mod=mod"
```

This will force Go to fetch the package directly, regardless of the vendor directory. After the package is downloaded, you can switch back to vendor mode if needed.

#### Missing packages

If you encounter an error indicating that the golang.org/x/text/language package is missing from your vendor directory, you can fetch it by following these steps:

1. Disable vendor mode: Go operates in vendor mode when the -mod=vendor flag is set. You'll need to disable vendor mode to fetch packages directly. Run the following command in your terminal:

```sh
export GOFLAGS="-mod=mod"
```

2. Fetch the missing package: Now that vendor mode is disabled, you can fetch the missing package by running:

```sh
go get golang.org/x/text/language
```

This command tells Go to fetch the golang.org/x/text/language package directly, regardless of the vendor directory

3. Re-enable vendor mode (if necessary): If you wish to switch back to vendor mode, you can do so by running:

```sh
export GOFLAGS="-mod=vendor"
```


### Compiling

To compile the provider, run `make build`. This will build the provider with sanity checks present in scripts directory and put the provider binary in `$GOPATH/bin` directory.

<strong>Important: </strong>To successfully use the provider you need to follow these steps:

- Copy or Symlink the provider from the `$GOPATH/bin` to `~/.terraform.d/plugins/terraform.local/CiscoDevNet/aci/<Version>/<architecture>/` for example:
  ```bash
  ln -s ~/go/bin/terraform-provider-aci ~/.terraform.d/plugins/terraform.local/CiscoDevNet/aci/2.3.0/linux_amd64/terraform-provider-aci
  ```

- Edit the Terraform Provider Configuration to use the local provider.

  ```hcl
  terraform {
    required_providers {
      aci = {
        source = "terraform.local/CiscoDevNet/aci"
        version = "2.3.0"
      }
    }
  }
  ```

<strong>NOTE:</strong> Currently only resource properties supports the reflecting manual changes made in CISCO ACI. Manual changes to relationship is not taken care by the provider.
