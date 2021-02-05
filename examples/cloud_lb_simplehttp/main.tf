resource "aci_rest" "cloudALB" {
  path    = "/api/node/mo/${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/clb-${var.name}.json"
  payload = <<EOF
  cloudLB:
        attributes:
          scheme: internet
          type: application
        children:
        - cloudRsLDevToCloudSubnet:
            attributes:
              tDn: ${var.subnet_a_dn}
        - cloudRsLDevToCloudSubnet:
            attributes:
              tDn: ${var.subnet_b_dn}
  EOF
}

resource "aci_rest" "serviceGraphCreate" {
  path       = "/api/node/mo/${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/AbsGraph-${var.name}.json"
  depends_on = [aci_rest.cloudALB]
  payload    = <<EOF
  vnsAbsGraph:
        attributes:
          name: ${var.name}
          type: cloud
        children:
        - vnsAbsTermNodeProv:
            attributes:
              name: T2
            children:
            - vnsAbsTermConn:
                attributes:
                  name: ProvTermConn
        - vnsAbsTermNodeCon:
            attributes:
              name: T1
            children:
            - vnsAbsTermConn:
                attributes:
                  name: ConsTermConn
        - vnsAbsNode:
            attributes:
              name: N0
              managed: "yes"
              funcType: GoTo
              funcTemplateType: ADC_ONE_ARM
            children:
            - vnsRsNodeToCloudLDev: 
                attributes:
                  tDn: ${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/clb-${var.name}
            - vnsAbsFuncConn:
                attributes:
                  name: provider
                  attNotify: "no"
                  connType: none
            - vnsAbsFuncConn:
                attributes:
                  name: consumer
                  attNotify: "no"
                  connType: none
        - vnsAbsConnection:
            attributes:
              connDir: provider
              connType: external
              name: CON1
              adjType: L3
            children:
            - vnsRsAbsConnectionConns:
                attributes:
                  tDn: ${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/AbsGraph-${var.name}/AbsNode-N0/AbsFConn-provider
            - vnsRsAbsConnectionConns:
                attributes:
                  tDn: ${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/AbsGraph-${var.name}/AbsTermNodeProv-T2/AbsTConn
        - vnsAbsConnection:
            attributes:
              connDir: provider
              connType: external
              name: CON0
            children:
            - vnsRsAbsConnectionConns:
                attributes:
                  tDn: ${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/AbsGraph-${var.name}/AbsNode-N0/AbsFConn-consumer
            - vnsRsAbsConnectionConns:
                attributes:
                  tDn: ${join("", regex("(uni/tn-.*?)/", var.epg_dn))}/AbsGraph-${var.name}/AbsTermNodeCon-T1/AbsTConn
  EOF
}

resource "aci_rest" "appliedServiceGraph" {
  path       = "/api/node/mo/${join("", regex("(uni/tn-.*?)/", var.epg_dn))}.json"
  depends_on = [aci_rest.serviceGraphCreate]
  payload    = <<EOF
  fvTenant:
          attributes:
            name: ${join("", regex("/tn-(.*?)/", var.epg_dn))}
          children:
          - vnsAbsGraph:
              attributes:
                name: var.name
              children:
              - vnsAbsNode:
                  attributes:
                    name: N0
                  children:
                  - cloudSvcPolicy:
                      attributes:
                        contractName: ${join("", regex("/brc-(.*?)/", var.contract_subject_dn))}
                        subjectName: ${join("", regex("/subj-(.*)", var.contract_subject_dn))}
                        tenantName: ${join("", regex("/tn-(.*?)/", var.epg_dn))}
                      children: 
                      - cloudListener:
                          attributes:
                            name: http_listener
                            port: var.listenerPort
                            protocol: http
                          children:
                          - cloudListenerRule:
                              attributes:
                                default: "yes"
                                name: forward
                                priority: "999"
                              children:
                              - cloudRuleAction:
                                  attributes:
                                    epgdn: var.epg_dn
                                    port: var.hostPort
                                    protocol: http
                                    type: forward
          - vzBrCP:
              attributes:
                name: ${join("", regex("/brc-(.*?)/", var.contract_subject_dn))}
              children:
              - vzSubj:
                 attributes:
                   name: ${join("", regex("/subj-(.*)", var.contract_subject_dn))}
                 children:
                 - vzRsSubjGraphAtt:
                     attributes:
                       tnVnsAbsGraphName: var.name
  EOF
}

output "dnsname" {
  value = data.aws_lb.wwwalb.dns_name
}
