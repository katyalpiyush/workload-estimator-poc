package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForData calculates resources required for the Data service.
func EstimateResourcesForData(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	ram = calculateDataRAM(dataset, workload)
	cpu = calculateDataCPU(dataset, workload)
	disk = calculateDataDisk(dataset, workload)
	diskIO = calculateDataDiskIO(dataset, workload)

	return ram, cpu, disk, diskIO
}

// calculateDataRAM computes the RAM requirement for the Data service. (verified)
func calculateDataRAM(dataset models.Dataset, workload models.Workload) float64 {
	// Constants
	const ttlExpiration = 0
	const avgKeySize = 0
	const inboundXdcrStreams = 0
	const outboundXdcrStreams = 0
	const purgeFrequency = 3
	const numReplicas = 1
	const bucketType = "Couchbase"
	const bucketTypeCouchbase = 56
	const bucketTypeEphemeral = 72
	const compressionRatio = 0.3
	const evictionPolicy = "Full"
	const jemallocBinSize = 0.25
	const highWaterMark = 0.85

	// Step 1: Calculate Expiry Ops Per Second
	var expiryOpsPerSec float64
	if ttlExpiration > 0 {
		expiryOpsPerSec = math.Round(float64(dataset.NoOfDocuments) / (float64(ttlExpiration) * 24 * 3600))
	} else {
		expiryOpsPerSec = 0
	}

	// Step 2: Calculate Tombstone Space (In bytes)
	xdcrFactor := 60 * math.Max(1, float64(inboundXdcrStreams+outboundXdcrStreams))
	tombstoneSpace := (float64(avgKeySize) + xdcrFactor) * float64(purgeFrequency) * float64(numReplicas+1) * 
		(float64(workload.DeletesPerSec) + expiryOpsPerSec) * 60 * 60 * 24

	// Step 3: Calculate Active Metadata and Keyset Size (In bytes)
	var bucketMetadataSize int
	if bucketType == "Couchbase" {
		bucketMetadataSize = bucketTypeCouchbase
	} else {
		bucketMetadataSize = bucketTypeEphemeral
	}
	totalActiveMetadataSize := float64(dataset.NoOfDocuments * int64(bucketMetadataSize))
	totalActiveKeysetSize := float64(dataset.NoOfDocuments * avgKeySize)
	totalActiveMetadataKeysetSize := totalActiveMetadataSize + totalActiveKeysetSize

	// Step 4: Calculate Replica Metadata and Keyset Size
	totalReplicaMetadataSize := float64(dataset.NoOfDocuments * int64(bucketMetadataSize) * numReplicas)
	totalReplicaKeysetSize := float64(dataset.NoOfDocuments * avgKeySize * numReplicas)
	totalReplicaMetadataKeysetSize := totalReplicaMetadataSize + totalReplicaKeysetSize

	// Step 5: Calculate Active and Replica Dataset Sizes (In Bytes)
	activeDatasetSize := float64(dataset.NoOfDocuments) * float64(dataset.AverageDocumentSize) * (1 - compressionRatio)
	replicaDatasetSize := activeDatasetSize * float64(numReplicas)

	// Step 6: Calculate Total Memory Required (In Bytes)
	var totalMemoryRequired float64
	if evictionPolicy == "Value" {
		totalMemoryRequired = ((float64(dataset.ResidentRatio) / 100) * (activeDatasetSize + replicaDatasetSize)) + totalActiveMetadataKeysetSize + totalReplicaMetadataKeysetSize
	} else { // Eviction policy = 'Full'
		totalMemoryRequired = (float64(dataset.ResidentRatio) / 100) * ((activeDatasetSize + replicaDatasetSize) + totalActiveMetadataKeysetSize + totalReplicaMetadataKeysetSize)
	}

	// Step 7: Calculate Total + Jemalloc Bin Size + Tombstones
	var totalWithJemallocAndTombstones float64
	if evictionPolicy == "Value" {
		totalWithJemallocAndTombstones = totalMemoryRequired + (totalMemoryRequired * jemallocBinSize)
	} else { // Eviction policy = 'Full'
		totalWithJemallocAndTombstones = totalMemoryRequired + (totalMemoryRequired * jemallocBinSize * (float64(dataset.ResidentRatio) / 100))
	}

	// Add tombstone space if bucket type is Ephemeral
	if bucketType == "Ephemeral" {
		totalWithJemallocAndTombstones += tombstoneSpace
	}

	// Step 8: Calculate Total RAM Quota
	totalRAMQuota := totalWithJemallocAndTombstones / highWaterMark

	// Step 9: Calculate Final Total RAM (convert bytes to GB)
	totalRAM := totalRAMQuota / 1024 / 1024 / 1024

	// Step 10: Upper bound ram value
	totalRAM = math.Ceil(totalRAM);

	return totalRAM
}

