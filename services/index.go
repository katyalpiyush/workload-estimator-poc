package services

import "workload-estimator-poc/models"

// EstimateResourcesForIndex calculates resources required for the Index service.
func EstimateResourcesForIndex(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(dataset.NoOfDocuments * dataset.AverageDocumentSize * dataset.PercentIndexesOfDataset) / (100 * 1024 * 1024 * 1024)
	cpu = float64(workload.SQLQueriesPerSec) / 100.0
	disk = ram * 2
	diskIO = float64(workload.SQLQueriesPerSec) * 5

	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
