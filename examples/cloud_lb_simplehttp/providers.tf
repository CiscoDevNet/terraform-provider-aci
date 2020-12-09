
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.12.0"
    }
    aci = {
      source  = "CiscoDevNet/aci"
      version = "0.5.0"
    }
  }

}
