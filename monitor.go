package gogpu

// GPUInfo stores basic static hardware information
type GPUInfo struct {
	Vendor      string
	Model       string
	TotalMemory uint64 // Total VRAM in Bytes
}

// GPUStatus stores real-time dynamic metrics (nil means unsupported or unavailable)
type GPUStatus struct {
	// Basic Condition
	Temperature *uint   // GPU Core Temperature (Celsius)
	CoreClock   *uint   // GPU Core Clock (MHz)
	MemoryClock *uint   // VRAM Clock (MHz)
	MemoryUsed  *uint64 // Used VRAM (Bytes)
	PowerDraw   *uint   // Current Power Draw (Watts)
	FanSpeed    *uint   // Fan Speed Percentage (%)

	// Engine Utilization (0~100%)
	Usage3D      *uint // 3D/Render Engine Utilization
	UsageCompute *uint // Compute/CUDA Engine Utilization
	UsageEncoder *uint // Video Encoder (NVENC/VCE) Utilization
	UsageDecoder *uint // Video Decoder (NVDEC/UVD) Utilization
}

// GPUMonitor defines the unified monitoring interface
type GPUMonitor interface {
	GetGPUs() ([]GPUInfo, error)
	GetStatus(index int) (GPUStatus, error)
	Close() error
}
