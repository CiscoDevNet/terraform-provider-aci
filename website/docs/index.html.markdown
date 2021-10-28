---
layout: "aci"
page_title: "Provider: ACI"
sidebar_current: "docs-aci-index"
description: |-
  The Cisco ACI provider is used to interact with the resources provided by Cisco APIC.
  The provider needs to be configured with the proper credentials before it can be used.
---

Application Centric Infrastructure (ACI)
-----------------------------------------  
The Cisco Application Centric Infrastructure (ACI) allows application requirements to define the network. This architecture simplifies, optimizes, and accelerates the entire application deployment life cycle.  

Application Policy Infrastructure Controller (APIC)
--------------------------------------------------
The APIC manages the scalable ACI multi-tenant fabric. The APIC provides a unified point of automation and management, policy programming, application deployment, and health monitoring for the fabric. The APIC, which is implemented as a replicated synchronized clustered controller, optimizes performance, supports any application anywhere, and provides unified operation of the physical and virtual infrastructure.
The APIC enables network administrators to easily define the optimal network for applications. Data center operators can clearly see how applications consume network resources, easily isolate and troubleshoot application and infrastructure problems, and monitor and profile resource usage patterns.
The Cisco Application Policy Infrastructure Controller (APIC) API enables applications to directly connect with a secure, shared, high-performance resource pool that includes network, compute, and storage capabilities.

Cisco ACI Provider
------------
The Cisco ACI terraform provider is used to interact with resources provided by Cisco APIC. The provider needs to be configured with proper credentials to authenticate with Cisco APIC.

Authentication
--------------
The Provider supports authentication with Cisco APIC in 3 ways:  

 1. Authentication with user-id and password.  
 example:  

 ```hcl
    provider "aci" {
      # cisco-aci user name
      username = "admin"
      # cisco-aci password
      password = "password"
      # cisco-aci url
      url      = "https://my-cisco-aci.com"
      insecure = true
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

  3. Authentication with login-domain and password.
  example:

  ```hcl
    provider "aci" {
      username = "apic:Demo_domain\\\\admin"
      # private_key = "path to private key"
      # cert_name = "user-cert"
      password = "password"
      url = "url"
      insecure = true
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

Example Usage
------------
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
}

resource "aci_tenant" "test-tenant" {
  name        = "test-tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "test-app" {
  tenant_dn   = "${aci_tenant.test-tenant.id}"
  name        = "test-app"
  description = "This app profile is created by terraform"
}
```

Argument Reference
------------------
Following arguments are supported with Cisco ACI terraform provider.

 * `username` - (Required) This is the Cisco APIC username, which is required to authenticate with CISCO APIC.
 * `password` - (Optional) Password of the user mentioned in username argument. It is required when you want to use token-based authentication.
 * `private-key` - (Optional) Path to the private key for which x509 certificate has been calculated for the user mentioned in `username`.
 * `url` - (Required) URL for CISCO APIC.
 * `insecure` - (Optional) This determines whether to use insecure HTTP connection or not. Default value is `true`.  
 * `validate_relation_dn` - (Optional) Flag to validate if a object with added relation Dn exists in the APIC. Type: Bool, Default: "true". 
 * `cert_name` - (Optional) Certificate name for the User in Cisco ACI.
 * `proxy_url` - (Optional) Proxy Server URL with port number.
 * `proxy_creds` - (Optional) Proxy server credentials in the form of username:password.


~> NOTE: `password` or `private-key` either of one is required.
