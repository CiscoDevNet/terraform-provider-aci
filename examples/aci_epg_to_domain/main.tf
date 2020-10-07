#configure provider with your cisco aci credentials.
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

# provider "aci" {
#   username    = ""
#   private_key = ""
#   cert_name   = ""
#   url         = ""
#   insecure    = true
# }