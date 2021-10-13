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

resource "aci_saml_provider" "example" {
    name                      = "example"
    name_alias                = "saml_provider_alias"
    description               = "From Terraform"
    annotation                = "orchestrator:terraform"
    entity_id                 = "entity_id_example" 
    gui_banner_message        = "gui_banner_message_example"
    https_proxy               = "https_proxy_example"
    id_p                      = "adfs"
    key                       = "key_example"
    metadata_url              = "metadata_url_example"
    monitor_server            = "disabled"
    monitoring_password       = "monitoring_password_example"
    monitoring_user           = "default"
    retries                   = "1"
    sig_alg                   = "SIG_RSA_SHA256"
    timeout                   = "5"
    tp                        = "tp_example"
    want_assertions_encrypted = "yes"
    want_assertions_signed    = "yes"
    want_requests_signed      = "yes"
    want_response_signed      = "yes"
    aaa_rs_prov_to_epp        = aci_resource.example.id
    aaa_rs_sec_prov_to_epg    = aci_resource.example.id
}