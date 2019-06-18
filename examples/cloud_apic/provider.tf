provider "aci" {
  username = "" <APIC username>
  password = "" <APIC pwd>
  url      = "" <cloud APIC URL>
  insecure = true
}

provider "aws" {
  region     = "us-west-1"
  access_key = ""
  secret_key = ""
}
