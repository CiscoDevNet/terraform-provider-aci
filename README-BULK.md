# Added BULK mode to BD and EPG

- **Branch source: v2.7.0**
- **Goal: improve Create, Read and Update performance for Bridge Domain and Endpoint Group resources**

This code adds BULK mode to the CREATE, READ, and UPDATE operations for Bridge Domain and Endpoint Group resources.
BULK mode refers to the ability to CREATE, READ and UPDATE the available relationships for the above resources in a single bulk operation instead of one at a time.
This brings absolute performance increases in the READ operation, performance increases proportional to the number of relations involved in the CREATE and UPDATE operations.
The BULK mode can be enabled operation by operation resource by resource, to facilitate debugging and allow for a correct evaluation of performance while still having the traditional operating mode by default.

## Added examples for BULK operations

In the files bd.tf, epg.tf, bd1.tf.ste, bd2.tf.ste, epg1.tf.ste, epg2.tf.ste, available in the directory examples/bulk, configure the Terrafom <count> Meta-Argument consistently with your strategy (we tested the performance of BULK operations by setting <count=1000> instead of the default <count=1>).
There are two families of tests, with the purpose of modifying the parameters to test the BULK update operation.
The t1-test-suite.sh and t2-test-suite.sh scripts reconfigure the environment with the corresponding test family.
The run.sh script runs the selected test suite, the plan.sh script runs the scheduling (mainly used to check the performance of BULK Read operations).
Both scripts report the start time and end time automatically.

## Warnings

This plugin adds the bulk operations fields to Terraform's state file.
Bulk operations are experimental, be careful using them in production setups.
