package services

import "workload-estimator-poc/models"

// EstimateResourcesForQuery calculates resources required for the Query service.
func EstimateResourcesForQuery(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(workload.SQLQueriesPerSec) * 0.05 // Dummy calculation
	cpu = float64(workload.SQLQueriesPerSec) / 80.0
	disk = ram * 1.5
	diskIO = float64(workload.SQLQueriesPerSec) * 3

	// Normalize by number of nodes
	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
