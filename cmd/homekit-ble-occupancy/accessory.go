package main

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

// OccupancySensor implements an occupancy sensor accessory.
type OccupancySensor struct {
	*accessory.Accessory
	OccupancySensor *service.OccupancySensor
}

// NewOccupancySensor returns an occupancy sensor accessory with one occupancy
// sensor service.
func NewOccupancySensor(info accessory.Info) *OccupancySensor {
	acc := &OccupancySensor{
		Accessory:       accessory.New(info, accessory.TypeSensor),
		OccupancySensor: service.NewOccupancySensor(),
	}

	acc.AddService(acc.OccupancySensor.Service)

	return acc
}
