
resource "aci_key_ring" "full_example" {
  admin_state           = "completed"
  annotation            = "annotation"
  certificate           = <<EOT
-----BEGIN CERTIFICATE-----
MIIDUzCCAjugAwIBAgIUbIvQNyz27t/xrIW2dgMwlGWqJvowDQYJKoZIhvcNAQEL
BQAwODERMA8GA1UEAwwIVXNlciBBQkMxFjAUBgNVBAoMDUNpc2NvIFN5c3RlbXMx
CzAJBgNVBAYTAlVTMCAXDTI2MDMxNzEwMjY0MVoYDzIxMjYwMjIxMTAyNjQxWjA4
MREwDwYDVQQDDAhVc2VyIEFCQzEWMBQGA1UECgwNQ2lzY28gU3lzdGVtczELMAkG
A1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCeulLqFk1C
H6wKdn0/sNrpDOfkjxeIfsT388BV2FH9XWBYumk8iAIRLp2Tae0+dKOr50t0q9xX
pCmC0dT4vF6h7NYQIihCNE//LpijuGbw+iG8YsavCr8wUAKdkuwzdA/3MiXRACzE
WY3CYuhReGnSsjXUVvuaSb02jm80eCGDWUxm/KqH1qO/vlrzSU7NORvNrqQpKPA7
e4ukjGD7k6jO8c5+DWWC6aOAgD6Orpb58yK7jUPSlQsYHtoKmJb245lpE5yMWg8e
ZFfJIPh+cdP4gXcWVPSKJ/nMAj/vEhSp0WK6CJwu8+W/nhdvkemQR0oKb3XcJikV
ZBj3OO0HK1HlAgMBAAGjUzBRMB0GA1UdDgQWBBT9EPk1GYfyISGUsTX/+7mREkPV
djAfBgNVHSMEGDAWgBT9EPk1GYfyISGUsTX/+7mREkPVdjAPBgNVHRMBAf8EBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAzIrT1/ih7/BP363+90gRA29q4nA8RBWF6
JC1TVExRbgVtJOEbEUywTcUxeW1Q3bSSwW836qEWJpRcwvu103yZgn0yIWubndtb
pArujgGt5irDwQmwW9zVUPWrJkkUkNezEZzM9WrRh+fvepvlDQwq6Pqd75BpmiOr
ZPfliCH5wX8Hjlg3cFeQviRLEMjkLtaDmmDfTFRQYCjKohTSgZdd9aauzgsYOq/2
S2sHwyJR8MTUyxawvCRJP8q5v3KGLVUzU6vRtEbZUx5FvjNkgxLgsSP4+buQi7FR
aY0EpluqxPc3Wsl5AkIpFAf97vrRPGTPpSZig5LX9eS0Rmh4eJ10
-----END CERTIFICATE-----
EOT
  description           = "description_1"
  elliptic_curve        = "none"
  key                   = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCeulLqFk1CH6wK
dn0/sNrpDOfkjxeIfsT388BV2FH9XWBYumk8iAIRLp2Tae0+dKOr50t0q9xXpCmC
0dT4vF6h7NYQIihCNE//LpijuGbw+iG8YsavCr8wUAKdkuwzdA/3MiXRACzEWY3C
YuhReGnSsjXUVvuaSb02jm80eCGDWUxm/KqH1qO/vlrzSU7NORvNrqQpKPA7e4uk
jGD7k6jO8c5+DWWC6aOAgD6Orpb58yK7jUPSlQsYHtoKmJb245lpE5yMWg8eZFfJ
IPh+cdP4gXcWVPSKJ/nMAj/vEhSp0WK6CJwu8+W/nhdvkemQR0oKb3XcJikVZBj3
OO0HK1HlAgMBAAECggEAGRrupa2/Vy0RQK+7DF7ZpltvNhxXzt6reBhVRylUK9cg
DQihlScMmfvIhUS6IdMVrB6E4GrqX/D2y9qKbTPGDVBwudPdOujGml2cBjvStFyr
sf3hilVl0KtnA9X1hrHdxe29mNAGmfZxey8PqwsjmbqsKZipv56CzzqPZjogpYMz
Q0Vu9ek2Y8Piw4TCepixAkFHlYtzpfzxFJ1xapam8o+mOyi+5FwVbx2QnA31Y+6p
cjzXnXnWl3aRNTTfvhz6Z/cZmDJ+sd/ATpk4/jt8biR4mOaYNS4k0cZX0PpjKP3K
hoR9DPU5qRDmgoShF7TMkBRiM2oG2Nxd794unRZMJQKBgQDbfdx/ThXhLnCbl98B
hN2JfLLuGvgO0NvDjlViPwmwnjIotiOArn4MP79m6CS/nKgW6wzuXfIJQcxPUslC
coloK9TDSya5nD1tqFkukhPukc3Q5Sb41pZzhksTjZgG3dohwqqRvcmnZPpxssa5
Bqv64zhiPPHhLE9whLNJ0oPl2wKBgQC5IRaQcyLXLHFjd2JIn2KrWAf9T3StUeE1
03y/IgZtwC6Kt8MvLmy3LENbHQdISibY+4IKoe+oEI/jjvNwM9sBwbcqm3L20+mY
i6W5r6qYufc58+sVtn2PqL8DK47xrqsTYnYIO6XxgavVIxGFAlgm0WDS7/tyZFZ0
oJ63uvuTPwKBgE8BgOr6Cnohozr+cbE0SCIDFs0KPBvpJhHAHA/fLPe8GcX5udHJ
/WkfUSATk5a9JuwI84ChpEucuZQb8oHOhJNQo6cgV/IbwSjFnkRbJH2NUg5NTbfo
VgODZWbGYuCl3qykS41mST3N5TAj79AODL2kKFmEInSw5G9V9Msv0XZdAoGAEzp7
lH9Q9BZ6pIEm5TIg1nkrQ0U4cjQZ9zRDNbr7/fRDIUda75Cb3B6t1E3cjsac6Faf
OCl/se4ec91KLbJFIhaTxsokk7yI+74tdW7ogjp2kj9igHvW6M/3HwYsL7AbtsS4
S7yeTMpSJa4hyLXooAeJTf99F3GShUSVl7HFJZMCgYBrBTmJrHDhIAwIuzLwZfpr
R5jCxJmzpnEYjS7EBZqslTnB/Abr79SSezRPaR3H7OE/Ig4kSym0ekyjGPQbDV0T
+gpPD7GMWzFZ1RWQ6aJZ5bgwQACemsfYbWuV7bf9jx5Phle/KAjOzP8NgnqwSD3B
Hd067NWQRwOjTHgato7Yxg==
-----END PRIVATE KEY-----
EOT
  key_type              = "RSA"
  modulus               = "mod2048"
  name                  = "test_name"
  name_alias            = "name_alias_1"
  owner_key             = "owner_key_1"
  owner_tag             = "owner_tag_1"
  regenerate            = "no"
  certificate_authority = "test_name"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
}

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_key_ring" "full_example_tenant" {
  parent_dn             = aci_tenant.example.id
  admin_state           = "completed"
  annotation            = "annotation"
  certificate           = <<EOT
