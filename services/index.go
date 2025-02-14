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

	cpu = calculateIndexCPU(workload, nodes)
	disk = calculateIndexDisk(ram, nodes)
	diskIO = calculateIndexDiskIO(workload, nodes)
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

// calculateIndexCPU computes the CPU requirement for the Index service.
func calculateIndexCPU(workload models.Workload, nodes int) float64 {
	cpu := float64(workload.SQLQueriesPerSec) / 100.0
	return cpu / float64(nodes) // Normalize by number of nodes
}

// calculateIndexDisk computes the disk space required for indexing.
func calculateIndexDisk(ram float64, nodes int) float64 {
	disk := ram * 2 // Disk requirement is typically twice the RAM
	return disk / float64(nodes) // Normalize by number of nodes
}

// calculateIndexDiskIO computes the disk I/O requirement for the Index service.
func calculateIndexDiskIO(workload models.Workload, nodes int) float64 {
	diskIO := float64(workload.SQLQueriesPerSec) * 5
	return diskIO / float64(nodes) // Normalize by number of nodes
}
