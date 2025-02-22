package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForQuery calculates resources required for the Query service.
func EstimateResourcesForQuery(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = CalculateQueryRAM()
	cpu = CalculateQueryCPU(workload)
	disk = CalculateQueryDisk()
	diskIO = CalculateQueryDiskIO()

	return ram, cpu, disk, diskIO
}

// CalculateQueryRAM estimates the RAM required for the Query service. (verified)
func CalculateQueryRAM() float64 {
	return 0
}

// CalculateQueryCPU estimates the CPU required for the Query service. (verified)
func CalculateQueryCPU(workload models.Workload) float64 {
	// Constants
	const simpleQueryThroughputPerSecStaleOk = 0
	const SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK = 14000.0 / 24.0
	const simpleQueryThroughputPerSecStaleFalse = 0
	const SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE = 700.0 / 24.0
	
	const mediumQueryThroughputPerSecStaleOk = 0
	const MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK = 1500.0 / 24.0
	var mediumQueryThroughputPerSecStaleFalse = workload.SQLQueriesPerSec
	const MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE = 400.0 / 24.0

	const complexQueryThroughputPerSecStaleOk = 0
	const COMPLEX_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK = 700.0 / 24.0
	const complexQueryThroughputPerSecStaleFalse = 0
	const COMPLEX_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE = 200.0 / 24.0

	// Step 1: Simple Query CPU Calculation
	simpleQueryCPU := round((simpleQueryThroughputPerSecStaleOk / SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK) + (simpleQueryThroughputPerSecStaleFalse / SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE), 2)

	// Step 2: Medium Query CPU Calculation
	mediumQueryCPU := round((mediumQueryThroughputPerSecStaleOk / MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK)	+ (float64(mediumQueryThroughputPerSecStaleFalse) / MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE), 2)

	// Step 3: Complex Query CPU Calculation
	complexQueryCPU := round((complexQueryThroughputPerSecStaleOk / COMPLEX_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK) + (complexQueryThroughputPerSecStaleFalse / COMPLEX_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE), 2)

	// Step 4: Calculate Total CPU
	totalCPU := math.Ceil(simpleQueryCPU + mediumQueryCPU + complexQueryCPU)

	return totalCPU
}

// CalculateQueryDisk estimates the disk space required for the Query service. (verified)
func CalculateQueryDisk() float64 {
	return 0
}

// CalculateQueryDiskIO estimates the disk I/O required for the Query service. (verified)
func CalculateQueryDiskIO() float64 {
	return 0
}