-----BEGIN CERTIFICATE-----
MIIDUzCCAjugAwIBAgIUbIvQNyz27t/xrIW2dgMwlGWqJvowDQYJKoZIhvcNAQEL
BQAwODERMA8GA1UEAwwIVXNlciBBQkMxFjAUBgNVBAoMDUNpc2NvIFN5c3RlbXMx
CzAJBgNVBAYTAlVTMCAXDTI2MDMxNzEwMjY0MVoYDzIxMjYwMjIxMTAyNjQxWjA4
MREwDwYDVQQDDAhVc2VyIEFCQzEWMBQGA1UECgwNQ2lzY28gU3lzdGVtczELMAkG
A1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCeulLqFk1C
H6wKdn0/sNrpDOfkjxeIfsT388BV2FH9XWBYumk8iAIRLp2Tae0+dKOr50t0q9xX
pCmC0dT4vF6h7NYQIihCNE//LpijuGbw+iG8YsavCr8wUAKdkuwzdA/3MiXRACzE
WY3CYuhReGnSsjXUVvuaSb02jm80eCGDWUxm/KqH1qO/vlrzSU7NORvNrqQpKPA7
e4ukjGD7k6jO8c5+DWWC6aOAgD6Orpb58yK7jUPSlQsYHtoKmJb245lpE5yMWg8e
ZFfJIPh+cdP4gXcWVPSKJ/nMAj/vEhSp0WK6CJwu8+W/nhdvkemQR0oKb3XcJikV
ZBj3OO0HK1HlAgMBAAGjUzBRMB0GA1UdDgQWBBT9EPk1GYfyISGUsTX/+7mREkPV
djAfBgNVHSMEGDAWgBT9EPk1GYfyISGUsTX/+7mREkPVdjAPBgNVHRMBAf8EBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAzIrT1/ih7/BP363+90gRA29q4nA8RBWF6
JC1TVExRbgVtJOEbEUywTcUxeW1Q3bSSwW836qEWJpRcwvu103yZgn0yIWubndtb
pArujgGt5irDwQmwW9zVUPWrJkkUkNezEZzM9WrRh+fvepvlDQwq6Pqd75BpmiOr
ZPfliCH5wX8Hjlg3cFeQviRLEMjkLtaDmmDfTFRQYCjKohTSgZdd9aauzgsYOq/2
S2sHwyJR8MTUyxawvCRJP8q5v3KGLVUzU6vRtEbZUx5FvjNkgxLgsSP4+buQi7FR
aY0EpluqxPc3Wsl5AkIpFAf97vrRPGTPpSZig5LX9eS0Rmh4eJ10
-----END CERTIFICATE-----
EOT
  description           = "description_1"
  elliptic_curve        = "none"
  key                   = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCeulLqFk1CH6wK
