package calculator

import (
	"math"
	"workload-estimator-poc/models"
)

// RAM is split equally among services (excluding query) on capella
func CalculateRAM(ramList []float64, services []string, nodes int64) float64 {

	noOfServicesInGroup := 0
	for _, service := range services {
		if service != "query" {
			noOfServicesInGroup++
		}
	}
	if noOfServicesInGroup == 0 {
		noOfServicesInGroup = 1
	}

	var osMemoryReserved float64 = 0.2
	var ramAvailable float64 = 0.0
	for _, ram := range ramList {
		var ramAvailableForService float64 = ram / float64(nodes)
		if ramAvailable < ramAvailableForService {
			ramAvailable = ramAvailableForService
		}
	}
	var ramHardwareWithoutOS float64 = ramAvailable * float64(noOfServicesInGroup)

	var ramHardware float64 = ramHardwareWithoutOS / (1 - osMemoryReserved)

	return ramHardware
}

// CalculateCPU takes a list of CPU values and calculates the total CPU hardware required per node
func CalculateCPU(cpuList []float64, nodes int64) float64 {
	if nodes == 0 {
		return 0
	}
	var totalCPU float64
	for _, cpu := range cpuList {
		totalCPU += cpu
	}
	totalCPU /= float64(nodes)
	return math.Ceil(totalCPU)
}

// CalculateDisk takes a list of Disk values and calculates the total Disk storage based on disk type
func CalculateDisk(diskList []float64, nodes int64) float64 {
	if nodes == 0 {
		return 0
	}
	var totalDisk float64
	for _, disk := range diskList {
		totalDisk += disk
	}
	totalDisk /= float64(nodes)
	totalDisk = math.Ceil(totalDisk)
	if totalDisk < 50 {
		totalDisk = 50
	}
	if totalDisk > 16000 {
		totalDisk = 16000
	}
	return totalDisk
}

// CalculateDiskIO takes a list of Disk I/O values and calculates the total Disk I/O
func CalculateDiskIO(diskIOList []float64, nodes int64, diskType string) float64 {
	if nodes == 0 {
		return 0
	}
	var totalDiskIO float64
	for _, diskIO := range diskIOList {
		totalDiskIO += diskIO
	}
	totalDiskIO /= float64(nodes)
	// Adjust totalDiskSpace based on DiskType
	switch diskType {
	case "gp3":
		if totalDiskIO > 16000 {
			totalDiskIO = 16000
		}
	case "io2":
		if totalDiskIO > 64000 {
			totalDiskIO = 64000
		}
	}
	if(totalDiskIO < 3000){
		totalDiskIO = 3000
	}
	return totalDiskIO
}

// findClosestInstance selects the closest instance based on the calculated CPU and RAM
func findClosestInstance(totalCPU, totalRAM float64) models.Instance {
	for _, instance := range models.Instances {
		if float64(instance.VCPU) >= totalCPU && float64(instance.RAM) >= totalRAM {
			return instance
		}
	}
	// If no such instance is found, return the last instance (maximum available)
	// This case is just a fallback, usually shouldn't be needed due to how instances are sorted.
	return models.Instances[len(models.Instances)-1]
}
