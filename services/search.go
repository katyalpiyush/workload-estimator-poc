package services

import "workload-estimator-poc/models"

// EstimateResourcesForSearch calculates resources required for the Search service.
func EstimateResourcesForSearch(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(dataset.NoOfDocuments * dataset.AverageDocumentSize * dataset.PercentFullTextSearchOfDataset) / (100 * 1024 * 1024 * 1024)
	cpu = float64(workload.ReadPerSec+workload.WritesPerSec) / 150.0
	disk = ram * 2.5
	diskIO = float64(workload.ReadPerSec+workload.WritesPerSec) * 8

	// Normalize by number of nodes
	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
