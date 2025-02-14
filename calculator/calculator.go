package calculator

import (
	"workload-estimator-poc/models"
	"workload-estimator-poc/services"
	"fmt"
	"log"
	"encoding/json"
)

func EstimateResources(request models.ComputeRequest) models.ComputeResponse {
	// Apply default values
	// Log the received request
	requestJSON, err := json.MarshalIndent(request, "", "  ")
	if err == nil {
		log.Printf("Received Request: \n%s\n", string(requestJSON))
	} else {
		log.Println("Failed to log request")
	}

	applyDefaults(&request)

	// Log the received request
	requestJSON2, err := json.MarshalIndent(request, "", "  ")
	if err == nil {
		log.Printf("Received Request: \n%s\n", string(requestJSON2))
	} else {
		log.Println("Failed to log request")
	}

	var serviceGroupResults []models.ServiceGroupResult

	// Iterate over service groups
	for _, group := range request.ServiceGroups {
		nodes := group.NoOfNodes

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
				ram, cpu, disk, diskIO = services.EstimateResourcesForIndex(request.Dataset, request.Workload, nodes)
			case "query":
				ram, cpu, disk, diskIO = services.EstimateResourcesForQuery(request.Dataset, request.Workload, nodes)
			case "search":
				ram, cpu, disk, diskIO = services.EstimateResourcesForSearch(request.Dataset, request.Workload, nodes)
			case "eventing":
				ram, cpu, disk, diskIO = services.EstimateResourcesForEventing(request.Dataset, request.Workload, nodes)
			case "analytics":
				ram, cpu, disk, diskIO = services.EstimateResourcesForAnalytics(request.Dataset, request.Workload, nodes)
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
			ServiceGroupType: fmt.Sprintf("Services: %v", group.Services),  
			EstimatedRAM:     fmt.Sprintf("%.2f GB", totalRAM),
			EstimatedCPU:     fmt.Sprintf("%.2f vCPUs", totalCPU),
			EstimatedDisk:    fmt.Sprintf("%.2f GB", totalDisk),
			EstimatedDiskIO:  fmt.Sprintf("%d IOPS", int(totalDiskIO)),
		})
	}

	// Return the results for each service group
	return models.ComputeResponse{
		ServiceGroupsResults: serviceGroupResults,
	}
}