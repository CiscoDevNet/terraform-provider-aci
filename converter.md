# Cisco ACI Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 
  - v0.12 and higher (ACI Provider v1.0.0 or higher)
  - v0.11.x or below (ACI Provider v0.7.1 or below)

- [Go](https://golang.org/doc/install) Latest Version


### Pre-Requirements

1. Install latest version of [Go](http://www.golang.org)

### Setting up Terraform

1. A .tf configuration file is required to describe the providers, resources, and data sources that represent your infrastructure using the HashiCorp Configuration Language(HCL). 

The file should follow this format:

terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  # APIC Username
  username = "admin"
  # APIC Password
  password = "!v3G@!4@Y"
  # APIC URL
  url      = "https://sandboxapicdc.cisco.com"
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
    name        = "example_tenant"
    description = "This tenant is created by terraform"
}

resource "aci_annotation" "terraform_annotation" {
  parent_dn = "uni/tn-example_tenant"
  key       = "test_key"
  value     = "test_value"
}

### Set up directory for aci_converter.go

1. While in the directory with <main.tf>,  run the following commands:

terraform init
terraform plan -out = <bin_file_name.bin>
terraform show -json <bin_file_name.bin> > <json_file_name.json>

- In this instance, <json_file_name.json> will be used as the <INPUT_FILE.json> for aci_converter.go.


### Compiling 

1.  Run aci_converter.go  using the following command:

go run aci_converter.go <INPUT_FILE.json> <OUTPUT_FILE.json>

- The <OUTPUT_FILE.json> is the file the ACI Payload will be written to.
