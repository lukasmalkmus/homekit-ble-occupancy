package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/axiomhq/pkg/version"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/util"
	"github.com/mdp/qrterminal/v3"
	"github.com/shirou/gopsutil/host"
	"tinygo.org/x/bluetooth"
)

const setupID = "BLEO"

var (
	port    = flag.Int("port", 0, "port on which accessory is reachable")
	pin     = flag.String("pin", "", "pin to use (must be 8 characters of 0-9)")
	storage = flag.String("storage", "./homekit-ble-occupancy", "storage folder to use")
)

// Device contains device metadata.
type Device struct {
	// LastSeen is the time the device was last seen.
	LastSeen time.Time
	// RSSILimit is the upper limit for the device to get tracked.
	RSSILimit int
}

func main() {
	flag.Parse()

	// If no pin is given, create a random one. In any case, validate the pin.
	if *pin == "" {
		var err error
		if *pin, err = generateRandomPin(); err != nil {
			log.Fatal(err)
		}
	}
	humanPin, err := hc.ValidatePin(*pin)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Bluetooth adapter.
	adapter := bluetooth.DefaultAdapter
	if err = adapter.Enable(); err != nil {
		log.Fatal(err)
	}

	// Get the host machines info.
	info, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	// If no arguments are given, scan for Bluetooth devices and print their
	// address and signal strength.
	if flag.NArg() == 0 {
		log.Print("Scanning for devices...")
		if scanErr := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			// if device.LocalName() != "" {
			log.Printf("%s: %s (%d)\n", device.LocalName(), device.Address.String(), device.RSSI*-1)
			// }
		}); scanErr != nil {
			log.Fatal(scanErr)
		}
		return
	}

	// Create the list of devices we want to keep track of. Arguments must be
	// the bluetooth address and the RSSI separated by a plus ('+') character.
	// The RSSI value is positive as opossed to the meaning of RSSI being
	// negatively donated.
	bluetoothDevices := make(map[string]*Device, flag.NArg())
	bluetoothDevicesMtx := new(sync.RWMutex)
	for _, deviceConfig := range flag.Args() {
		parts := strings.Split(deviceConfig, "+")
		if len(parts) != 2 {
			log.Fatal("expected address and positive rssi separated by '+'")
		}

		address := parts[0]
		rssiStr := parts[1]

		rssi, convErr := strconv.Atoi(rssiStr)
		if convErr != nil {
			log.Fatal(convErr)
		}

		bluetoothDevices[address] = &Device{
			RSSILimit: rssi,
		}
	}

	// Create the HomeKit accessory.
	acc := NewOccupancySensor(accessory.Info{
		Name:             "BLE Occupancy Sensor",
		Manufacturer:     "Lukas Malkmus",
		Model:            info.Platform,
		SerialNumber:     info.PlatformVersion,
		FirmwareRevision: version.Revision(),
	})
	occupied := acc.OccupancySensor.OccupancyDetected

	// Setup the config and transport.
	config := hc.Config{
		StoragePath: *storage,
		Port:        strconv.Itoa(*port),
		Pin:         *pin,
		SetupId:     setupID,
	}

	transport, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	// Generate and print the HomeKit setup code on the terminal.
	xhm, err := config.XHMURI(util.SetupFlagIP)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Scan the QR-Code below or enter the setup code %q manaually.\n", humanPin)
	qrterminal.Generate(xhm, qrterminal.M, os.Stdout)

	// Scan for devices and if a device being tracked is found, make sure its
	// signal strength is strong enough before updating its last seen time.
	go func() {
		if err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			bluetoothDevicesMtx.Lock()
			defer bluetoothDevicesMtx.Unlock()

			if meta, ok := bluetoothDevices[device.Address.String()]; ok && (meta.RSSILimit+int(device.RSSI)) > 0 {
				meta.LastSeen = time.Now()
			}
		}); err != nil {
			log.Fatal(err)
		}
	}()

	// Regularly check, if the device has recently been seen and update the
	// occupancy status accordingly.
	go func() {
		for {
			select {
			case <-time.After(time.Second):
			default:
			}

			bluetoothDevicesMtx.RLock()
			for _, meta := range bluetoothDevices {
				if time.Since(meta.LastSeen) < time.Second*10 {
					occupied.SetValue(characteristic.OccupancyDetectedOccupancyDetected)
				} else {
					occupied.SetValue(characteristic.OccupancyDetectedOccupancyNotDetected)
				}
			}
			bluetoothDevicesMtx.RUnlock()
		}
	}()

	// On termination, stop the transport and adapter.
	hc.OnTermination(func() {
		<-transport.Stop()

		if err := adapter.StopScan(); err != nil {
			log.Fatal(err)
		}
	})

	transport.Start()
}

func generateRandomPin() (string, error) {
	const (
		length  = 8
		charset = "0123456789"
	)

	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		ret[i] = charset[num.Int64()]
	}

	return string(ret), nil
}
