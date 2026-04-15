//go:build linux
// +build linux

package linux

import "fmt"

// GPUInfo stores basic static hardware information (local to linux platform package)
type GPUInfo struct {
	Vendor      string
	Model       string
	TotalMemory uint64
}

// GPUStatus stores real-time dynamic metrics (local to linux platform package)
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

// SysfsMonitor reads GPU metrics from the Linux sysfs virtual filesystem
type SysfsMonitor struct{}

func NewSysfsMonitor() *SysfsMonitor {
	return &SysfsMonitor{}
}

func (m *SysfsMonitor) GetGPUs() ([]GPUInfo, error) {
	// TODO: Scan /sys/class/drm/card<N>/device for GPU list
	// Read vendor/device IDs from /sys/class/drm/card<N>/device/vendor and map to names
	var infos []GPUInfo
	return infos, fmt.Errorf("linux sysfs GetGPUs not implemented yet")
}

func (m *SysfsMonitor) GetStatus(index int) (GPUStatus, error) {
	// TODO: Read primary utilization from /sys/class/drm/card<N>/device/gpu_busy_percent
	// TODO: Read temperature from /sys/class/hwmon/hwmon<N>/temp1_input
	// TODO: Read frequency from /sys/class/hwmon/hwmon<N>/freq1_input
	return GPUStatus{}, fmt.Errorf("linux sysfs GetStatus not implemented yet")
}

func (m *SysfsMonitor) Close() error {
	return nil
}
