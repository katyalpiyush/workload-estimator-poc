package services

import "workload-estimator-poc/models"

// EstimateResourcesForAnalytics calculates resources required for the Analytics service.
func EstimateResourcesForAnalytics(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = float64(dataset.NoOfDocuments * dataset.AverageDocumentSize * dataset.PercentOperationalAnalyticsOfDataset) / (100 * 1024 * 1024 * 1024)
	cpu = float64(workload.SQLQueriesPerSec) / 50.0
	disk = ram * 3
	diskIO = float64(workload.SQLQueriesPerSec) * 6

	// Normalize by number of nodes
	ram /= float64(nodes)
	cpu /= float64(nodes)
	disk /= float64(nodes)
	diskIO /= float64(nodes)

	return ram, cpu, disk, diskIO
}
