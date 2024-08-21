package convert_funcs

func dn_to_class(prefix string) string {

	mapping := map[string]string{

		"tagAnnotation-rsHPathAtt": "infraRsHPathAtt",

		"tagTag-rsHPathAtt": "infraRsHPathAtt",

		"tagAnnotation-rscons": "fvRsCons",

		"tagTag-rscons": "fvRsCons",

		"mgmtRsOoBCons-instp": "mgmtInstP",

		"tagAnnotation-instp": "mgmtInstP",

		"tagTag-instp": "mgmtInstP",

		"tagAnnotation-rsmonitorToExporter": "netflowRsMonitorToExporter",

		"tagTag-rsmonitorToExporter": "netflowRsMonitorToExporter",

		"tagAnnotation-rtmap": "pimRouteMapPol",

		"tagTag-rtmap": "pimRouteMapPol",

		"tagAnnotation-crtrn": "fvCrtrn",

		"tagTag-crtrn": "fvCrtrn",

		"tagAnnotation-rsconsIf": "fvRsConsIf",

		"tagTag-rsconsIf": "fvRsConsIf",

		//	"tagAnnotation-crtrn": "fvSCrtrn",

		//	"tagTag-crtrn": "fvSCrtrn",

		"tagAnnotation-subnet": "mgmtSubnet",

		"tagTag-subnet": "mgmtSubnet",

		"tagAnnotation-qosdpppol": "qosDppPol",

		"tagTag-qosdpppol": "qosDppPol",

		"fvCrtrn-epg": "fvAEPg",

		"fvRsAEPgMonPol-epg": "fvAEPg",

		"fvRsBd-epg": "fvAEPg",

		"fvRsCons-epg": "fvAEPg",

		"fvRsConsIf-epg": "fvAEPg",

		"fvRsCustQosPol-epg": "fvAEPg",

		"fvRsDomAtt-epg": "fvAEPg",

		"fvRsDppPol-epg": "fvAEPg",

		"fvRsFcPathAtt-epg": "fvAEPg",

		"fvRsIntraEpg-epg": "fvAEPg",

		"fvRsNodeAtt-epg": "fvAEPg",

		"fvRsPathAtt-epg": "fvAEPg",

		"fvRsProtBy-epg": "fvAEPg",

		"fvRsProv-epg": "fvAEPg",

		"fvRsSecInherited-epg": "fvAEPg",

		"fvRsTrustCtrl-epg": "fvAEPg",

		"tagAnnotation-epg": "fvAEPg",

		"tagTag-epg": "fvAEPg",

		"tagAnnotation-dnsattr": "fvDnsAttr",

		"tagTag-dnsattr": "fvDnsAttr",

		"tagAnnotation-recordpol": "netflowRecordPol",

		"tagTag-recordpol": "netflowRecordPol",

		"fvFBRMember-fbrg": "fvFBRGroup",

		"fvFBRoute-fbrg": "fvFBRGroup",

		"tagAnnotation-fbrg": "fvFBRGroup",

		"tagTag-fbrg": "fvFBRGroup",

		"tagAnnotation-rsfcPathAtt": "fvRsFcPathAtt",

		"tagTag-rsfcPathAtt": "fvRsFcPathAtt",

		"tagAnnotation-pfx": "fvFBRoute",

		"tagTag-pfx": "fvFBRoute",

		"tagAnnotation-oobbrc": "vzOOBBrCP",

		"tagTag-oobbrc": "vzOOBBrCP",

		"tagAnnotation-qoscustom": "qosCustomPol",

		"tagTag-qoscustom": "qosCustomPol",

		"tagAnnotation-idgattr": "fvIdGroupAttr",

		"tagTag-idgattr": "fvIdGroupAttr",

		"tagAnnotation-rspathAtt": "fvRsPathAtt",

		"tagTag-rspathAtt": "fvRsPathAtt",

		"tagAnnotation-provlbl": "l3extProvLbl",

		"tagTag-provlbl": "l3extProvLbl",

		"tagAnnotation-trustctrlpol": "fhsTrustCtrlPol",

		"tagTag-trustctrlpol": "fhsTrustCtrlPol",

		"tagAnnotation-epiptag": "fvEpIpTag",

		"tagTag-epiptag": "fvEpIpTag",

		"tagAnnotation-epmactag": "fvEpMacTag",

		"tagTag-epmactag": "fvEpMacTag",

		"tagAnnotation-ipattr": "fvIpAttr",

		"tagTag-ipattr": "fvIpAttr",

		"tagAnnotation-nodesidp": "mplsNodeSidP",

		"tagTag-nodesidp": "mplsNodeSidP",

		"tagAnnotation-prof": "rtctrlProfile",

		"tagTag-prof": "rtctrlProfile",

		"tagAnnotation-rsprov": "fvRsProv",

		"tagTag-rsprov": "fvRsProv",

		"tagAnnotation-conslbl": "l3extConsLbl",

		"tagTag-conslbl": "l3extConsLbl",

		"netflowRsMonitorToExporter-monitorpol": "netflowMonitorPol",

		"netflowRsMonitorToRecord-monitorpol": "netflowMonitorPol",

		"tagAnnotation-monitorpol": "netflowMonitorPol",

		"tagTag-monitorpol": "netflowMonitorPol",

		"tagAnnotation-rsdomAtt": "fvRsDomAtt",

		"tagTag-rsdomAtt": "fvRsDomAtt",

		"tagAnnotation-rsooBCons": "mgmtRsOoBCons",

		"tagTag-rsooBCons": "mgmtRsOoBCons",

		"tagAnnotation-rsintraEpg": "fvRsIntraEpg",

		"tagTag-rsintraEpg": "fvRsIntraEpg",

		"tagAnnotation-rsprotBy": "fvRsProtBy",

		"tagTag-rsprotBy": "fvRsProtBy",

		"tagAnnotation-rsoutToFBRGroup": "l3extRsOutToFBRGroup",

		"tagTag-rsoutToFBRGroup": "l3extRsOutToFBRGroup",

		"fvRsCons-esg": "fvESg",

		"fvRsConsIf-esg": "fvESg",

		"fvRsIntraEpg-esg": "fvESg",

		"fvRsProv-esg": "fvESg",

		"fvRsScope-esg": "fvESg",

		"fvRsSecInherited-esg": "fvESg",

		"tagAnnotation-esg": "fvESg",

		"tagTag-esg": "fvESg",

		"tagAnnotation-macattr": "fvMacAttr",

		"tagTag-macattr": "fvMacAttr",

		"tagAnnotation-rssecInherited": "fvRsSecInherited",

		"tagTag-rssecInherited": "fvRsSecInherited",

		"tagAnnotation-rtmapentry": "pimRouteMapEntry",

		"tagTag-rtmapentry": "pimRouteMapEntry",

		"tagAnnotation-rsnodeAtt": "fvRsNodeAtt",

		"tagTag-rsnodeAtt": "fvRsNodeAtt",

		"tagAnnotation-nexthop": "fvFBRMember",

		"tagTag-nexthop": "fvFBRMember",

		"tagAnnotation-rsredistributePol": "l3extRsRedistributePol",

		"tagTag-rsredistributePol": "l3extRsRedistributePol",

		"tagAnnotation-vmattr": "fvVmAttr",

		"tagTag-vmattr": "fvVmAttr",

		"infraRsHPathAtt-hpaths": "infraHPathS",

		"infraRsPathToAccBaseGrp-hpaths": "infraHPathS",

		"tagAnnotation-hpaths": "infraHPathS",

		"tagTag-hpaths": "infraHPathS",
	}

	if className, found := mapping[prefix]; found {
		return className
	}

	return ""
}
