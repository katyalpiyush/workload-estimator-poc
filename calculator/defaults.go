package calculator

import "workload-estimator-poc/models"

// applyDefaults assigns default values for dataset, workload, and service groups based on workloadNature
func applyDefaults(request *models.ComputeRequest) {
	// Preserve user-provided values for noOfDocuments and averageDocumentSize
	noOfDocuments := request.Dataset.NoOfDocuments
	averageDocumentSize := request.Dataset.AverageDocumentSize

	switch request.WorkloadNature {
	case "read":
		request.ServiceGroups = []models.ServiceGroup{
			{Services: []string{"data", "index", "query"}, NoOfNodes: 3, DiskType: "gp3"},
			{Services: []string{"search"}, NoOfNodes: 2, DiskType: "gp3"},
		}
		request.Dataset = models.Dataset{
			NoOfDocuments: noOfDocuments,
      AverageDocumentSize: averageDocumentSize,
			ResidentRatio:  70,
			PercentIndexesOfDataset:  15,
			PercentFullTextSearchOfDataset:  15,
		}
		request.Workload = models.Workload{
			ReadPerSec:  5000,
			WritesPerSec:  100,
			DeletesPerSec:  50,
			SQLQueriesPerSec:  2000,
		}

	case "write":
		request.ServiceGroups = []models.ServiceGroup{
			{Services: []string{"data"}, NoOfNodes: 3, DiskType: "gp3"},
			{Services: []string{"eventing"}, NoOfNodes: 2, DiskType: "gp3"},
		}
		request.Dataset = models.Dataset{
			NoOfDocuments: noOfDocuments,
      AverageDocumentSize: averageDocumentSize,
			ResidentRatio:  80,
			PercentIndexesOfDataset:  0,
			PercentFullTextSearchOfDataset:  0,
			PercentOperationalAnalyticsOfDataset:  0,
		}
		request.Workload = models.Workload{
			ReadPerSec:  500,
			WritesPerSec:  5000,
			DeletesPerSec:  1000,
			SQLQueriesPerSec:  500,
		}

	case "readwrite":
		request.ServiceGroups = []models.ServiceGroup{
			{Services: []string{"data", "index", "query"}, NoOfNodes: 3, DiskType: "gp3"},
			{Services: []string{"search", "eventing"}, NoOfNodes: 2, DiskType: "gp3"},
		}
		request.Dataset = models.Dataset{
			NoOfDocuments: noOfDocuments,
      AverageDocumentSize: averageDocumentSize,
			ResidentRatio:  60,
			PercentIndexesOfDataset:  40,
			PercentFullTextSearchOfDataset:  0,
			PercentOperationalAnalyticsOfDataset:  0,
		}
		request.Workload = models.Workload{
			ReadPerSec:  2000,
			WritesPerSec:  2000,
			DeletesPerSec:  500,
			SQLQueriesPerSec:  1500,
		}

	case "override":
		// Do nothing, use user-provided values
	}
}
