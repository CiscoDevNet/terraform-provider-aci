package data

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitalizeLogger()

type DataStore struct {
	// A map containing all the information about the classes required to render the templates.
	Classes map[string]Class
	// The host from which the meta data is retrieved.
	metaHost string
	// A list of all the classes that have been retrieved from the remote location.
	// This is used to avoid retrieving the same class multiple times.
	retrievedClasses []string
}

func NewDataStore() *DataStore {
	dataStore := &DataStore{Classes: make(map[string]Class)}
	// Set the meta data host for retrieval of meta files.
	dataStore.setMetaHost()
	// Check if classes are set in the environment variable 'GEN_ACI_TF_META_CLASSES' and retrieve the meta files for those classes.
	dataStore.retrieveEnvMetaClassesFromRemote()
	// Refresh the meta files from the remote location if specified in the environment variable 'GEN_ACI_TF_META_REFRESH'.
	// If the environment variable is not set, the default is to not refresh the meta data.
	dataStore.refreshMetaFiles()
	// Load classes into the data store for all meta files in the constMetaPath directory.
	dataStore.loadClasses()
	return dataStore
}

func (ds *DataStore) setMetaHost() {
	// Check if the meta data host is set in the constEnvMetaHost environment variable.
	// If it is set, use it as the host for the meta data retrieval.
	// If it is not set, use the default host defined in constPubhubDevnetHost.
	host := os.Getenv(constEnvMetaHost)
	if host == "" {
		host = constPubhubDevnetHost
	}
	ds.metaHost = host
	genLogger.Info(fmt.Sprintf("Meta data host set to: %s.", host))
}

func (ds *DataStore) retrieveEnvMetaClassesFromRemote() {
	// Retrieve the meta data for the classes specified in the constEnvMetaClasses environment variable.
	classNames := strings.Split(os.Getenv(constEnvMetaClasses), ",")

	// Avoid attempting retrieval for unset or set to empty ("") environment variable.
	if classNames[0] != "" {
		genLogger.Debug(fmt.Sprintf("Retrieving meta files for classes: %s.", classNames))
		for _, className := range classNames {
			// ENHANCEMENT: Concurrently retrieve/write the meta data.
			// Retrieve the meta file for the class from the remote location.
			ds.retrieveMetaFileFromRemote(className)
		}
	}
}

func (ds *DataStore) refreshMetaFiles() {
	// Check if the meta data should be refreshed from remote location as specified in the environment variable.
	metaRefresh := os.Getenv(constEnvMetaRefresh)
	// If the environment variable is not set, the default is to not refresh the meta data.
	if metaRefresh != "" {
		var refresh bool
		var err error
		// If the environment variable is set, parse it to a boolean value.
		// If the parsing fails, log a warning and skip retrieval of meta data.
		// Parsing accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
		// Any other value returns an error, which is logged and the retrieval is skipped.
		refresh, err = strconv.ParseBool(metaRefresh)
		if err != nil {
			genLogger.Warn(fmt.Sprintf("Refreshing of meta is skipped due to error: %s.", err.Error()))
		} else if refresh {
			genLogger.Debug(fmt.Sprintf("Refreshing meta data from remote location: %s.", ds.metaHost))
			for _, className := range utils.GetFileNamesFromDirectory(constMetaPath, true) {
				// ENHANCEMENT: Concurrently retrieve/write the meta data.
				// Retrieve the meta file for the class from the remote location.
				// Only retrieve the meta file if it is not already retrieved.
				if !slices.Contains(ds.retrievedClasses, className) {
					ds.retrieveMetaFileFromRemote(className)
				}
			}
			genLogger.Debug(fmt.Sprintf("Succesfully refreshed meta data from remote location: %s.", ds.metaHost))
		}

	}
}

func (ds *DataStore) retrieveMetaFileFromRemote(className string) {
	shortName, packageName := splitClassNameToPackageNameAndShortName(className)
	url := fmt.Sprintf(constMetaFileUrl, ds.metaHost, packageName, shortName)
	genLogger.Debug(fmt.Sprintf("Retrieving meta data for class '%s' from: %s.", className, url))

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Get(url)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during retrieval of meta file for class '%s': %s.", className, err.Error()))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during reading of file for class '%s': %s.", className, err.Error()))
	}

	outputFile, err := os.Create(fmt.Sprintf("%s/%s.json", constMetaPath, className))
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during creation of file for class '%s': %s.", className, err.Error()))
	}

	defer outputFile.Close()
	_, err = outputFile.Write(resBody)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during writing to file for class '%s': %s.", className, err.Error()))
	}

	ds.retrievedClasses = append(ds.retrievedClasses, className)

	genLogger.Debug(fmt.Sprintf("Succesfully wrote meta data for class '%s' to: %s.", className, outputFile.Name()))
}

func (ds *DataStore) loadClasses() {
	// Load the meta data for all classes in the meta directory.
	genLogger.Debug(fmt.Sprintf("Loading classes from: %s.", constMetaPath))
	for _, className := range utils.GetFileNamesFromDirectory(constMetaPath, true) {
		// Create a new class object and add it to the data store.
		ds.Classes[className] = *NewClass(className)
	}
	genLogger.Debug(fmt.Sprintf("Succesfully loaded classes from: %s.", constMetaPath))
}
