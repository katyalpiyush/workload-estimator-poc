package models

// Instance represents an available instance type with vCPU and RAM
type Instance struct {
	VCPU int `json:"vcpu"`
	RAM  int `json:"ram"`
}

// Instances holds all the available instances
var Instances = []Instance{
	{VCPU: 4, RAM: 16},
	{VCPU: 4, RAM: 32},
	{VCPU: 8, RAM: 16},
	{VCPU: 8, RAM: 32},
	{VCPU: 8, RAM: 64},
	{VCPU: 16, RAM: 32},
	{VCPU: 16, RAM: 64},
	{VCPU: 16, RAM: 128},
	{VCPU: 32, RAM: 128},
	{VCPU: 32, RAM: 256},
	{VCPU: 36, RAM: 72},
	{VCPU: 48, RAM: 96},
	{VCPU: 48, RAM: 192},
	{VCPU: 48, RAM: 384},
	{VCPU: 64, RAM: 256},
	{VCPU: 64, RAM: 512},
	{VCPU: 72, RAM: 144},
	{VCPU: 96, RAM: 192},
	{VCPU: 96, RAM: 384},
	{VCPU: 96, RAM: 768},
}