// calculateDataCPU computes the CPU requirement for the Data service. (verified)
func calculateDataCPU(dataset models.Dataset, workload models.Workload) float64 {
	// Constants
	const ttlExpiration = 0
	const inboundXdcrStreams = 0
	const outboundXdcrStreams = 0
	const numberReplicas = 1
	const storageEngine = "Couchstore"  // or "Magma"
	const guardrails_cpu_per_bucket_min = 0.2
	const minimum_number_of_cores_one_bucket = 4

	// Step 1: Compute Expiry Ops Per Second (Same as RAM calculation)
	var expiryOpsPerSec float64
	if ttlExpiration > 0 {
		expiryOpsPerSec = math.Round(float64(dataset.NoOfDocuments) / (float64(ttlExpiration) * 24 * 3600))
	} else {
		expiryOpsPerSec = 0
	}

	// Step 2: Calculate CPU
	// For storage engine of type "Couchstore"
	cpu := float64(inboundXdcrStreams) + float64(outboundXdcrStreams) + (((float64(workload.WritesPerSec) + float64(workload.DeletesPerSec) + expiryOpsPerSec) * numberReplicas) / 10000)
	// For storage engine of type "Magma"
	// Can be added (currently ignored for simplicity)

	// Step 3: Additions based on storage engine type
	if storageEngine == "Couchstore" {
		cpu += 0.4
	}

	// Step 4: Round to 1 decimal place
	cpu = round(cpu, 1)

	// Step 5: Check for gaurdrails minimum
	gaurdrails_minimum := 1 * guardrails_cpu_per_bucket_min				// here 1 corresponds to the number of buckets currently taking it as a single bucket

	// Step 6: Get max cpu
	cpu = max(gaurdrails_minimum, cpu)

	// Step 7: Additions based on storage engine type
	if storageEngine == "Couchstore"{
		cpu += (minimum_number_of_cores_one_bucket - 1)
	}

	// Step 8: Upper bound cpu value
	cpu = math.Ceil(cpu)

	return cpu
}

// utility function
// note : this function is a bit different from python's round function in sizing calculator as it rounds
// to the next integer without any context of the digits being odd or even when ending with 5
// but is recommended as it do not truncate certain values which might result in a lesser value than expected
func round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}

