package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForEventing calculates resources required for the Eventing service.
func EstimateResourcesForEventing(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = CalculateEventingRAM()
	cpu = CalculateEventingCPU(dataset, workload)
	disk = CalculateEventingDisk()
	diskIO = CalculateEventingDiskIO()

	return ram, cpu, disk, diskIO
}

// CalculateEventingRAM estimates the RAM required for the Eventing service. (verified)
func CalculateEventingRAM() float64 {
	return 0
}

// CalculateEventingCPU estimates the CPU required for the Eventing service. (verified)
func CalculateEventingCPU(dataset models.Dataset, workload models.Workload) float64 {
	const ttlExpiration = 0.0																	// in days
	const numberOfHandlers = 1.0															// under advanced section - default to 1
	const sourceBucketMutationRateFactor = 115000.0 / 24.0		// constant under eventing calculation
	const handlerCountPerCoreFactor = 5.0 / 24.0							// constant under eventing calculation
	const percentageDocsInFunction = 0.0	/ 100.0							// Needs to be taken from user (currently taken as constant)
	const numberOfReadOpsPerExecution = 0.0										// under advanced section - default to 0
	const numberOfWriteOpsPerExecution = 0.0									// under advanced section - default to 0
	const numberOfDeleteOpsPerExecution = 0.0									// under advanced section - default to 0
	const bucketOpsPerCoreFactor = 83000.0 / 24.0							// constant under eventing calculation
	const timersPerCoreFactor = 26500.0 / 24.0								// constant under eventing calculation
	const numberOfTimersCreatedPerExecution = 0.0							// under advanced section - default to 0
	const n1qlPerCoreFactor = 8500.0 / 24.0										// constant under eventing calculation
	const numberOfN1qlQueriesPerExecution = 0.0								// under advanced section - default to 0
	const logPerCoreFactor = 105500.0 / 24.0									// constant under eventing calculation
	const numberOfLogStatementsPerExecution = 0.0							// under advanced section - default to 0
	const curlPerCoreFactor = 1100.0 / 24.0										// constant under eventing calculation
	const numberOfCurlStatementsPerExecution = 0.0						// under advanced section - default to 0

	// Step 1: Calculate Expiry Ops Per Second
	var expiryOpsPerSec float64
	if ttlExpiration > 0 {
		expiryOpsPerSec = math.Round(float64(dataset.NoOfDocuments) / (float64(ttlExpiration) * 24 * 3600))
	} else {
		expiryOpsPerSec = 0
	}

	// Step 2: Calculate Mutation Rate Per Second
	var mutationRatePerSec float64 = float64(workload.WritesPerSec) + float64(workload.DeletesPerSec) + expiryOpsPerSec

	// Step 3: Calculate CPU required for source bucket mutation
	var sourceBucketMutationRateCpuRequired float64 = round((mutationRatePerSec * numberOfHandlers) / sourceBucketMutationRateFactor, 3)

	// Step 4: Calculate CPU required for number of handlers
	var numberOfHandlersCpuRequired float64 = round(numberOfHandlers / handlerCountPerCoreFactor, 3)

	// Step 5: Calculate Bucket Operations CPU required
	var bucketOpsCpuRequired float64 = round(((numberOfReadOpsPerExecution + numberOfWriteOpsPerExecution + numberOfDeleteOpsPerExecution) * (mutationRatePerSec * percentageDocsInFunction)) / bucketOpsPerCoreFactor, 3)

	// Step 6: Calculate Timer Operations CPU required
	var timerOpsCpuRequired float64 = round((mutationRatePerSec * percentageDocsInFunction) / timersPerCoreFactor * numberOfTimersCreatedPerExecution, 3)

	// Step 7: Calculate N1QL Operations CPU required
	var n1qlOpsCpuRequired float64 = round((mutationRatePerSec * percentageDocsInFunction) / n1qlPerCoreFactor * numberOfN1qlQueriesPerExecution, 3)

	// Step 8: Calculate Log Operations CPU required
	var logOpsCpuRequired float64 = round((mutationRatePerSec * percentageDocsInFunction) / logPerCoreFactor * numberOfLogStatementsPerExecution, 3)

	// Step 9: Calculate CURL Operations CPU required
	var curlOpsCpuRequired float64 = round((mutationRatePerSec * percentageDocsInFunction) / curlPerCoreFactor * numberOfCurlStatementsPerExecution, 3)

	// Step 10: Calculate total CPU required
	var totalCpuRequired = sourceBucketMutationRateCpuRequired + numberOfHandlersCpuRequired + bucketOpsCpuRequired + timerOpsCpuRequired + n1qlOpsCpuRequired + logOpsCpuRequired + curlOpsCpuRequired

	// Step 11: Apply upper bound
	var clusterCpuRequired = math.Ceil(totalCpuRequired)

	return clusterCpuRequired
}

// CalculateEventingDisk estimates the Disk required for the Eventing service. (verified)
func CalculateEventingDisk() float64 {
	return 0
}

// CalculateEventingDiskIO estimates the Disk I/O required for the Eventing service. (verified)
func CalculateEventingDiskIO() float64 {
	return 0
}