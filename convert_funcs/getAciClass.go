package convert_funcs

type GetAciClass func(string) string

func AciClassMap(prefix string) string {
	mapping := map[string]string{
		"rsprotBy":            "FvRsProtBy",
		"rtmapentry":          "PimRouteMapEntry",
		"eptags":              "FvEpIpTag",
		"idgattr":             "FvIdGroupAttr",
		"ipattr":              "FvIpAttr",
		"macattr":             "FvMacAttr",
		"rscons":              "FvRsCons",
		"rsconsIf":            "FvRsConsIf",
		"trustctrlpol":        "FhsTrustCtrlPol",
		"eptags":              "FvEpMacTag",
		"nexthop":             "FvFBRMember",
		"qoscustom":           "QosCustomPol",
		"rspathAtt":           "FvRsPathAtt",
		"crtrn":               "FvSCrtrn",
		"rsoutToFBRGroup":     "L3extRsOutToFBRGroup",
		"rsooBCons":           "MgmtRsOoBCons",
		"monitorpol":          "NetflowMonitorPol",
		"rsdomAtt":            "FvRsDomAtt",
		"rsintraEpg":          "FvRsIntraEpg",
		"rtmap":               "PimRouteMapPol",
		"prof":                "RtctrlProfile",
		"nodesidp":            "MplsNodeSidP",
		"recordpol":           "NetflowRecordPol",
		"tn":                  "VzOOBBrCP",
		"rsfcPathAtt":         "FvRsFcPathAtt",
		"vmattr":              "FvVmAttr",
		"rsnodeAtt":           "FvRsNodeAtt",
		"annotationKey":       "TagAnnotation",
		"fbrg":                "FvFBRGroup",
		"pfx":                 "FvFBRoute",
		"tagKey":              "TagTag",
		"tn":                  "MgmtInstP",
		"rsmonitorToExporter": "NetflowRsMonitorToExporter",
		"qosdpppol":           "QosDppPol",
		"rsprov":              "FvRsProv",
		"provlbl":             "L3extProvLbl",
		"epg":                 "FvAEPg",
		"crtrn":               "FvCrtrn",
		"esg":                 "FvESg",
		"rsredistributePol":   "L3extRsRedistributePol",
		"subnet":              "MgmtSubnet",
		"dnsattr":             "FvDnsAttr",
		"rssecInherited":      "FvRsSecInherited",
		"conslbl":             "L3extConsLbl",
	}

	if class, found := mapping[prefix]; found {
		return class
	}
	return ""
}
