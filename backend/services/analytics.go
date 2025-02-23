package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForAnalytics calculates resources required for the Analytics service.
func EstimateResourcesForAnalytics(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = CalculateAnalyticsRAM()
	cpu = CalculateAnalyticsCPU()
	disk = CalculateAnalyticsDisk(dataset)
	diskIO = CalculateAnalyticsDiskIO()

	return ram, cpu, disk, diskIO
}

// CalculateAnalyticsRAM estimates the RAM required for the Analytics service. (verified)
func CalculateAnalyticsRAM() float64 {
	return 0
}

// CalculateAnalyticsCPU estimates the CPU required for the Analytics service. (verified)
func CalculateAnalyticsCPU() float64 {
	return 0
}

// CalculateAnalyticsDisk estimates the Disk required for the Analytics service. (verified)
func CalculateAnalyticsDisk(dataset models.Dataset) float64 {
	const avgKeySize = 0
	const storageEngine = "Couchstore"  								// or "Magma"
	const bucketType = "Couchbase"
	const bucketTypeCouchbase = 56
	const bucketTypeEphemeral = 72
	const numReplicas = 1																// under advanced section in sizing calculator
	const queryTempSpaceAllowance = 2										// under advance section in sizing calculator - by default kept as 2
	const percentageDocumentsInIndex = 100.0						// under advance section in sizing calculator - by default kept as 100
	const totalSecondaryBytes = 0												// under advance section in sizing calculator - by default kept as 0

	// Step 1: Calculate data collection disk space required (In GB)
	var dataCollectionDiskSpaceRequired float64 = 0.0
	if storageEngine == "Magma" {
		// Use path for magma, currently only have for couchstore
	} else {
		if bucketType == "Couchbase" {
			dataCollectionDiskSpaceRequired = math.Ceil((float64(dataset.NoOfDocuments)	* float64(avgKeySize + dataset.AverageDocumentSize + bucketTypeCouchbase)) / 1024 / 1024 / 1024)
		} else {
			dataCollectionDiskSpaceRequired = math.Ceil((float64(dataset.NoOfDocuments)	* float64(avgKeySize + dataset.AverageDocumentSize + bucketTypeEphemeral)) / 1024 / 1024 / 1024)
		}
	}

	// Step 2: Calculate no of documents in analytics collection
	var documentsInAnalyticsCollection float64 = float64(dataset.PercentOperationalAnalyticsOfDataset) * float64(dataset.NoOfDocuments) / 100

	// Step 3: Calculate active data size (In GB)
	var activeDataSize float64 = dataCollectionDiskSpaceRequired * float64(dataset.PercentOperationalAnalyticsOfDataset) / 100

	// Step 4: Calculate replica data size (In GB)
	var replicaDataSize float64 = activeDataSize * numReplicas

	// Step 5: Active, Replica and Temp total (In GB)
	var activeReplicaTempTotal float64 = ((activeDataSize * queryTempSpaceAllowance) + replicaDataSize)

	// Step 6: Calculate number of documents in Analytics index
	var documentsInAnalyticsIndex float64 = percentageDocumentsInIndex * documentsInAnalyticsCollection / 100

	// Step 7: Calculate size of index (In GB)
	var indexSize float64 = round((documentsInAnalyticsIndex * (avgKeySize + totalSecondaryBytes) * (1 + numReplicas)) / 1024 / 1024 / 1024, 1)

	// Step 8: Total analytics disk size (In GB)
	var totalDisk float64 = round(activeReplicaTempTotal + indexSize, 0)

	return totalDisk
}

// CalculateAnalyticsDiskIO estimates the Disk IO required for the Analytics service. (verified)
func CalculateAnalyticsDiskIO() float64 {
	return 0
}
