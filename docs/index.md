---
layout: "aci"
page_title: "Provider: ACI"
sidebar_current: "docs-aci-index"
description: |-
  The Cisco ACI provider is used to interact with the resources provided by Cisco APIC.
  The provider needs to be configured with the proper credentials before it can be used.
---

# Application Centric Infrastructure (ACI)

The Cisco Application Centric Infrastructure (ACI) allows application requirements to define the network. This architecture simplifies, optimizes, and accelerates the entire application deployment life cycle.

# Application Policy Infrastructure Controller (APIC)

The APIC manages the scalable ACI multi-tenant fabric. The APIC provides a unified point of automation and management, policy programming, application deployment, and health monitoring for the fabric. The APIC, which is implemented as a replicated synchronized clustered controller, optimizes performance, supports any application anywhere, and provides unified operation of the physical and virtual infrastructure.
The APIC enables network administrators to easily define the optimal network for applications. Data center operators can clearly see how applications consume network resources, easily isolate and troubleshoot application and infrastructure problems, and monitor and profile resource usage patterns.
The Cisco Application Policy Infrastructure Controller (APIC) API enables applications to directly connect with a secure, shared, high-performance resource pool that includes network, compute, and storage capabilities.

# Cisco ACI Provider

The Cisco ACI terraform provider is used to interact with resources provided by Cisco APIC. The provider needs to be configured with proper credentials to authenticate with Cisco APIC.

## Authentication

The Provider supports authentication with Cisco APIC in 3 ways:

 1. Authentication with user-id and password.
 example:

 ```hcl
provider "aci" {
  username = "admin"
  password = "password"
  url      = "https://my-cisco-aci.com"
}
 ```

 In this method, it will obtain an authentication token from Cisco APIC and will use that token to authenticate. A limitation with this approach is APIC counts the request to authenticate and threshold it to avoid DOS attack. After too many attempts this authentication method may fail as the threshold will be exceeded.
 To avoid the above-mentioned problem Cisco APIC supports signature-based authentication.

 2. Signature Based authentication.
    * x509 certificate has been created and added it to the user in Cisco APIC.
    * With the help of private key that has been used to calculate the certificate, a signature has been calculated and passed with the request. This signature will be used to authenticate the user.
    example.

```
provider "aci" {
  username = "admin"
  private_key = "path to private key"
  cert_name = "user-cert"
  url      = "https://my-cisco-aci.com"
}
```

3. Authentication with login-domain and password.
   example:

```hcl
provider "aci" {
  username = "apic:Demo_domain\\\\admin"
  password = "password"
  url = "url"
}
```

### How to add Certificate to the Cisco APIC local user ###

* Generate certificate via below command.

```shell
$ openssl req -new -newkey rsa:1024 -days 36500 -nodes -x509 -keyout admin.key -out admin.crt -subj '/CN=Admin/O=Your Company/C=US'
```

* Add the X.509 certificate to your ACI AAA local user at ADMIN » AAA.

* Click AAA Authentication. Check that in the Authentication field the Realm field displays Local.

* Expand Security Management » Local Users
Click the name of the user you want to add a certificate to, in the User Certificates area
Click the + sign and in the Create X509 Certificate enter a certificate name in the Name field. Copy and paste your X.509 certificate in the Data field.
Use this certificate name as the value of the "cert_name" argument.

## Example Usage

```hcl
terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "password"
  url      = "https://my-cisco-aci.com"
  insecure = true
}

resource "aci_tenant" "example" {
  name = "example_tenant"
}
```

## Schema

NOTE: either 'password' OR 'private_key' and 'cert_name' must be provided for the ACI provider

### Optional

- `cert_name` (String) Certificate name for the User in Cisco ACI.
  - Environment variable: `ACI_CERT_NAME`
- `insecure` (Boolean) Allow insecure HTTPS client.
  - Default: `true`
  - Environment variable: `ACI_INSECURE`
- `password` (String) Password for the APIC Account.
  - Environment variable: `ACI_PASSWORD`
- `private_key` (String) Private key path for signature calculation.
  - Environment variable: `ACI_PRIVATE_KEY`
- `proxy_creds` (String) Proxy server credentials in the form of username:password.
  - Environment variable: `ACI_PROXY_CREDS`
- `proxy_url` (String) Proxy Server URL with port number.
  - Environment variable: `ACI_PROXY_URL`
- `retries` (Number) Number of retries for REST API calls.
  - Default: `2`
  - Environment variable: `ACI_RETRIES`
- `url` (String) URL of the Cisco ACI web interface. This can also be set as the ACI_URL environment variable.
  - Environment variable: `ACI_URL`
- `username` (String) Username for the APIC Account.
  - Environment variable: `ACI_USERNAME`
- `validate_relation_dn` (Boolean) Flag to validate if a object with entered relation Dn exists in the APIC.
  - Default: `true`
  - Environment variable: `ACI_VAL_REL_DN`
- `annotation` (String) Global annotation for the provider. 
  - Environment variable: `ACI_ANNOTATION`
