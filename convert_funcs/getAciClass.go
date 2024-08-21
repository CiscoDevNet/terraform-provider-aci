package convert_funcs

func GetAciClass(prefix string) string {
	mapping := map[string]string{
		"rsHPathAtt":          "infraRsHPathAtt",
		"rscons":              "fvRsCons",
		"instp":               "mgmtInstP",
		"rsmonitorToExporter": "netflowRsMonitorToExporter",
		"rtmap":               "pimRouteMapPol",
		"crtrn":               "fvCrtrn",
		"rsconsIf":            "fvRsConsIf",
		//"crtrn":               "fvSCrtrn",
		"subnet":            "mgmtSubnet",
		"qosdpppol":         "qosDppPol",
		"epg":               "fvAEPg",
		"dnsattr":           "fvDnsAttr",
		"annotationKey":     "tagAnnotation",
		"recordpol":         "netflowRecordPol",
		"fbrg":              "fvFBRGroup",
		"rsfcPathAtt":       "fvRsFcPathAtt",
		"pfx":               "fvFBRoute",
		"oobbrc":            "vzOOBBrCP",
		"qoscustom":         "qosCustomPol",
		"idgattr":           "fvIdGroupAttr",
		"rspathAtt":         "fvRsPathAtt",
		"provlbl":           "l3extProvLbl",
		"trustctrlpol":      "fhsTrustCtrlPol",
		"epiptag":           "fvEpIpTag",
		"epmactag":          "fvEpMacTag",
		"ipattr":            "fvIpAttr",
		"nodesidp":          "mplsNodeSidP",
		"prof":              "rtctrlProfile",
		"rsprov":            "fvRsProv",
		"conslbl":           "l3extConsLbl",
		"monitorpol":        "netflowMonitorPol",
		"rsdomAtt":          "fvRsDomAtt",
		"rsooBCons":         "mgmtRsOoBCons",
		"rsintraEpg":        "fvRsIntraEpg",
		"rsprotBy":          "fvRsProtBy",
		"rsoutToFBRGroup":   "l3extRsOutToFBRGroup",
		"esg":               "fvESg",
		"macattr":           "fvMacAttr",
		"rssecInherited":    "fvRsSecInherited",
		"rtmapentry":        "pimRouteMapEntry",
		"rsnodeAtt":         "fvRsNodeAtt",
		"tagKey":            "tagTag",
		"nexthop":           "fvFBRMember",
		"rsredistributePol": "l3extRsRedistributePol",
		"vmattr":            "fvVmAttr",
		"hpaths":            "infraHPathS",
	}

	if class, found := mapping[prefix]; found {
		return class
	}
	return ""
}
