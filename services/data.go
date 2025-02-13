package services

import (
	"workload-estimator-poc/models"
)

// EstimateResourcesForData calculates resources required for the Data service.
func EstimateResourcesForData(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(dataset.NoOfDocuments * dataset.AverageDocumentSize * dataset.ResidentRatio) / (1024 * 1024 * 1024)
	
	cpu = float64(workload.ReadPerSec+workload.WritesPerSec+workload.DeletesPerSec) / 100.0
	disk = float64(dataset.NoOfDocuments*dataset.AverageDocumentSize) / (1024 * 1024 * 1024)
	diskIO =float64(workload.ReadPerSec + workload.WritesPerSec + workload.DeletesPerSec) * 10

	// Normalize by number of nodes
	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
