package windows

import (
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

type WMIMonitor struct {
}

func NewWMIMonitor() *WMIMonitor {
	return &WMIMonitor{}
}

// Win32_VideoController WMI structure
type win32_VideoController struct {
	Name       string
	AdapterRAM uint32
}

// win32_PerfFormattedData_GPUEngine WMI structure
type win32_PerfFormattedData_GPUEngine struct {
	Name                  string
	UtilizationPercentage uint64
}

func (m *WMIMonitor) GetGPUs() ([]GPUInfo, error) {
	var vcs []win32_VideoController
	q := wmi.CreateQuery(&vcs, "")
	err := wmi.Query(q, &vcs)
	if err != nil {
		return nil, fmt.Errorf("WMI Win32_VideoController query failed: %v", err)
	}

	var infos []GPUInfo
	for _, vc := range vcs {
		// Filter out virtual GPUs like Indirect Display Drivers (IDD) which usually have 0 AdapterRAM
		if vc.AdapterRAM == 0 || strings.Contains(strings.ToLower(vc.Name), "idddriver") {
			continue
		}

		vendor, model := parseVendorAndModel(vc.Name)
		infos = append(infos, GPUInfo{
			Vendor:      vendor,
			Model:       model,
			TotalMemory: uint64(vc.AdapterRAM), // Note: AdapterRAM might overflow 4GB limits in WMI
		})
	}
	return infos, nil
}

func (m *WMIMonitor) GetStatus(index int) (GPUStatus, error) {
	// This is a skeleton implementation. In reality, it should filter by indexing and Names.
	var engines []win32_PerfFormattedData_GPUEngine
	q := wmi.CreateQuery(&engines, "")
	err := wmi.Query(q, &engines)
	if err != nil {
		return GPUStatus{}, fmt.Errorf("WMI GPUEngine query failed: %v", err)
	}

	status := GPUStatus{}
	var total3D uint64
	var totalCompute uint64
	var totalEncode uint64
	var totalDecode uint64

	for _, eng := range engines {
		name := strings.ToLower(eng.Name)
		if strings.Contains(name, "3d") {
			total3D += eng.UtilizationPercentage
		} else if strings.Contains(name, "compute") || strings.Contains(name, "cuda") {
			totalCompute += eng.UtilizationPercentage
		} else if strings.Contains(name, "videoencode") {
			totalEncode += eng.UtilizationPercentage
		} else if strings.Contains(name, "videodecode") {
			totalDecode += eng.UtilizationPercentage
		}
	}

	usage3d := uint(total3D)
	usageCompute := uint(totalCompute)
	usageEncode := uint(totalEncode)
	usageDecode := uint(totalDecode)

	status.Usage3D = &usage3d
	status.UsageCompute = &usageCompute
	status.UsageEncoder = &usageEncode
	status.UsageDecoder = &usageDecode

	return status, nil
}

func (m *WMIMonitor) Close() error {
	return nil
}

// parseVendorAndModel basic splitting for Vendor and Model
func parseVendorAndModel(name string) (string, string) {
	name = strings.TrimSpace(name)
	parts := strings.SplitN(name, " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "Unknown", name
}
