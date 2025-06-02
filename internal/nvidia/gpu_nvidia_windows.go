package nvidia

import (
	"fmt"
	"github.com/Henry-GongZY/GoGPU/internal/nvidia/nvml-go"
)

// GPUMonitor GPU监控器
type NVGPUMonitor struct {
	nvml        *nvml.API
	devices     []nvml.Device
	devicecount uint32
}

// GPUInfo GPU信息结构体
type GPUInfo struct {
	DecoderUtilization uint32
	SamplingPeriodUs   uint32
	Temperature        uint32
}

// NewGPUMonitor 创建新的GPU监控器实例
func NewNVGPUMonitor(dllPath string) (*NVGPUMonitor, error) {
	monitor := &NVGPUMonitor{}

	// 加载NVML库
	n, err := nvml.New(dllPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load NVML: %w", err)
	}
	monitor.nvml = n

	// 初始化NVML
	if err := monitor.nvml.Init(); err != nil {
		err := monitor.nvml.Shutdown()
		if err != nil {
			return nil, err
		} // 清理资源
		return nil, fmt.Errorf("failed to initialize NVML: %w", err)
	}

	// 获得GPU数量
	count, err := monitor.nvml.DeviceGetCount()
	if err != nil {
		return nil, fmt.Errorf("Failed to get device count: %v", err)
	}
	monitor.devicecount = count

	// 设置所有GPU设备句柄
	var device nvml.Device
	for i := 0; i <= int(count); i++ {
		device, err = monitor.nvml.DeviceGetHandleByIndex(uint32(i))
		if err != nil {
			return nil, fmt.Errorf("Failed to get device %d: %v", i, err)
		}
		monitor.devices = append(monitor.devices, device)
	}
	return monitor, nil
}

// DeviceCount 查找本机所拥有的 Nvidia GPU 的数量
func (m *NVGPUMonitor) DeviceCount() uint32 {
	return uint32(len(m.devices))
}

// GetGPUUtilization 获得GPU利用率
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

// Close 关闭并清理资源
func (m *NVGPUMonitor) Close() error {
	if m.nvml != nil {
		if err := m.nvml.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown NVML: %w", err)
		}
	}
	return nil
}
