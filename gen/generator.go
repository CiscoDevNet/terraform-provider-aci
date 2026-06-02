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

	// Policy: per-class authoring mistakes (missing roles, unresolved placeholders,
	// invalid YAML enum values) accumulate on ctx.Diagnostics and are reported as a
	// single summary below. Hard failures (I/O, JSON, schema parsing, version
	// parsing) still return error immediately and abort generation via Fatalf.
	ctx := data.NewContext()

	genLogger.Trace("Initializing Data Store")
	_, err := data.NewDataStore(ctx)
	if err != nil {
		genLogger.Fatalf("Error during initialization of Data Store: %s", err.Error())
	}

	if err := ctx.Diagnostics.Error(); err != nil {
		genLogger.Fatalf("Generation failed:\n%s", err)
	}

	genLogger.Trace("Initialization of Data Store complete")
}
