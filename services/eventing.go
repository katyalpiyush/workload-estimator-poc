package services

import "workload-estimator-poc/models"

// EstimateResourcesForEventing calculates resources required for the Eventing service.
func EstimateResourcesForEventing(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(workload.WritesPerSec) * 0.1 // Dummy calculation
	cpu = float64(workload.WritesPerSec) / 200.0
	disk = ram * 1.2
	diskIO = float64(workload.WritesPerSec) * 5

	// Normalize by number of nodes
	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
