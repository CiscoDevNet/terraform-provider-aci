package data

const (
	// Environment variables for logging configuration.
	// Determine if the meta data in the meta directory should be refreshed from remote location
	// The default is to not refresh the meta data.
	envMetaRefresh = "GEN_ACI_TF_META_REFRESH_ALL"
	// The host from which the meta data is retrieved, when not set, the default is the DevNet base URL.
	envMetaHost = "GEN_ACI_TF_META_HOST"
	// The environment variable that contains the list of classes to be retrieved.
	// This can be used to retrieve the meta data for classes that are not available in the meta directory.
	// This can be used to retrieve the meta data for classes that are available in the meta directory and require a refresh.
	// The classes are separated by a comma, ex "fvTenant,fvAp".
	// When GEN_ACI_TF_META_CLASSES is not set, no meta data is retrieved.
	envMetaClasses   = "GEN_ACI_TF_META_CLASSES"
	pubhubDevnetHost = "pubhub.devnetcloud.com/media/model-doc-latest/docs"
	metaFileUrl      = "https://%s/doc/jsonmeta/%s/%s.json"
	// The path to the meta data files.
	metaPath = "./gen/meta"
)