// calculateDataDisk computes the Disk Space requirement for the Data service. (verified)
func calculateDataDisk(dataset models.Dataset, workload models.Workload) float64 {
	// Constants
	const ttlExpiration = 0					// in days
	const avgKeySize = 0						// in bytes
	const inboundXdcrStreams = 0
	const outboundXdcrStreams = 0
	const purgeFrequency = 3
	const numReplicas = 1
	const bucketType = "Couchbase"
	const bucketTypeCouchbase = 56
	const bucketTypeEphemeral = 72
	const compressionRatio = 0.3
	const storageEngine = "Couchstore" // or "Magma"
	var appendOnlyMultiplier = 3

	// Step -1: Set appendOnlyMultiplier
	if storageEngine == "Magma" {
		appendOnlyMultiplier = 2
	}

	// Step 0: Choose process based on type of storage engine
	// If storage engine is of type "Magma" then use calculate_magma_disk_space function in sizing calculator
	// Currently using the process for couchstore only:

	// Step 1: Compute Expiry Ops Per Second (Same as RAM calculation)
	var expiryOpsPerSec float64
	if ttlExpiration > 0 {
		expiryOpsPerSec = math.Round(float64(dataset.NoOfDocuments) / float64(ttlExpiration) / 24 / 3600)			// here ttl is expected in days - converted to seconds
	} else {
		expiryOpsPerSec = 0
	}

	// Step 2: Compute Tombstone Space (In Bytes)
	// count number of streams configured
	var count = 0
	if inboundXdcrStreams > 0 {
		count++
	}
	if outboundXdcrStreams > 0 {
		count++
	}	
	// calculate tombstone space
	tombstoneSpace := math.Round((float64(avgKeySize) + (60 * math.Max(1, float64(count)))) *
		float64(purgeFrequency) * float64(numReplicas+1) * (float64(workload.DeletesPerSec) + expiryOpsPerSec) * 60 * 60 * 24)

	// Step 3: Compute Metadata & Keyset Sizes (In Bytes)
	// active data
	activeMetadataSize := float64(dataset.NoOfDocuments) * bucketTypeCouchbase
	if bucketType == "Ephemeral" {
		activeMetadataSize = float64(dataset.NoOfDocuments) * bucketTypeEphemeral
	}
	activeKeysetSize := float64(dataset.NoOfDocuments) * float64(avgKeySize)
	totalActiveMetadataKeysetSize := activeMetadataSize + activeKeysetSize
	
	// replica data
	replicaMetadataSize := (activeMetadataSize * numReplicas)
	replicaKeysetSize := float64(dataset.NoOfDocuments) * float64(avgKeySize) * float64(numReplicas)
	totalReplicaMetadataKeysetSize := replicaMetadataSize + replicaKeysetSize

	totalMetadataKeysetSize := totalActiveMetadataKeysetSize + totalReplicaMetadataKeysetSize

	// Step 4: Compute Dataset Sizes (In Bytes)
	activeDatasetSize := float64(dataset.NoOfDocuments) * float64(dataset.AverageDocumentSize)
	replicaDatasetSize := activeDatasetSize * float64(numReplicas)

	// Step 5: Compute Size on Disk
	var sizeOnDisk float64
	if bucketType == "Ephemeral" {
		sizeOnDisk = 0
	} else {
		result := activeDatasetSize + replicaDatasetSize
		if compressionRatio != 0 {
			result *= (1 - compressionRatio)
		}
		sizeOnDisk = ((result + totalMetadataKeysetSize) * float64(appendOnlyMultiplier)) + tombstoneSpace
	}

	// Step 6: Convert size into GB and round off to ceil value
	sizeOnDisk = math.Ceil(sizeOnDisk / 1024 / 1024 / 1024)

	return sizeOnDisk
}

// calculateDataDiskIO computes the Disk I/O requirement for the Data service. (verified)
func calculateDataDiskIO( dataset models.Dataset, workload models.Workload) float64 {
	// Constants
	const ttlExpiration = 0   // in days
	const numReplicas = 1
	const bucketType = "Couchbase"

	// Step 0: If storage engine is "Magma" then use different process as per sizing calculator need to use calculate_magma_disk_io function
	// Current process is for storage engine as "Couchstore"

	// Step 1: Compute Expiry Ops Per Second (Same as RAM Calculation)
	var expiryOpsPerSec float64
	if ttlExpiration > 0 {
		expiryOpsPerSec = math.Round(float64(dataset.NoOfDocuments) / (float64(ttlExpiration) * 24 * 3600))
	} else {
		expiryOpsPerSec = 0
	}

	// Step 2: Compute Disk I/O
	var diskIO float64
	if bucketType == "Ephemeral" {
		diskIO = 0
	} else {
		diskIO = (float64(workload.WritesPerSec) + float64(workload.DeletesPerSec) + expiryOpsPerSec) * float64(numReplicas+1)
	}

	return diskIO
}