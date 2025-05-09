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
var genLogger = logger.InitializeLogger()

type DataStore struct {
	// A map containing all the information about the classes required to render the templates.
	Classes map[string]Class
	// The client used to retrieve the meta data from the remote location.
	client *http.Client
	// The host from which the meta data is retrieved.
	metaHost string
	// A list of all the classes that have been retrieved from the remote location.
	// This is used to avoid retrieving the same class multiple times.
	retrievedClasses []string
}

func NewDataStore() (*DataStore, error) {
	dataStore := &DataStore{
		Classes: make(map[string]Class),
		client:  &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
	}
	// Set the meta data host for retrieval of meta files.
	dataStore.setMetaHost()
	// Check if classes are set in the environment variable 'GEN_ACI_TF_META_CLASSES' and retrieve the meta files for those classes.
	err := dataStore.retrieveEnvMetaClassesFromRemote()
	if err != nil {
		return nil, err
	}
	// Refresh the meta files from the remote location if specified in the environment variable 'GEN_ACI_TF_META_REFRESH'.
	// If the environment variable is not set, the default is to not refresh the meta data.
	err = dataStore.refreshMetaFiles()
	if err != nil {
		return nil, err
	}
	// Load classes into the data store for all meta files in the constMetaPath directory.
	err = dataStore.loadClasses()
	if err != nil {
		return nil, err
	}
	return dataStore, nil
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

func (ds *DataStore) retrieveEnvMetaClassesFromRemote() error {
	// Retrieve the meta data for the classes specified in the constEnvMetaClasses environment variable.
	classNames := strings.Split(os.Getenv(constEnvMetaClasses), ",")

	// Avoid attempting retrieval for unset or set to empty ("") environment variable.
	if classNames[0] != "" {
		genLogger.Debug(fmt.Sprintf("Retrieving meta files for classes: %s.", classNames))
		for _, className := range classNames {
			// ENHANCEMENT: Concurrently retrieve/write the meta data.
			// Retrieve the meta file for the class from the remote location.
			err := ds.retrieveMetaFileFromRemote(className)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ds *DataStore) refreshMetaFiles() error {
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
				err = ds.retrieveMetaFileFromRemote(className)
				if err != nil {
					return err
				}
			}
			genLogger.Debug(fmt.Sprintf("Successfully refreshed meta data from remote location: %s.", ds.metaHost))
		}
	}
	return nil
}

func (ds *DataStore) retrieveMetaFileFromRemote(className string) error {
	// Only retrieve the meta file if it is not already retrieved.
	if !slices.Contains(ds.retrievedClasses, className) {
		packageName, shortName, err := splitClassNameToPackageNameAndShortName(className)
		if err != nil {
			return err
		}
		url := fmt.Sprintf(constMetaFileUrl, ds.metaHost, packageName, shortName)
		genLogger.Debug(fmt.Sprintf("Retrieving meta data for class '%s' from: %s.", className, url))

		res, err := ds.client.Get(url)
		if err != nil {
			genLogger.Error(fmt.Sprintf("Error during retrieval of meta file for class '%s': %s.", className, err.Error()))
			return err
		}

		outputFile, err := os.Create(fmt.Sprintf("%s/%s.json", constMetaPath, className))
		if err != nil {
			genLogger.Error(fmt.Sprintf("Error during creation of file for class '%s': %s.", className, err.Error()))
			return err
		}

		defer outputFile.Close()
		_, err = io.Copy(outputFile, res.Body)
		if err != nil {
			genLogger.Error(fmt.Sprintf("Error during writing to file for class '%s': %s.", className, err.Error()))
			return err
		}

		ds.retrievedClasses = append(ds.retrievedClasses, className)

		genLogger.Debug(fmt.Sprintf("Successfully wrote meta data for class '%s' to: %s.", className, outputFile.Name()))
	}

	return nil
}

func (ds *DataStore) loadClasses() error {
	// Load the meta data for all classes in the meta directory.
	genLogger.Debug(fmt.Sprintf("Loading classes from: %s.", constMetaPath))
	for _, className := range utils.GetFileNamesFromDirectory(constMetaPath, true) {
		// Create a new class object and add it to the data store.
		classDetails, err := NewClass(className)
		if err != nil {
			return err
		}
		ds.Classes[className] = *classDetails
	}
	genLogger.Debug(fmt.Sprintf("Successfully loaded classes from: %s.", constMetaPath))
	return nil
}
