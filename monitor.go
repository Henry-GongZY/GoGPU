package gogpu

import (
	"context"
	"time"

	"github.com/Henry-GongZY/GoGPU/internal/watcher"
)

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

// StatusEvent represents a single polling event for a specific GPU
type StatusEvent struct {
	Index  int
	Status GPUStatus
	Err    error
}

// platformBackend is implemented by OS-specific files (e.g. gogpu_windows.go)
type platformBackend interface {
	GetGPUs() ([]GPUInfo, error)
	GetStatus(index int) (GPUStatus, error)
	Close() error
}

// GPUMonitor provides the main unified monitoring API.
type GPUMonitor struct {
	backend platformBackend
}

// GetGPUs returns the discovered physical GPUs.
func (m *GPUMonitor) GetGPUs() ([]GPUInfo, error) {
	return m.backend.GetGPUs()
}

// GetStatusOnce performs a one-off query for the current status of the specified GPU.
func (m *GPUMonitor) GetStatusOnce(index int) (GPUStatus, error) {
	return m.backend.GetStatus(index)
}

// GetStatus starts a background goroutine that polls the specified GPU status.
// It returns a read-only channel emitting StatusEvents.
// Polling automatically stops and the channel closes when ctx is canceled.
func (m *GPUMonitor) GetStatus(ctx context.Context, index int, interval time.Duration) <-chan StatusEvent {
	return watcher.Poll(ctx, interval, func() StatusEvent {
		status, err := m.GetStatusOnce(index)
		return StatusEvent{
			Index:  index,
			Status: status,
			Err:    err,
		}
	})
}

// GetAllStatus starts a background goroutine that polls all detected GPUs.
func (m *GPUMonitor) GetAllStatus(ctx context.Context, interval time.Duration) <-chan []StatusEvent {
	gpus, err := m.GetGPUs()
	if err != nil {
		out := make(chan []StatusEvent, 1)
		out <- []StatusEvent{{Index: -1, Err: err}}
		close(out)
		return out
	}

	return watcher.Poll(ctx, interval, func() []StatusEvent {
		events := make([]StatusEvent, len(gpus))
		for i := range gpus {
			status, err := m.GetStatusOnce(i)
			events[i] = StatusEvent{
				Index:  i,
				Status: status,
				Err:    err,
			}
		}
		return events
	})
}

// Close releases any underlying hardware resources (e.g. freeing NVML API wrappers).
func (m *GPUMonitor) Close() error {
	return m.backend.Close()
}
