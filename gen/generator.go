package main

import (
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitalizeLogger()

func main() {
	// Ensure that the log file is closed at the end of the main().
	defer genLogger.CloseLogFile()
	genLogger.Info("Implement Generator")
}
