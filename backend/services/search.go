package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForSearch calculates total resources required for the Search service.
func EstimateResourcesForSearch(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = CalculateSearchRAM()
	cpu = CalculateSearchCPU()
	disk = CalculateSearchDisk(dataset, workload)
	diskIO = CalculateSearchDiskIO()

	return ram, cpu, disk, diskIO
}

// CalculateSearchRAM calculates RAM required for the Search service. (verified)
func CalculateSearchRAM() float64 {
	// Constants
	const maxSize = 0
	const maxFrom = 0
	const searchResultsSize = 112
	const documentMatchStructure = 160
	const scansPerSecond = 0						// Currently taking as constant but required from user

	// RAM Calculation (In GB)
	ram := math.Ceil(round(((maxSize + maxFrom + searchResultsSize) * float64(documentMatchStructure)) / ((1024 * 1024 * 1024)) * scansPerSecond, 2))
	return ram
}

// CalculateSearchCPU calculates CPU required for the Search service. (verified)
func CalculateSearchCPU() float64 {
	// Constants
	const numPartitions = 1 				// under advanced menu in sizing calculator - By default 1, Assuming a single index with a single partition

	// Step 1: CPU Calculation
	cpu := math.Ceil(numPartitions)

	return cpu
}

// CalculateSearchDisk calculates Disk space required for the Search service. (verified)
func CalculateSearchDisk(dataset models.Dataset, workload models.Workload) float64 {
	// Constants
	const index = false
	const store = false
	const includeInAllFields = false
	const includeTermVectors = false
	const docValues = false
	const avgKeySize = 0
	const avgFieldLength = 0
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

	// Step 2: Number of documents in index
	numberOfDocuments := float64(dataset.NoOfDocuments)*(float64(dataset.PercentIndexesOfDataset) / 100)

	// Step 3: Index Size Calculation (In MB)
	var indexSize float64 = 0.0
	if countIfAll == 0 && !index{
		indexSize = round(((float64(dataset.NoOfDocuments) * avgKeySize) + (numberOfDocuments * avgFieldLength * fieldLength * 1.3)) / (1024 * 1024), 0)
	} else {
		indexSize = round(((float64(dataset.NoOfDocuments) * avgKeySize) + (numberOfDocuments * avgFieldLength * fieldLength * float64(countIfAll) * 1.5 * fieldLength)) / (1024 * 1024), 0)
	}

	// Step 4: Disk Space Required
	diskSpace := indexSize * (float64(numReplicas) + 1)

	// Step 5: Convert to GB
	diskSpace = math.Ceil(diskSpace / 1024)

	return diskSpace
}

// CalculateSearchDiskIO calculates Disk I/O required for the Search service. (verified)
func CalculateSearchDiskIO() float64 {
	return 0
}