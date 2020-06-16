provider "aci" {
  username = "admin"
  password = "cisco123"
  url      = "https://192.168.10.102"
  insecure = true
}

# provider "aci" {
#   username    = "admin"
#   private_key = <<EOF

# EOF
#   # private_key = "/Users/nirav.katarmal/Documents/github/aci_test/admin.key"
#   cert_name   = "admin"
#   url         = "https://192.168.10.102"
#   insecure    = true
# }