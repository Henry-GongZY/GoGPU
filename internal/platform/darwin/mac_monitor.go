//go:build darwin
// +build darwin

package darwin

import "fmt"

// GPUInfo stores basic static hardware information (local to darwin platform package)
type GPUInfo struct {
	Vendor      string
	Model       string
	TotalMemory uint64
}

// GPUStatus stores real-time dynamic metrics (local to darwin platform package)
type GPUStatus struct {
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

// MacMonitor reads GPU metrics on macOS via system_profiler or IOKit
type MacMonitor struct{}

func NewMacMonitor() *MacMonitor {
	return &MacMonitor{}
}

func (m *MacMonitor) GetGPUs() ([]GPUInfo, error) {
	// TODO: Use "system_profiler SPDisplaysDataType" or CGO IOKit to enumerate GPUs
	var infos []GPUInfo
	return infos, fmt.Errorf("macOS GetGPUs not implemented yet")
}

func (m *MacMonitor) GetStatus(index int) (GPUStatus, error) {
	// TODO: For Intel/AMD GPUs, use CGO IOKit; for Apple Silicon, parse powermetrics output
	return GPUStatus{}, fmt.Errorf("macOS GetStatus not implemented yet")
}

func (m *MacMonitor) Close() error {
	return nil
}