dn0/sNrpDOfkjxeIfsT388BV2FH9XWBYumk8iAIRLp2Tae0+dKOr50t0q9xXpCmC
0dT4vF6h7NYQIihCNE//LpijuGbw+iG8YsavCr8wUAKdkuwzdA/3MiXRACzEWY3C
YuhReGnSsjXUVvuaSb02jm80eCGDWUxm/KqH1qO/vlrzSU7NORvNrqQpKPA7e4uk
jGD7k6jO8c5+DWWC6aOAgD6Orpb58yK7jUPSlQsYHtoKmJb245lpE5yMWg8eZFfJ
IPh+cdP4gXcWVPSKJ/nMAj/vEhSp0WK6CJwu8+W/nhdvkemQR0oKb3XcJikVZBj3
OO0HK1HlAgMBAAECggEAGRrupa2/Vy0RQK+7DF7ZpltvNhxXzt6reBhVRylUK9cg
DQihlScMmfvIhUS6IdMVrB6E4GrqX/D2y9qKbTPGDVBwudPdOujGml2cBjvStFyr
sf3hilVl0KtnA9X1hrHdxe29mNAGmfZxey8PqwsjmbqsKZipv56CzzqPZjogpYMz
Q0Vu9ek2Y8Piw4TCepixAkFHlYtzpfzxFJ1xapam8o+mOyi+5FwVbx2QnA31Y+6p
cjzXnXnWl3aRNTTfvhz6Z/cZmDJ+sd/ATpk4/jt8biR4mOaYNS4k0cZX0PpjKP3K
hoR9DPU5qRDmgoShF7TMkBRiM2oG2Nxd794unRZMJQKBgQDbfdx/ThXhLnCbl98B
hN2JfLLuGvgO0NvDjlViPwmwnjIotiOArn4MP79m6CS/nKgW6wzuXfIJQcxPUslC
coloK9TDSya5nD1tqFkukhPukc3Q5Sb41pZzhksTjZgG3dohwqqRvcmnZPpxssa5
Bqv64zhiPPHhLE9whLNJ0oPl2wKBgQC5IRaQcyLXLHFjd2JIn2KrWAf9T3StUeE1
03y/IgZtwC6Kt8MvLmy3LENbHQdISibY+4IKoe+oEI/jjvNwM9sBwbcqm3L20+mY
i6W5r6qYufc58+sVtn2PqL8DK47xrqsTYnYIO6XxgavVIxGFAlgm0WDS7/tyZFZ0
oJ63uvuTPwKBgE8BgOr6Cnohozr+cbE0SCIDFs0KPBvpJhHAHA/fLPe8GcX5udHJ
/WkfUSATk5a9JuwI84ChpEucuZQb8oHOhJNQo6cgV/IbwSjFnkRbJH2NUg5NTbfo
VgODZWbGYuCl3qykS41mST3N5TAj79AODL2kKFmEInSw5G9V9Msv0XZdAoGAEzp7
lH9Q9BZ6pIEm5TIg1nkrQ0U4cjQZ9zRDNbr7/fRDIUda75Cb3B6t1E3cjsac6Faf
OCl/se4ec91KLbJFIhaTxsokk7yI+74tdW7ogjp2kj9igHvW6M/3HwYsL7AbtsS4
S7yeTMpSJa4hyLXooAeJTf99F3GShUSVl7HFJZMCgYBrBTmJrHDhIAwIuzLwZfpr
R5jCxJmzpnEYjS7EBZqslTnB/Abr79SSezRPaR3H7OE/Ig4kSym0ekyjGPQbDV0T
+gpPD7GMWzFZ1RWQ6aJZ5bgwQACemsfYbWuV7bf9jx5Phle/KAjOzP8NgnqwSD3B
Hd067NWQRwOjTHgato7Yxg==
-----END PRIVATE KEY-----
EOT
  key_type              = "RSA"
  modulus               = "mod2048"
  name                  = "test_name"
  name_alias            = "name_alias_1"
  owner_key             = "owner_key_1"
  owner_tag             = "owner_tag_1"
  regenerate            = "no"
  certificate_authority = "test_name"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
}
