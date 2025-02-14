package services

import "workload-estimator-poc/models"

// EstimateResourcesForSearch calculates total resources required for the Search service.
func EstimateResourcesForSearch(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = CalculateSearchRAM()
	cpu = CalculateSearchCPU(workload)
	disk = CalculateSearchDisk(ram)
	diskIO = CalculateSearchDiskIO(workload)

	return ram, cpu, disk, diskIO
}

// CalculateSearchRAM calculates RAM required for the Search service.
func CalculateSearchRAM() float64 {
	// Constants
	const maxSize = 0
	const maxFrom = 0
	const searchResultsSize = 112
	const documentMatchStructure = 160
	const scansPerSecond = 0						// Currently taking as constant but required from user

	// RAM Calculation
	ram := (float64(maxSize+maxFrom+searchResultsSize) * float64(documentMatchStructure)) * scansPerSecond / 1024 / 1024 /1024

	return ram
}


// CalculateSearchCPU calculates CPU required for the Search service.
func CalculateSearchCPU() float64 {
	// Constants
	const numPartitions = 1 // Assuming a single index with a single partition

	// CPU Calculation
	cpu := float64(numPartitions)

	return cpu
}

// CalculateSearchDisk calculates Disk space required for the Search service.
func CalculateSearchDisk(ram float64) float64 {
	const diskMultiplier = 2.5
	disk := ram * diskMultiplier
	return disk
}

// CalculateSearchDiskIO calculates Disk I/O required for the Search service.
func CalculateSearchDiskIO(workload models.Workload) float64 {
	const diskIOMultiplier = 8.0
	diskIO := float64(workload.ReadPerSec+workload.WritesPerSec) * diskIOMultiplier
	return diskIO
}