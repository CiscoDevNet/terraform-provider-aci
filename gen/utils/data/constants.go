package data

const (
	// Environment variables for logging configuration.
	// The path to the definition files.
	constDefinitionsPath = "../../../gen/definitions"
	// The path to the global definition file.
	constGlobalDefinitionFilePath = "../../../gen/definitions/global.yaml"
	// Determine if the meta data in the meta directory should be refreshed from remote location
	// The default is to not refresh the meta data.
	constEnvMetaRefresh = "GEN_ACI_TF_META_REFRESH_ALL"
	// The host from which the meta data is retrieved, when not set, the default is the DevNet base URL.
	constEnvMetaHost = "GEN_ACI_TF_META_HOST"
	// The environment variable that contains the list of classes to be retrieved.
	// This can be used to retrieve the meta data for classes that are not available in the meta directory.
	// This can be used to retrieve the meta data for classes that are available in the meta directory and require a refresh.
	// The classes are separated by a comma, ex "fvTenant,fvAp".
	// When GEN_ACI_TF_META_CLASSES is not set, no meta data is retrieved.
	constEnvMetaClasses   = "GEN_ACI_TF_META_CLASSES"
	constPubhubDevnetHost = "pubhub.devnetcloud.com/media/model-doc-latest/docs"
	constMetaFileUrl      = "https://%s/doc/jsonmeta/%s/%s.json"
	// The path to the meta data files.
	constMetaPath = "./gen/meta"
)
