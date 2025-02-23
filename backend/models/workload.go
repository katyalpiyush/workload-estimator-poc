package models

// ServiceGroup represents a group of services, their node count, and disk type
type ServiceGroup struct {
	Name      		string   		`json:"name"`
	Services  		[]string 		`json:"services"`
	NoOfNodes 		int64    		`json:"no_of_nodes"`
	DiskType  		string   		`json:"disk_type"`
}

// Dataset represents the dataset characteristics for estimation
type Dataset struct {
	NoOfDocuments                    		 		int64 		`json:"no_of_documents"`
	AverageDocumentSize              		 		int64 		`json:"average_document_size"`
	ResidentRatio                    		 		int64 		`json:"resident_ratio"`
	PercentIndexesOfDataset          		 		int64 		`json:"percent_indexes_of_dataset"`
	PercentFullTextSearchOfDataset   		 		int64 		`json:"percent_full_text_search_of_dataset"`
	PercentOperationalAnalyticsOfDataset 		int64 		`json:"percent_operational_analytics_of_dataset"`
}

// Workload represents the workload characteristics for estimation
type Workload struct {
	ReadPerSec       		int64 				`json:"read_per_sec"`
	WritesPerSec     		int64 				`json:"writes_per_sec"`
	DeletesPerSec    		int64 				`json:"deletes_per_sec"`
	SQLQueriesPerSec 		int64 				`json:"sql_queries_per_sec"`
}

// ComputeRequest is the input request format for the workload estimation
type ComputeRequest struct {
	ServiceGroups 		[]ServiceGroup 	`json:"service_groups"`
	Dataset       		Dataset        	`json:"dataset"`
	Workload      		Workload       	`json:"workload"`
	WorkloadNature 		string        	`json:"workload_nature"`
}

// Summary holds the overview of the estimations of all the service groups
type Summary struct {
	ClusterOption				string				`json:"cluster_option"`
	NodesAllocated			int64					`json:"nodes_allocated"`
	ServiceGroups				int64					`json:"service_groups"`
	Services 						[]string			`json:"services"`
	WorkloadType				string				`json:"workload_type"`
}

// ServiceGroupResult holds the resource estimates for each service group
type ServiceGroupResult struct {
	Services 				 		[]string  		`json:"services"`
	Nodes 					 		int64 				`json:"nodes"`
	EstimatedRAM     		int64  				`json:"estimated_ram"`
	EstimatedCPU     		int64  				`json:"estimated_cpu"`
	DiskType						string 				`json:"disk_type"`
	EstimatedDisk    		int64  				`json:"estimated_disk"`
	EstimatedDiskIO  		int64  				`json:"estimated_disk_io"`
}

// ComputeResponse is the output response format containing results for all service groups
type ComputeResponse struct {
	Summary							  Summary								`json:"summary"`
	ServiceGroupsResults  []ServiceGroupResult 	`json:"service_groups_results"`
}