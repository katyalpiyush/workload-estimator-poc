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

// CalculateQueryRAM estimates the RAM required for the Query service.
func CalculateQueryRAM() float64 {
	return 0
}

// CalculateQueryCPU estimates the CPU required for the Query service.
func CalculateQueryCPU(workload models.Workload) float64 {
	// Constants
	const SimpleQueryThroughputPerSecStaleOk = 0
	const SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK = 14000.0 / 24.0
	const SimpleQueryThroughputPerSecStaleFalse = 0
	const SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE = 700.0 / 24.0

	const MediumQueryThroughputPerSecStaleOk = 0
	const MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK = 1500.0 / 24.0
	const MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE = 400.0 / 24.0

	// Simple Query CPU Calculation
	simpleQueryCPU := (float64(SimpleQueryThroughputPerSecStaleOk) / SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK) +
		(float64(SimpleQueryThroughputPerSecStaleFalse) / SIMPLE_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE)

	// Medium Query CPU Calculation
	mediumQueryCPU := (float64(MediumQueryThroughputPerSecStaleOk) / MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_OK) +
		(float64(workload.SQLQueriesPerSec) / MEDIUM_QUERY_QUERIES_PER_SEC_PER_CORE_STALE_FALSE)

	// Complex Query CPU Required is always 0 as there is no formula defined
	complexQueryCPU := 0.0

	// Calculate Total CPU
	totalCPU := math.Ceil(simpleQueryCPU + mediumQueryCPU + complexQueryCPU)

	return totalCPU
}

// CalculateQueryDisk estimates the disk space required for the Query service.
func CalculateQueryDisk() float64 {
	return 0
}

// CalculateQueryDiskIO estimates the disk I/O required for the Query service.
func CalculateQueryDiskIO() float64 {
	return 0
}