package test

// RelationInfoFvRsCons represents a named relation to a contract.
var RelationInfoFvRsCons = map[string]interface{}{
	"label": "contract",
	"relationInfo": map[string]interface{}{
		"type":   "named",
		"fromMo": "fv:EPg",
		"toMo":   "vz:BrCP",
	},
}

// RelationInfoNetflowRsExporterToCtx represents an explicit relation from exporter to context.
var RelationInfoNetflowRsExporterToCtx = map[string]interface{}{
	"label": "netflow exporter",
	"relationInfo": map[string]interface{}{
		"type":   "explicit",
		"fromMo": "netflow:AExporterPol",
		"toMo":   "fv:Ctx",
	},
	"isCreatableDeletable": "always",
}

// GlobalMetaDefinitionVrf returns a NoMetaFile map with VRF mapping.
func GlobalMetaDefinitionVrf() map[string]string {
	return map[string]string{
		"fvCtx": "vrf",
	}
}

// GlobalMetaDefinitionContract returns a NoMetaFile map with contract mapping.
func GlobalMetaDefinitionContract() map[string]string {
	return map[string]string{
		"vzBrCP": "contract",
	}
}

// GlobalMetaDefinitionNetflow returns a NoMetaFile map with netflow and VRF mapping.
func GlobalMetaDefinitionNetflow() map[string]string {
	return map[string]string{
		"fvCtx":               "vrf",
		"netflowAExporterPol": "netflow_exporter_policy",
	}
}
