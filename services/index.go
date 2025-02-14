package services

import (
	"fmt"
	"math"
	"workload-estimator-poc/models"
)

// EstimateResourcesForIndex calculates resources required for the Index service.
func EstimateResourcesForIndex(dataset models.Dataset, workload models.Workload, nodes int) (ram, cpu, disk, diskIO float64) {
	ram = calculateIndexRAM(dataset, 4)
	// NEED TO CHANGE THE -----------^ PARAMETER (CURRENTLY HARDCODED)

	cpu = calculateIndexCPU(workload)
	disk = calculateIndexDisk(dataset)
	diskIO = calculateIndexDiskIO()
	return ram, cpu, disk, diskIO
}

// calculateIndexRAM computes the RAM required for indexing.
func calculateIndexRAM(dataset models.Dataset, cpuAvailable int) float64 {
	// Constants
	const plasmaKeySize = 0
	const documentKeyIDSize = 0
	const numIndexes = 1
	const purgeRatio = 0.2
	const residentRatio = 0.1
	const compression = 0.5
	const compressionRatio = 2.5
	const jemallocFragmentation = 0.4

	numDocsIndex := float64(dataset.NoOfDocuments) * float64(dataset.PercentIndexesOfDataset) / 100
	fmt.Printf("Number of Documents in Index: %f\n", numDocsIndex)

	// Item Size
	itemSize := plasmaKeySize
	if plasmaKeySize == 0 {
		itemSize = documentKeyIDSize + 12 + 17
	}
	fmt.Printf("Item Size: %d\n", itemSize)

	// Items per Page
	itemSizeFactor := (2.0/3.0)*float64(itemSize) + (1.0/3.0)*(float64(itemSize)+50)
	itemsPerPage := math.Floor(192.0 * 1024.0 / itemSizeFactor)
	if itemsPerPage >= 300 {
		itemsPerPage = 300
	}
	fmt.Printf("Items per Page: %f\n", itemsPerPage)

	// Page Size (In Bytes)
	pageSize := (float64(itemSize) + 50) * itemsPerPage
	fmt.Printf("Page Size: %f\n", pageSize)

	// Number of Shards
	numShards := math.Ceil(float64(numIndexes)/5) * 2
	fmt.Printf("Number of Shards: %f\n", numShards)

	// Index Memory (In GB)
	indexMemory := (float64(itemSize) + 50) * numDocsIndex * (1+purgeRatio) * 2 / (1024 * 1024 * 1024)
	fmt.Printf("Index Memory: %f GB\n", indexMemory)

	// Index Memory After Compaction
	indexMemoryAfterCompaction := (indexMemory*1024*1024*1024 - (numDocsIndex*(1+purgeRatio)*2.0/3.0*40.0*2.0)) / (1024 * 1024 * 1024)
	fmt.Printf("Index Memory After Compaction: %f GB\n", indexMemoryAfterCompaction)

	// Index Memory After Resident Ratio
	indexMemoryAfterResidentRatio := indexMemoryAfterCompaction * residentRatio
	fmt.Printf("Index Memory After Resident Ratio: %f GB\n", indexMemoryAfterResidentRatio)

	// Index Memory After Compression
	indexMemoryAfterCompression := indexMemoryAfterResidentRatio * (compression/compressionRatio + (1 - compression))
	fmt.Printf("Index Memory After Compression: %f GB\n", indexMemoryAfterCompression)

	// Index Memory After Skiplist
	indexMemoryAfterSkiplist := indexMemoryAfterCompression + numDocsIndex/itemsPerPage*(float64(itemSize)+32)*2/(1024*1024*1024)
	fmt.Printf("Index Memory After Skiplist: %f GB\n", indexMemoryAfterSkiplist)

	// Index Memory After Buffer Overhead
	indexMemoryAfterBufferOverhead := (indexMemoryAfterSkiplist*1024 + numShards*4 + pageSize*float64(cpuAvailable)*numShards*2*1.1/1024/1024) / 1024
	fmt.Printf("Index Memory After Buffer Overhead: %f GB\n", indexMemoryAfterBufferOverhead)

	// Index Memory After Jemalloc Fragmentation
	indexMemoryAfterJemalloc := indexMemoryAfterBufferOverhead + indexMemoryAfterBufferOverhead*jemallocFragmentation/(1-jemallocFragmentation)
	fmt.Printf("Index Memory After Jemalloc Fragmentation: %f GB\n", indexMemoryAfterJemalloc)

	// Recommended RAM Quota
	ramQuota := indexMemoryAfterJemalloc * 1.05
	fmt.Printf("Recommended RAM Quota: %f GB\n", ramQuota)

	return ramQuota
}

// calculateIndexCPU computes the CPU required for indexing.
func calculateIndexCPU(workload models.Workload) float64 {
	// Constants
	const arrayIndexSizeOfEachElement = 0
	const mutationRatePerSecond = 0
	const mutationIngestThroughputPerCore = 12500
	const avgIndexScansPerSecond = 0
	const scanThroughputPerCore = 9000
	const arrayLength = 0

	// Plasma Expected Cores Required
	var plasmaExpectedCores float64
	if arrayIndexSizeOfEachElement == 0 {
		plasmaExpectedCores = (float64(mutationRatePerSecond) / float64(mutationIngestThroughputPerCore)) +
			(float64(avgIndexScansPerSecond) / float64(scanThroughputPerCore))
	} else {
		plasmaExpectedCores = ((float64(mutationRatePerSecond) * float64(arrayLength)) / float64(mutationIngestThroughputPerCore)) +
			(float64(avgIndexScansPerSecond) / float64(scanThroughputPerCore))
	}

	// Total Required CPU
	totalRequiredCPU := plasmaExpectedCores * 1.2

	return totalRequiredCPU
}

// calculateIndexDisk computes the disk space required for indexing.
func calculateIndexDisk(dataset models.Dataset) float64 {
	// Constants
	const documentKeyIDSize = 0
	const totalSecondaryBytes = 0
	const snappyCompression = 0.8
	const fragmentation = 0.3

	// Number of documents in index
	numDocsIndex := float64(dataset.NoOfDocuments) * float64(dataset.PercentIndexesOfDataset) / 100

	// Plasma Disk Size 
	plasmaDiskSize := (((numDocsIndex * 2 / 400) * ((documentKeyIDSize + totalSecondaryBytes + 56) * 4)) +
		((documentKeyIDSize + totalSecondaryBytes + 16) * numDocsIndex * 2)) * 2

	// Plasma Compression
	plasmaCompression := plasmaDiskSize * snappyCompression

	// Plasma Fragmentation
	plasmaFragmentation := plasmaDiskSize * fragmentation

	// Plasma Expected Max Disk Usage
	plasmaExpectedMaxDiskUsage := plasmaFragmentation + plasmaCompression

	// Plasma 20% DGM Recommended Disk Quota
	plasmaDGMRecommendedDiskQuota := plasmaExpectedMaxDiskUsage * 1.3
	
	// Convert to GB
	plasmaDGMRecommendedDiskQuota = plasmaDGMRecommendedDiskQuota / 1024 / 1024 /1024

	return plasmaDGMRecommendedDiskQuota
}

// calculateIndexDiskIO computes the disk I/O requirement for the Index service.
func calculateIndexDiskIO() float64 {
	return 0
}
