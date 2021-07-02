resource "aci_x509_certificate" "test_cert" {
    local_user_dn = "uni/userext/user-admin"
    description = "From Terraform"
      
    name = "test_1"
    annotation  = "x509_certificate_tag"    
    name_alias  = "alias_name"
    data =<<EOF
-----BEGIN CERTIFICATE-----
MIIB4TCCAUoCCQCUp2SLleZjbjANBgkqhkiG9w0BAQsFADA0MQ4wDAYDVQQDDAVB
ZG1pbjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzAgFw0yMDAz
MzAxMzQ1MDRaGA8yMTIwMDMwNjEzNDUwNFowNDEOMAwGA1UEAwwFQWRtaW4xFTAT
BgNVBAoMDFlvdXIgQ29tcGFueTELMAkGA1UEBhMCVVMwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBALKlRyT0/Sx7bw+G79h+VHArPL5A+ZONgHrEu/G8US9dvMEc
HZXZ8ctC6yB8ILFkR2bierj3X+AyOu+247ne/hvXhzXQnWfk90167TewsFtF+RCg
SRM3nIQVZxwqCzXX4FqdH58e5Bi3DJudFLVB7pWWSjPfOxdCwDcSU3QssOJPAgMB
AAEwDQYJKoZIhvcNAQELBQADgYEAArlbmTrHMxRXxWAT0z2OBCjMSwy60Ef2YK16
P9ItPy3dddfIzKEVlC89Rws2nzbYYbohrSy5EMruE7e5fidjMYJLTvmyFsnv8JUt
xBF3y0UmaVfGYa0M5mWSqKAD7tSoK5QJjq+Me+sNP0Rdoj5QBDlttAnhLide8FE6
Q1R8g1Y=
-----END CERTIFICATE-----
EOF
}

terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}