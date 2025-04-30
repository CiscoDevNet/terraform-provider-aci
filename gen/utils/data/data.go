package data

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitalizeLogger()

type DataStore struct {
	Classes  map[string]Class
	metaHost string
}

func NewDataStore() *DataStore {
	dataStore := &DataStore{Classes: make(map[string]Class)}
	// Set the meta data host for retrieval of meta files.
	dataStore.setMetaHost()
	// Check if classes are set in the environment variable 'GEN_ACI_TF_META_CLASSES' and retrieve the meta files for those classes.
	dataStore.retrieveEnvClassesMetaFromRemote()
	// Initialize classes in data store for all meta files in the metaPath directory.
	dataStore.initializeMeta()

	return dataStore
}

func (ds *DataStore) setMetaHost() {
	// Check if the meta data host is set in the environment variables.
	// If it is set, use it as the host for the meta data.
	// If it is not set, use the default devnet host.
	host := os.Getenv(envMetaHost)
	if host == "" {
		host = pubhubDevnetHost
	}
	ds.metaHost = host
	genLogger.Info(fmt.Sprintf("Meta data host set to: %s.", host))
}

func (ds *DataStore) retrieveEnvClassesMetaFromRemote() {
	// Retrieve the meta data for the classes specified in the environment 'GEN_ACI_TF_META_CLASSES' variable
	classNames := strings.Split(os.Getenv(envMetaClasses), ",")
	if classNames[0] != "" {
		genLogger.Info(fmt.Sprintf("Retrieving meta files for classes: %s.", classNames))
		ds.setClasses(classNames, true)
	}
}

func (ds *DataStore) initializeMeta() {
	var refresh bool
	var err error
	// Check if the meta data should be refreshed from remote location as specified in the environment variable.
	metaRefreshEnv := os.Getenv(envMetaRefresh)
	// If the environment variable is not set, the default is to not refresh the meta data.
	if metaRefreshEnv != "" {
		// If the environment variable is set, parse it to a boolean value.
		// If the parsing fails, log a warning and skip retrieval of meta data.
		refresh, err = strconv.ParseBool(metaRefreshEnv)
		if err != nil {
			genLogger.Warn(fmt.Sprintf("Refreshing of meta is skipped due to error: %s.", err.Error()))
		}
	}
	ds.setClasses(utils.GetFileNamesFromDirectory(metaPath), refresh)
}

func (ds *DataStore) setClasses(classNames []string, retrieve bool) {
	// Retrieve the meta data for the classes provided.
	for _, className := range classNames {
		// Remove the ".json" suffix from the class name if it exists, which is the case when the file name is passed during refresh.
		if strings.HasSuffix(className, ".json") {
			className = strings.Replace(className, ".json", "", 1)
		}

		// When the class is already set, skip the process because it was retrieved during retrieveEnvClassesFromRemote.
		if _, ok := ds.Classes[className]; ok {
			genLogger.Trace(fmt.Sprintf("Meta data for class '%s' already set. Skipping...", className))
			continue
		}

		// Create a new class object and add it to the data store.
		// It is set to the datastore here to avoid retrieving the meta data again.
		// The creation of the class also sets the class name and package name which are required for the meta file download.
		classDetails := *NewClass(className)
		ds.Classes[className] = classDetails

		if retrieve {
			// ENHANCEMENT: Concurrently retrieve/write the meta data.
			// Retrieve the meta file for the class from the remote location.
			ds.retrieveMetaFileFromRemote(classDetails)
		}
	}
}

func (ds *DataStore) retrieveMetaFileFromRemote(classDetails Class) {
	url := fmt.Sprintf(metaFileUrl, ds.metaHost, classDetails.ClassNamePackage, classDetails.ClassNameShort)
	genLogger.Trace(fmt.Sprintf("Retrieving meta data for class '%s' from: %s.", classDetails.ClassName, url))

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Get(url)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during retrieval of meta file for class '%s': %s.", classDetails.ClassName, err.Error()))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during reading of file for class '%s': %s.", classDetails.ClassName, err.Error()))
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%s.json", metaPath, classDetails.ClassName))
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during creation of file for class '%s': %s.", classDetails.ClassName, err.Error()))
	}

	defer outputFile.Close()
	_, err = outputFile.Write(resBody)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during writing to file for class '%s': %s.", classDetails.ClassName, err.Error()))
	}

	genLogger.Trace(fmt.Sprintf("Succesfully wrote meta data for class '%s' to: %s.", classDetails.ClassName, outputFile.Name()))
}

func (ds *DataStore) LoadMetaFiles() {
	for className, classDetails := range ds.Classes {
		genLogger.Trace(fmt.Sprintf("Loading meta data for class '%s' to data store.", className))
		// ENHANCEMENT: Concurrently load the meta data.
		classDetails.LoadMetaFile()
	}
}
