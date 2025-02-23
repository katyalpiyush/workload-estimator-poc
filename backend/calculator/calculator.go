package calculator

import (
	"workload-estimator-poc/models"
	"workload-estimator-poc/services"
	"fmt"
	"log"
	"encoding/json"
)

const CLUSTER_OPTION = "Custom"

func EstimateResources(request models.ComputeRequest) models.ComputeResponse {
	// Log the received request
	requestJSON, err := json.MarshalIndent(request, "", "  ")
	if err == nil {
		log.Printf("Received Request: \n%s\n", string(requestJSON))
		} else {
			log.Println("Failed to log request")
		}
		
	// Apply default values
	applyDefaults(&request)

	// Log the received request
	requestJSON2, err := json.MarshalIndent(request, "", "  ")
	if err == nil {
		log.Printf("Received Request: \n%s\n", string(requestJSON2))
	} else {
		log.Println("Failed to log request")
	}

	var serviceGroupResults []models.ServiceGroupResult
	var nodesAllocated int64 = 0
	var servicesAll []string

	// Iterate over service groups
	for _, group := range request.ServiceGroups {
		nodes := group.NoOfNodes
		nodesAllocated += nodes
		servicesAll = append(servicesAll, group.Services...)

		// Lists to store resources for each service in the current group
		var ramList, cpuList, diskList, diskIOList []float64

		// Iterate over services within the service group
		for _, service := range group.Services {
			var ram, cpu, disk, diskIO float64

			// Call the corresponding service calculator
			switch service {
			case "data":
				ram, cpu, disk, diskIO = services.EstimateResourcesForData(request.Dataset, request.Workload)
			case "index":
				ram, cpu, disk, diskIO = services.EstimateResourcesForIndex(request.Dataset, request.Workload)
			case "query":
				ram, cpu, disk, diskIO = services.EstimateResourcesForQuery(request.Dataset, request.Workload)
			case "search":
				ram, cpu, disk, diskIO = services.EstimateResourcesForSearch(request.Dataset, request.Workload)
			case "eventing":
				ram, cpu, disk, diskIO = services.EstimateResourcesForEventing(request.Dataset, request.Workload)
			case "analytics":
				ram, cpu, disk, diskIO = services.EstimateResourcesForAnalytics(request.Dataset, request.Workload, int(nodes))
			}

			// Add resources to corresponding lists
			ramList = append(ramList, ram)
			cpuList = append(cpuList, cpu)
			diskList = append(diskList, disk)
			diskIOList = append(diskIOList, diskIO)
		}

		// Print the lists for debugging
		fmt.Printf("\nService Group: %v\n", group.Services)
		fmt.Printf("RAM List: %v\n", ramList)
		fmt.Printf("CPU List: %v\n", cpuList)
		fmt.Printf("Disk List: %v\n", diskList)
		fmt.Printf("Disk IO List: %v\n\n", diskIOList)

		// Calculate the total resources for this service group
		totalRAM := CalculateRAM(ramList, len(group.Services), nodes)
		totalCPU := CalculateCPU(cpuList, len(group.Services), nodes)
		totalDisk := CalculateDisk(diskList, nodes)
		totalDiskIO := CalculateDiskIO(diskIOList, nodes, group.DiskType)

		// Find the closest instance for RAM and CPU
		selectedInstance := findClosestInstance(totalCPU, totalRAM)

		// Update the RAM and CPU with the selected instance's values
		totalRAM = float64(selectedInstance.RAM)
		totalCPU = float64(selectedInstance.VCPU)

		// Store the result for this service group
		serviceGroupResults = append(serviceGroupResults, models.ServiceGroupResult{
			Services: 				group.Services,
			Nodes:						group.NoOfNodes,
			EstimatedRAM:			int64(totalRAM),
			EstimatedCPU:     int64(totalCPU),
			DiskType: 				group.DiskType,
			EstimatedDisk:    int64(totalDisk),
			EstimatedDiskIO:  int64(totalDiskIO),
		})
	}

	// Create summary for all the service groups
	var summary models.Summary;
	summary.ClusterOption = CLUSTER_OPTION
	summary.NodesAllocated = nodesAllocated
	summary.ServiceGroups = int64(len(request.ServiceGroups))
	summary.Services = servicesAll
	summary.WorkloadType = request.WorkloadNature

	// Return the results for all the service groups
	return models.ComputeResponse{
		Summary: summary,
		ServiceGroupsResults: serviceGroupResults,
	}
}