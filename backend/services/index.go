package services

import (
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForIndex calculates resources required for the Index service.
func EstimateResourcesForIndex(dataset models.Dataset, workload models.Workload) (ram, cpu, disk, diskIO float64) {
	cpu = calculateIndexCPU()
	ram = calculateIndexRAM(dataset, cpu)
	disk = calculateIndexDisk(dataset)
	diskIO = calculateIndexDiskIO()
	return ram, cpu, disk, diskIO
}

// calculateIndexRAM computes the RAM required for indexing. (verified)
func calculateIndexRAM(dataset models.Dataset, cpuAvailable float64) float64 {
	// Constants
	const avgKeySize = 0
	const mutationIngestRate = 0                           // comes under the advanced section of sizing calculator
	const primaryIndex = false                             // under advanced section - is false by default as not suitable for production
	const arrayLength = 0                                  // under advanced section - is by default kept to 0
	const arrayIndexElementSize = 0                        // under advanced section - is by default kept to 0
	const sizeOfNonArrayFields = 0                         // under advanced section (Composite Array Index Size) - is by default kept to 0
	const totalSecondaryBytes = 0                          // under advanced section - is by default kept to 0
	const indexDataStructureSize = 114                     // defined as constant in index ram calculations (In bytes)
	const defaultIndex = false                             // under api request payload - is by default false
	const maxMutationQueueSizeOverhead = 256 * 1024 * 1024 // defined as constant in index calculations config
	const fixCostIndexerCommBuffers = 100 * 1024 * 1024    // defined as constant in index calculations config
	const tempAllocationProtobuf = 150 * 1024 * 1024       // defined as constant in index calculations config
	const numReplicas = 1                                  // under advanced section - is by default kept to 1
	const residentRatio = 0.1                              // using it as a constant as used by sizing calculator - its not equal to the resident ratio used for documents
	// const indexType = "plasma"     // or "moi"

	// Step 0: Choose process based on index type
	// use moi_ram_requirement as per sizing calculator
	// currently using index type as "plasma"

	// in sizing calculator this is set via the advanced options
	// the percentage indexes of dataset should ideally work but is not working from the basic section
	absoluteDocumentsInIndex := float64(dataset.NoOfDocuments) * float64(dataset.PercentIndexesOfDataset) / 100

	// FOLLOWING CALCULATIONS ARE FOR COUCHBASE VERSION 7.0 AS USED BY SIZING CALCULATOR

	// Step 1: Versions generated due to MVCC (Multi-Version Concurrency Control)
	versionsGeneratedMvcc := min(mutationIngestRate*20*60, absoluteDocumentsInIndex*3)

	// Step 2: Plasma Memory Usage (In bytes)
	// secondary index
	var plasmaMemUsageSecIdx float64 = 0
	if !primaryIndex && arrayLength == 0 {
		plasmaMemUsageSecIdx = (totalSecondaryBytes + avgKeySize + indexDataStructureSize) * (absoluteDocumentsInIndex + versionsGeneratedMvcc) * 2
	}
	// primary index
	var plasmaMemUsagePrimIdx float64 = 0
	if primaryIndex {
		plasmaMemUsagePrimIdx = (avgKeySize + indexDataStructureSize) * (absoluteDocumentsInIndex + versionsGeneratedMvcc)
	}
	// array index
	var plasmaMemUsageArrIdx float64 = 0
	if !primaryIndex && arrayLength != 0 {
		plasmaMemUsageArrIdx = (((indexDataStructureSize + avgKeySize + (arrayIndexElementSize * arrayLength) + sizeOfNonArrayFields) + ((indexDataStructureSize + avgKeySize + arrayIndexElementSize + sizeOfNonArrayFields) * arrayLength)) * (absoluteDocumentsInIndex + versionsGeneratedMvcc)) * 1.2
	}

	// Step 3: Plasma Write Buffer (In bytes)
	var plasmaWriteBuffer float64 = 2 * 2 * 2 * 1024 * 1024
	if !defaultIndex {
		plasmaWriteBuffer /= 100
	}

	// Step 4: Memory overhead incoming mutation buffers (In bytes)
	var memOverheadMutationBuffer float64 = math.Ceil((totalSecondaryBytes + avgKeySize) * (mutationIngestRate / 1000) * 4000)
	if arrayLength != 0 {
		memOverheadMutationBuffer = math.Ceil((avgKeySize + (arrayIndexElementSize * arrayLength) + sizeOfNonArrayFields) * (mutationIngestRate / 1000) * 4000)
	}

	// Step 5: Calculate encode buffer overhead (In bytes)
	var encodeBufferOverhead float64 = 1794 * math.Ceil(float64(cpuAvailable)*1.2)
	if arrayLength != 0 {
		encodeBufferOverhead = 6660 * math.Ceil(float64(cpuAvailable)*1.2)
	}

	// Step 6: Total memory (In bytes)
	totalMemory := versionsGeneratedMvcc + plasmaMemUsageSecIdx + plasmaMemUsagePrimIdx + plasmaMemUsageArrIdx + plasmaWriteBuffer + memOverheadMutationBuffer + encodeBufferOverhead + maxMutationQueueSizeOverhead + fixCostIndexerCommBuffers + tempAllocationProtobuf

	// Step 7: Overhead for Golang memory management (In GB)
	expMaxMemUsageGB := (totalMemory * 1.05) / (1024 * 1024 * 1024)

	// Step 8: Memory for replicas (In GB)
	expMaxMemUsageGBReplicas := expMaxMemUsageGB + (expMaxMemUsageGB * numReplicas)

	// Step 9: Take into account resident ratio (In GB)
	expMaxMemUsageGBReplicasRR := expMaxMemUsageGBReplicas * residentRatio

	// Step 10: Recommended RAM quota (In GB)
	recommendedRamQuota := math.Ceil(max(expMaxMemUsageGBReplicasRR*1.05, 1))

	return recommendedRamQuota
}

// calculateIndexCPU computes the CPU required for indexing. (verified)
func calculateIndexCPU() float64 {
	// Constants
	const arrayLength = 0
	const defaultIndex = false
	const mutationIngestRate = 0
	const mutationIngestThroughputPerCore = 12500
	const scanRate = 0
	const scanThroughputPerCore = 9000
	const numReplicas = 1

	// Step 0: Choose the type of index.
	// Current process is for plasma index

	// Step 1: Calculate Mutation Cores
	var mutationCoresReq float64 = 0
	if arrayLength == 0 {
		if defaultIndex {
			mutationCoresReq = float64(mutationIngestRate) / float64(mutationIngestThroughputPerCore)
		}
	} else {
		if defaultIndex {
			mutationCoresReq = (float64(mutationIngestRate) * arrayLength) / float64(mutationIngestThroughputPerCore)
		}
	}

	// Step 2: Calculate Scan Cores
	scanCoresReq := float64(scanRate) / float64(scanThroughputPerCore)

	// Step 3: Index Cores Required
	indexCoresReq := mutationCoresReq + scanCoresReq

	// Step 4: Consider replicas
	indexCoresReq += (indexCoresReq * numReplicas)

	// Step 5: Recommended Cores
	recommendedCores := math.Ceil(max(indexCoresReq * 1.2, 1))

	return recommendedCores
}

// calculateIndexDisk computes the disk space required for indexing. (verified)
func calculateIndexDisk(dataset models.Dataset) float64 {
	// Constants
	const avgKeySize = 0
	const primaryIndex = false
	const arrayLength = 0
	const totalSecondaryBytes = 0
	const arrayIndexElementSize = 0
	const sizeOfNonArrayFields = 0
	const numReplicas = 1   

	// in sizing calculator this is set via the advanced options
	// the percentage indexes of dataset should ideally work but is not working from the basic section
	absoluteDocumentsInIndex := float64(dataset.NoOfDocuments) * float64(dataset.PercentIndexesOfDataset) / 100

	// Step 1: Calculate Disk Size (In Bytes)
	// primary index
	var diskSpacePrimaryIdx float64 = 0
	if primaryIndex {
		diskSpacePrimaryIdx = ((absoluteDocumentsInIndex * 2 / 400) * (avgKeySize + 56) * 4 + (avgKeySize + 16) * absoluteDocumentsInIndex * 2) *2
	}
	// secondary index
	var diskSpaceSecondaryIdx float64 = 0
	if !primaryIndex && arrayLength == 0 {
		diskSpaceSecondaryIdx = ((absoluteDocumentsInIndex * 2/ 400) * (totalSecondaryBytes + avgKeySize + 56) * 4 + (totalSecondaryBytes + avgKeySize + 16) * absoluteDocumentsInIndex * 2) * 2
	}
	// array index
	var diskSpaceArrayIdx float64 = 0
	if arrayLength != 0 && !primaryIndex {
		diskSpaceArrayIdx = math.Ceil(((absoluteDocumentsInIndex * 2/ 400) * (arrayIndexElementSize * arrayLength + sizeOfNonArrayFields + avgKeySize + 56) * 4 + (avgKeySize + arrayIndexElementSize * arrayLength + sizeOfNonArrayFields + 16) * absoluteDocumentsInIndex * 2) + (((absoluteDocumentsInIndex * arrayLength) * 2/ 400) * (arrayIndexElementSize + sizeOfNonArrayFields + avgKeySize + 56) * 4) + (avgKeySize + arrayIndexElementSize + sizeOfNonArrayFields + 16) * (absoluteDocumentsInIndex * arrayLength * 2))
	}
	var indexDiskUsage float64 = diskSpacePrimaryIdx + diskSpaceSecondaryIdx + diskSpaceArrayIdx

	// Step 2: Disk Size after snappy (In Bytes)
	var diskSizeAfterSnappy float64 = indexDiskUsage * 0.8

	// Step 3: Fragmentation (In Bytes)
	var fragmentation float64 = indexDiskUsage * 0.3
	
	// Step 4: Expected Max Disk Usage (In GB)
	var expMaxDiskUsage float64 = (diskSizeAfterSnappy + fragmentation) / 1024 / 1024 / 1024

	// Step 5: +20% for DGM and +10% reccomended overhead (In GB)
	var dgmOverheadDiskQuota float64 = expMaxDiskUsage * 1.3

	// Step 6: Consider replicas
	dgmOverheadDiskQuota += (dgmOverheadDiskQuota * numReplicas)

	// Step 7: Recommended Disk Quota
	var recommendedDiskQuota float64 = math.Ceil(max(dgmOverheadDiskQuota, 1))

	return recommendedDiskQuota
}

// calculateIndexDiskIO computes the disk I/O requirement for the Index service. (verified)
func calculateIndexDiskIO() float64 {
	return 0
}
