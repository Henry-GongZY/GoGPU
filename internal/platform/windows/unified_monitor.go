package windows

import (
	"fmt"
	"strings"

	"github.com/Henry-GongZY/GoGPU/internal/vendors/nvidia"
)

type UnifiedMonitor struct {
	wmiMonitor *WMIMonitor
	nvMonitor  *nvidia.NVGPUMonitor
	gpuCache   []GPUInfo
	isNV       []bool // maps to gpuCache index, indicates if we handle it via NVML
}

// NewUnifiedMonitor provides the main entry for internal callers.
func NewUnifiedMonitor() *UnifiedMonitor {
	m := &UnifiedMonitor{
		wmiMonitor: NewWMIMonitor(),
	}
	
	// Attempt to initialize NVIDIA NVML monitor
	nv, err := nvidia.NewNVGPUMonitor("")
	if err == nil {
		m.nvMonitor = nv
	} else {
		fmt.Printf("[DEBUG] NVML initialization failed: %v\n", err)
	}
	
	return m
}

func (m *UnifiedMonitor) GetGPUs() ([]GPUInfo, error) {
	// Always use WMI to scan all physical/virtual GPUs to maintain a stable order
	wmiInfos, err := m.wmiMonitor.GetGPUs()
	if err != nil {
		return nil, fmt.Errorf("UnifiedMonitor failed to probe hardware list: %v", err)
	}
	
	m.gpuCache = wmiInfos
	m.isNV = make([]bool, len(wmiInfos))

	var nvInfos []nvidia.NVGPUInfo
	if m.nvMonitor != nil {
		nvInfos, _ = m.nvMonitor.GetGPUs()
	}

	nvIndex := 0
	for i, info := range wmiInfos {
		vendor := strings.ToLower(info.Vendor)
		model := strings.ToLower(info.Model)
		
		// If WMI indicates NVIDIA, map this index to NVML
		if strings.Contains(vendor, "nvidia") || strings.Contains(model, "nvidia") || strings.Contains(model, "rtx") || strings.Contains(model, "gtx") {
			m.isNV[i] = true
			
			// WMI has a 4GB adapter RAM limit bug. If NVML succeeded, override the memory with accurate NVML reading!
			if nvIndex < len(nvInfos) {
				m.gpuCache[i].Vendor = "NVIDIA"
				m.gpuCache[i].Model = nvInfos[nvIndex].Model
				m.gpuCache[i].TotalMemory = nvInfos[nvIndex].TotalMemory
				nvIndex++
			}
		}
	}

	return m.gpuCache, nil
}

func (m *UnifiedMonitor) GetStatus(index int) (GPUStatus, error) {
	if index < 0 || index >= len(m.gpuCache) {
		return GPUStatus{}, fmt.Errorf("invalid GPU index: %d", index)
	}

	// If we identified this as an NVIDIA GPU and NVML initialized successfully
	if m.isNV[index] && m.nvMonitor != nil {
		// Calculate the corresponding NVML index from the WMI index
		nvIndex := 0
		for i := 0; i < index; i++ {
			if m.isNV[i] {
				nvIndex++
			}
		}
		
		nvStatus, err := m.nvMonitor.GetStatus(nvIndex)
		if err == nil {
			return GPUStatus{
				Temperature:  nvStatus.Temperature,
				CoreClock:    nvStatus.CoreClock,
				MemoryClock:  nvStatus.MemoryClock,
				MemoryUsed:   nvStatus.MemoryUsed,
				PowerDraw:    nvStatus.PowerDraw,
				FanSpeed:     nvStatus.FanSpeed,
				Usage3D:      nvStatus.Usage3D,
				UsageCompute: nvStatus.UsageCompute,
				UsageEncoder: nvStatus.UsageEncoder,
				UsageDecoder: nvStatus.UsageDecoder,
			}, nil
		}
		// If NVML encounters an unexpected error, silently fallback to WMI querying
	}

	// For specific hardware without hooks (e.g. Oray virtual cards, AMD, Intel), use WMI.
	return m.wmiMonitor.GetStatus(index)
}

func (m *UnifiedMonitor) Close() error {
	m.wmiMonitor.Close()
	if m.nvMonitor != nil {
		m.nvMonitor.Close()
	}
	return nil
}
