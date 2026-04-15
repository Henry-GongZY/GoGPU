//go:build windows
// +build windows

package gogpu

import "github.com/Henry-GongZY/GoGPU/internal/platform/windows"

type windowsMonitorWrapper struct {
	impl *windows.UnifiedMonitor
}

// NewMonitor is the main entry point to initialize GoGPU on Windows.
func NewMonitor() GPUMonitor {
	return &windowsMonitorWrapper{
		impl: windows.NewUnifiedMonitor(),
	}
}

func (w *windowsMonitorWrapper) GetGPUs() ([]GPUInfo, error) {
	winInfos, err := w.impl.GetGPUs()
	if err != nil {
		return nil, err
	}
	infos := make([]GPUInfo, len(winInfos))
	for i, v := range winInfos {
		infos[i] = GPUInfo{
			Vendor:      v.Vendor,
			Model:       v.Model,
			TotalMemory: v.TotalMemory,
		}
	}
	return infos, nil
}

func (w *windowsMonitorWrapper) GetStatus(index int) (GPUStatus, error) {
	winStatus, err := w.impl.GetStatus(index)
	if err != nil {
		return GPUStatus{}, err
	}
	return GPUStatus{
		Temperature:  winStatus.Temperature,
		CoreClock:    winStatus.CoreClock,
		MemoryClock:  winStatus.MemoryClock,
		MemoryUsed:   winStatus.MemoryUsed,
		PowerDraw:    winStatus.PowerDraw,
		FanSpeed:     winStatus.FanSpeed,
		Usage3D:      winStatus.Usage3D,
		UsageCompute: winStatus.UsageCompute,
		UsageEncoder: winStatus.UsageEncoder,
		UsageDecoder: winStatus.UsageDecoder,
	}, nil
}

func (w *windowsMonitorWrapper) Close() error {
	// Add close logic if platform requires
	return nil
}
