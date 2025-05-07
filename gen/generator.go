package main

import (
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/data"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitializeLogger()

func main() {
	genLogger.Trace("Starting Generator")
	// Ensure that the log file is closed at the end of the main().
	defer genLogger.CloseLogFile()

	genLogger.Trace("Initializing Data Store")
	_, err := data.NewDataStore()
	if err != nil {
		genLogger.Fatal("Error during initialization of Data Store: " + err.Error())
	}
	genLogger.Trace("Initialization of Data Store complete")
}
