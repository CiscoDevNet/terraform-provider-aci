package convert_funcs

type GetAciClass func(string) string

func AciClassMap(prefix string) string {
	mapping := map[string]string{
		"hpaths":              "infraHPathS",
		"macattr":             "fvMacAttr",
		"rscons":              "fvRsCons",
		"trustctrlpol":        "fhsTrustCtrlPol",
		"conslbl":             "l3extConsLbl",
		"qoscustom":           "qosCustomPol",
		"ipattr":              "fvIpAttr",
		"recordpol":           "netflowRecordPol",
		"rsintraEpg":          "fvRsIntraEpg",
		"rspathAtt":           "fvRsPathAtt",
		"rsmonitorToExporter": "netflowRsMonitorToExporter",
		"rsprotBy":            "fvRsProtBy",
		"vmattr":              "fvVmAttr",
		"subnet":              "mgmtSubnet",
		"epmactag":            "fvEpMacTag",
		"fbrg":                "fvFBRGroup",
		"idgattr":             "fvIdGroupAttr",
		"crtrn":               "fvCrtrn",
		"nexthop":             "fvFBRMember",
		"rsconsIf":            "fvRsConsIf",
		"monitorpol":          "netflowMonitorPol",
		"epg":                 "fvAEPg",
		"esg":                 "fvESg",
		"rssecInherited":      "fvRsSecInherited",
		"rtmapentry":          "pimRouteMapEntry",
		"tagKey":              "tagTag",
		"rsoutToFBRGroup":     "l3extRsOutToFBRGroup",
		//	"crtrn":               "fvSCrtrn",
		"rsredistributePol": "l3extRsRedistributePol",
		"instp":             "mgmtInstP",
		"annotationKey":     "tagAnnotation",
		"rsHPathAtt":        "infraRsHPathAtt",
		"qosdpppol":         "qosDppPol",
		"oobbrc":            "vzOOBBrCP",
		"pfx":               "fvFBRoute",
		"rsfcPathAtt":       "fvRsFcPathAtt",
		"rsprov":            "fvRsProv",
		"rsooBCons":         "mgmtRsOoBCons",
		"nodesidp":          "mplsNodeSidP",
		"dnsattr":           "fvDnsAttr",
		"prof":              "rtctrlProfile",
		"rsdomAtt":          "fvRsDomAtt",
		"rsnodeAtt":         "fvRsNodeAtt",
		"provlbl":           "l3extProvLbl",
		"epiptag":           "fvEpIpTag",
		"rtmap":             "pimRouteMapPol",
	}

	if class, found := mapping[prefix]; found {
		return class
	}
	return ""
}
