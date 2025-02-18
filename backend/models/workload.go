package models

// ServiceGroup represents a group of services, their node count, and disk type
type ServiceGroup struct {
	Name      string   `json:"name"`
	Services  []string `json:"services"`
	NoOfNodes int      `json:"no_of_nodes"`
	DiskType  string   `json:"disk_type"`
}

// Dataset represents the dataset characteristics for estimation
type Dataset struct {
	NoOfDocuments                    		 int `json:"no_of_documents"`
	AverageDocumentSize              		 int `json:"average_document_size"`
	ResidentRatio                    		 int `json:"resident_ratio"`
	PercentIndexesOfDataset          		 int `json:"percent_indexes_of_dataset"`
	PercentFullTextSearchOfDataset   		 int `json:"percent_full_text_search_of_dataset"`
	PercentOperationalAnalyticsOfDataset int `json:"percent_operational_analytics_of_dataset"`
}

// Workload represents the workload characteristics for estimation
type Workload struct {
	ReadPerSec       int `json:"read_per_sec"`
	WritesPerSec     int `json:"writes_per_sec"`
	DeletesPerSec    int `json:"deletes_per_sec"`
	SQLQueriesPerSec int `json:"sql_queries_per_sec"`
}

// ComputeRequest is the input request format for the workload estimation
type ComputeRequest struct {
	ServiceGroups []ServiceGroup `json:"service_groups"`
	Dataset       Dataset        `json:"dataset"`
	Workload      Workload       `json:"workload"`
	WorkloadNature string        `json:"workload_nature"`
}

// ServiceGroupResult holds the resource estimates for each service group
type ServiceGroupResult struct {
	ServiceGroupType string  `json:"service_group_type"`
	EstimatedRAM     string  `json:"estimated_ram"`
	EstimatedCPU     string  `json:"estimated_cpu"`
	EstimatedDisk    string  `json:"estimated_disk"`
	EstimatedDiskIO  string  `json:"estimated_disk_io"`
}

// ComputeResponse is the output response format containing results for all service groups
type ComputeResponse struct {
	ServiceGroupsResults []ServiceGroupResult `json:"service_groups_results"`
}
