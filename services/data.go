package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForData calculates resources required for the Data service.
func EstimateResourcesForData(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = calculateDataRAM(dataset, workload)
	cpu = calculateDataCPU(workload, nodes)
	disk = calculateDataDisk(dataset, nodes)
	diskIO = calculateDataDiskIO(workload, nodes)

	return ram, cpu, disk, diskIO
}

// calculateDataRAM computes the RAM requirement for the Data service.
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
		expiryOpsPerSec = float64(dataset.NoOfDocuments) / (float64(ttlExpiration) * 24 * 3600)
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
	totalActiveMetadataSize := float64(dataset.NoOfDocuments * bucketMetadataSize)
	totalActiveKeysetSize := float64(dataset.NoOfDocuments * avgKeySize)
	totalActiveMetadataKeysetSize := totalActiveMetadataSize + totalActiveKeysetSize

	// Step 4: Calculate Replica Metadata and Keyset Size
	totalReplicaMetadataSize := float64(dataset.NoOfDocuments * bucketMetadataSize * numReplicas)
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
		totalWithJemallocAndTombstones = totalMemoryRequired + (totalMemoryRequired * jemallocBinSize * float64(dataset.ResidentRatio))
	}

	// Add tombstone space if bucket type is Ephemeral
	if bucketType == "Ephemeral" {
		totalWithJemallocAndTombstones += tombstoneSpace
	}

	// Step 8: Calculate Total RAM Quota
	totalRAMQuota := totalWithJemallocAndTombstones / highWaterMark

	// Step 9: Calculate Final Total RAM (convert bytes to GB)
	totalRAM := totalRAMQuota / 1024 / 1024 / 1024

	// Normalize by number of nodes
	return totalRAM
}

// calculateDataCPU computes the CPU requirement for the Data service.
func calculateDataCPU(workload models.Workload, nodes int) float64 {
	cpu := float64(workload.ReadPerSec+workload.WritesPerSec+workload.DeletesPerSec) / 100.0
	return cpu / float64(nodes) // Normalize by number of nodes
}

// calculateDataDisk computes the Disk Space requirement for the Data service.
func calculateDataDisk(dataset models.Dataset, nodes int) float64 {
	disk := float64(dataset.NoOfDocuments*dataset.AverageDocumentSize) / (1024 * 1024 * 1024)
	return disk / float64(nodes) // Normalize by number of nodes
}

// calculateDataDiskIO computes the Disk I/O requirement for the Data service.
func calculateDataDiskIO(workload models.Workload, nodes int) float64 {
	diskIO := float64(workload.ReadPerSec+workload.WritesPerSec+workload.DeletesPerSec) * 10
	return diskIO / float64(nodes) // Normalize by number of nodes
}
