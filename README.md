
# Cisco ACI Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 
  - v0.12 and higher (ACI Provider v1.0.0 or higher)
  - v0.11.x or below (ACI Provider v0.7.1 or below)

- The [Go toolchain declared in `go.mod`](https://go.dev/doc/toolchain)

## Building The Provider

Clone the repository into any local directory:

```sh
git clone https://github.com/CiscoDevNet/terraform-provider-aci.git
cd terraform-provider-aci
```

Download the module dependencies and build the provider:

```sh
go mod download
make build
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

resource "aci_tenant" "example" {
  name        = "test-tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "example" {
  tenant_dn   = aci_tenant.example.id
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

1. Install the [Go toolchain declared in `go.mod`](https://go.dev/doc/toolchain).

### Existing resources and data-sources developed with terraform-plugin-sdk/v2

* Existing resources and data-sources are located in the [aci](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/aci) directory.

* Changes are made directly in the `provider.go`, `resource_*.go`, `data_source_*.go` and `utils` files. The [aci-go-client](https://github.com/ciscoecosystem/aci-go-client) is leveraged to construct payload constructs and handle REST communication towards the APIC.

* Documentation is manually maintained in the [legacy-docs](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/legacy-docs) directory and in the corresponding Terraform Registry pages under [docs](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/docs). The current generator does not copy legacy documentation into `docs`.

* Examples outside `gen/managed_files.json` are maintained manually in the
  [examples](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/examples)
  directory. Manifest-owned examples are regenerated from normalized class
  data and active templates.

#### Manual Development Process

The below steps should be followed for developing `terraform-plugin-sdk/v2` resources and data-sources:

1. Create issue (if not created yet) and comment that you will be working on the issue.

2. Fork the terraform-provider-aci repository.

3. Clone the forked code to your local machine.

4. Make changes to the files manually.
    * Code changes
    * Examples changes
    * Documentation changes 

5. Keep the matching files under `legacy-docs` and `docs` synchronized manually.

6. If generator inputs were changed, run `go generate` from the repository root and review every generated change. Normal generation also updates owned files under [internal/provider](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider) and `examples`.

7. Test the code.

8. Create PR for the code and request review from active maintainers.

9. Review process

### Resources and data-sources developed with terraform-plugin-framework

The terraform-plugin-framework implementation is located in
[internal/provider](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider).
Reusable generated ACI class models are located in
[`internal/provider/models`](https://github.com/CiscoDevNet/terraform-provider-aci/tree/main/internal/provider/models).

The generator is being migrated incrementally. A generated marker does not by
itself mean that the current generator owns a file: legacy generated files keep
their marker until their templates have been replaced. The exact files owned by
normal generation are recorded in [`gen/managed_files.json`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/main/gen/managed_files.json).
Currently, normal generation owns the class model files, provider index, and
the provider, minimum-resource, all-attributes-resource, and data-source
example files and associated supporting files recorded in the manifest. Resource
and data-source implementations, tests, and their Markdown documentation remain
static or legacy outputs until their replacement templates are introduced.

Generated files should not be edited manually. Fix their normalized metadata,
class definition, generator logic, or active template and regenerate the full
output set. Provider bootstrap and registry files and
`examples/provider/provider.tf` are intentionally static. The provider index
is rendered from its active template using that static example. The complete
current ownership contract is documented in
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
- `gen/definitions/example_files/` contains shared files copied beside generated
  resource examples when requested by a class definition.
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

Set `GEN_ACI_TF_META_REFRESH_ALL=1` to refresh every existing metadata file
before generation. Generator logging defaults to `INFO` on standard output.
`GEN_ACI_TF_LOG_LEVEL` accepts `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, or
`FATAL`; `GEN_ACI_TF_LOG_PATH` redirects output to a file opened in append
mode. For example:

```sh
GEN_ACI_TF_LOG_LEVEL=DEBUG GEN_ACI_TF_LOG_PATH=generator.log go generate
```

`internal/provider/annotation_unsupported.go` is a refresh-only generated file.
Refresh it explicitly from the complete ACI metadata:

```sh
GEN_ANNOTATION_UNSUPPORTED=1 go generate
```

When changing a definition produced by
`gen/scripts/migrate_class_definitions.go`, update the migration logic and
regenerate the definition instead of changing either its emitted YAML or the
legacy inputs under `gen/definitions/classes` and `gen/definitions/properties`:

```sh
go run gen/scripts/migrate_class_definitions.go
go generate
```

`gen/definitions/global.yaml` is the canonical source for shared configuration
and is intentionally not rewritten by the migration script. The legacy input
directories retain their `develop` layout so upstream changes and additions are
discovered automatically. After generation, review all generated changes and
run the relevant unit and acceptance tests before opening a pull request.

### Troubleshooting

#### Go Generate Fail

Run the generator with the Go toolchain declared in `go.mod`:

```sh
go generate
```

The repository commits its vendor directory. When a dependency is intentionally
added, removed, or changed, synchronize both the module files and vendor tree:

```sh
go mod tidy
go mod vendor
```

Review changes to `go.mod`, `go.sum`, and `vendor/` before committing them. Do
not use `go get` to repair an unchanged checkout; a missing vendored package in
that case indicates repository drift or an incomplete checkout.


### Compiling

To compile the provider, run `make build`. This runs the formatting check and
installs the provider binary into `GOBIN`, or into `GOPATH/bin` when `GOBIN` is
not set.

<strong>Important: </strong>To successfully use the provider you need to follow these steps:

- Copy or symlink the provider from `GOBIN` (or from `GOPATH/bin` when `GOBIN` is unset) to `~/.terraform.d/plugins/terraform.local/CiscoDevNet/aci/<Version>/<architecture>/`, for example:
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
