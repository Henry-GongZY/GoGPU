package nvidia

import (
	"fmt"

	"github.com/Henry-GongZY/GoGPU/internal/vendors/nvidia/nvml-go"
)

// NVGPUMonitor is the GPU monitor for NVIDIA cards
type NVGPUMonitor struct {
	nvml        *nvml.API
	devices     []nvml.Device
	devicecount uint32
}

// NVGPUInfo stores basic info for NVIDIA devices
type NVGPUInfo struct {
	Vendor      string
	Model       string
	TotalMemory uint64
}

// NVGPUStatus stores status metrics for NVIDIA devices
type NVGPUStatus struct {
	Temperature  *uint
	CoreClock    *uint
	MemoryClock  *uint
	MemoryUsed   *uint64
	PowerDraw    *uint
	FanSpeed     *uint
	Usage3D      *uint
	UsageCompute *uint
	UsageEncoder *uint
	UsageDecoder *uint
}

// NewNVGPUMonitor creates a new NVGPUMonitor instance and initializes NVML
func NewNVGPUMonitor(dllPath string) (*NVGPUMonitor, error) {
	monitor := &NVGPUMonitor{}

	if dllPath == "" {
		dllPath = "nvml.dll"
	}
	// Load NVML Library
	n, err := nvml.New(dllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load NVML: %w", err)
	}
	monitor.nvml = n

	// Initialize NVML
	if err := monitor.nvml.Init(); err != nil {
		err := monitor.nvml.Shutdown()
		if err != nil {
			return nil, err
		} // Clean up resources
		return nil, fmt.Errorf("failed to initialize NVML: %w", err)
	}

	// Get GPU count
	count, err := monitor.nvml.DeviceGetCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get device count: %v", err)
	}
	monitor.devicecount = count

	// Setup all GPU device handles
	var device nvml.Device
	for i := 0; i < int(count); i++ {
		device, err = monitor.nvml.DeviceGetHandleByIndex(uint32(i))
		if err != nil {
			return nil, fmt.Errorf("Failed to get device %d: %v", i, err)
		}
		monitor.devices = append(monitor.devices, device)
	}
	return monitor, nil
}

// DeviceCount gets the number of NVIDIA GPUs on the machine
func (m *NVGPUMonitor) DeviceCount() uint32 {
	return m.devicecount
}

// GetGPUUtilization gets GPU usage
func (m *NVGPUMonitor) GetGPUUtilization(index uint32) (uint32, uint32, error) {
	utilization, err := m.nvml.DeviceGetUtilizationRates(m.devices[index])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get gpu utilization of gpu %d: %w", index, err)
	}

	return utilization.GPU, utilization.Memory, nil
}

// GetEncoderUtilization 获得编码器利用率
func (m *NVGPUMonitor) GetEncoderUtilization(index uint32) (uint32, uint32, error) {
	utilization, samplingPeriodUs, err := m.nvml.DeviceGetEncoderUtilization(m.devices[index])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get encoder utilization of gpu %d: %w", index, err)
	}

	return utilization, samplingPeriodUs, nil
}

// GetDecoderUtilization 获取解码器利用率
func (m *NVGPUMonitor) GetDecoderUtilization(index uint32) (uint32, uint32, error) {
	utilization, samplingPeriodUs, err := m.nvml.DeviceGetDecoderUtilization(m.devices[index])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get decoder utilization of gpu %d: %w", index, err)
	}

	return utilization, samplingPeriodUs, nil
}

// GetTemperature 获取GPU温度
func (m *NVGPUMonitor) GetTemperature(index uint32) (uint32, error) {

	temp, err := m.nvml.DeviceGetTemperature(m.devices[index], nvml.TemperatureGPU)
	if err != nil {
		return 0, fmt.Errorf("failed to get temperature of gpu %d: %w", index, err)
	}

	return temp, nil
}

// GetFanSpeed 获取风扇速度
func (m *NVGPUMonitor) GetFanSpeed(index uint32) (uint32, error) {
	speed, err := m.nvml.DeviceGetFanSpeed(m.devices[index])
	if err != nil {
		return 0, fmt.Errorf("failed to get fan speed of gpu %d: %w", index, err)
	}
	return speed, nil
}

// Close gracefully shutdowns the NVML API
func (m *NVGPUMonitor) Close() error {
	if m.nvml != nil {
		if err := m.nvml.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown NVML: %w", err)
		}
	}
	return nil
}

// GetGPUs returns static hardware configuration
func (m *NVGPUMonitor) GetGPUs() ([]NVGPUInfo, error) {
	var infos []NVGPUInfo
	for i, device := range m.devices {
		name, err := m.nvml.DeviceGetName(device)
		if err != nil {
			name = fmt.Sprintf("NVIDIA GPU %d", i)
		}
		
		var totalMem uint64
		memInfo, err := m.nvml.DeviceGetMemoryInfo(device)
		if err == nil {
			totalMem = memInfo.Total
		}

		infos = append(infos, NVGPUInfo{
			Vendor:      "NVIDIA",
			Model:       name,
			TotalMemory: totalMem,
		})
	}
	return infos, nil
}

// GetStatus returns real-time metrics
func (m *NVGPUMonitor) GetStatus(index int) (NVGPUStatus, error) {
	if index < 0 || index >= len(m.devices) {
		return NVGPUStatus{}, fmt.Errorf("invalid gpu index %d", index)
	}

	status := NVGPUStatus{}

	// Utilization
	util, _, err := m.GetGPUUtilization(uint32(index))
	if err == nil {
		u := uint(util)
		status.Usage3D = &u
	}

	// Encoder
	encUtil, _, err := m.GetEncoderUtilization(uint32(index))
	if err == nil {
		u := uint(encUtil)
		status.UsageEncoder = &u
	}

	// Decoder
	decUtil, _, err := m.GetDecoderUtilization(uint32(index))
	if err == nil {
		u := uint(decUtil)
		status.UsageDecoder = &u
	}

	// Temperature
	temp, err := m.GetTemperature(uint32(index))
	if err == nil {
		t := uint(temp)
		status.Temperature = &t
	}

	// Fan speed
	fan, err := m.GetFanSpeed(uint32(index))
	if err == nil {
		f := uint(fan)
		status.FanSpeed = &f
	}

	// VRAM usage
	memInfo, err := m.nvml.DeviceGetMemoryInfo(m.devices[index])
	if err == nil {
		used := memInfo.Used
		status.MemoryUsed = &used
	}

	// Clocks
	coreClock, err := m.nvml.DeviceGetClockInfo(m.devices[index], nvml.ClockGraphics)
	if err == nil {
		c := uint(coreClock)
		status.CoreClock = &c
	}
	memClock, err := m.nvml.DeviceGetClockInfo(m.devices[index], nvml.ClockMem)
	if err == nil {
		mc := uint(memClock)
		status.MemoryClock = &mc
	}
	
	// Power draw
	power, err := m.nvml.DeviceGetPowerUsage(m.devices[index])
	if err == nil {
		p := uint(power / 1000) // Usually NVML returns milliwatts
		status.PowerDraw = &p
	}

	return status, nil
}
