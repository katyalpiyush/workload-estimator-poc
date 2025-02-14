package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForSearch calculates total resources required for the Search service.
func EstimateResourcesForSearch(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = CalculateSearchRAM()
	cpu = CalculateSearchCPU()
	disk = CalculateSearchDisk(dataset)
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
func CalculateSearchDisk(dataset models.Dataset) float64 {
	// Constants
	const index = false
	const store = false
	const includeInAllFields = false
	const includeTermVectors = false
	const docValues = false
	const avgKeyIDSize = 0
	const avgFieldLength = 1.21 // not sure about this - currently taking it equal to field length value
	const fieldLength = 1.21
	const numReplicas = 0

	// Step 1: Count If All (Sum of boolean values)
	countIfAll := 0
	if index {
		countIfAll++
	}
	if store {
		countIfAll++
	}
	if includeInAllFields {
		countIfAll++
	}
	if includeTermVectors {
		countIfAll++
	}
	if docValues {
		countIfAll++
	}

	// Step 2: Index Size Calculation
	var indexSize float64
	if countIfAll == 0 {
		indexSize = ((float64(dataset.NoOfDocuments)*float64(avgKeyIDSize)) +
			(float64(dataset.NoOfDocuments)*(float64(dataset.PercentIndexesOfDataset) / 100)*float64(avgFieldLength)*float64(fieldLength)*1.3))
	} else {
		indexSize = ((float64(dataset.NoOfDocuments)*float64(avgKeyIDSize)) +
			(float64(dataset.NoOfDocuments)*(float64(dataset.PercentIndexesOfDataset) / 100)*float64(avgFieldLength)*float64(fieldLength)*float64(countIfAll)*1.5*fieldLength))
	}

	// Step 3: Disk Space Required
	diskSpace := indexSize * (float64(numReplicas) + 1)

	// Convert to GB
	diskSpace = math.Ceil(diskSpace / 1024 / 1024 / 1024)
	return diskSpace
}


// CalculateSearchDiskIO calculates Disk I/O required for the Search service.
func CalculateSearchDiskIO(workload models.Workload) float64 {
	const diskIOMultiplier = 8.0
	diskIO := float64(workload.ReadPerSec+workload.WritesPerSec) * diskIOMultiplier
	return diskIO
}