# GoGPU

> [English Documentation](README.md)

一个轻量级、跨平台的 Go 语言 GPU 性能监控库。GoGPU 提供统一的接口，用于在 Windows、Linux 和 macOS 上查询实时 GPU 指标——自动为每个平台和厂商选择最佳的底层后端。

## 特性

- **统一 API** — 一套接口，全平台通用
- **Windows** — WMI 负责设备发现，NVML 提供 NVIDIA 显卡的高精度指标
- **Linux** *(计划中)* — sysfs / DRM，跨厂商通用监控
- **macOS** *(计划中)* — IOKit / `powermetrics`，支持 Apple Silicon 和独立显卡
- **虚拟 GPU 过滤** — 自动跳过非物理适配器（如远程显示驱动 Indirect Display Drivers）
- **厂商 SDK 优先** — 当 NVML 可用时优先使用，失败时自动降级至系统 API
- **空安全指标** — 每个可选字段均使用指针类型；`nil` 表示"此硬件不支持该指标"，而不是"值为零"

## 支持的监控指标

| 指标 | 字段 | Windows (NVIDIA) | Windows (其他) | Linux | macOS |
|---|---|---|---|---|---|
| GPU 核心温度 | `Temperature` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 核心频率 | `CoreClock` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 显存频率 | `MemoryClock` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 显存用量 | `MemoryUsed` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 功耗 | `PowerDraw` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 风扇转速 | `FanSpeed` | ✅ NVML* | ❌ | 🔜 | 🔜 |
| 3D / 渲染引擎占用率 | `Usage3D` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 计算 / CUDA 占用率 | `UsageCompute` | 🔜 | ❌ | 🔜 | 🔜 |
| 编码器占用率 | `UsageEncoder` | ✅ NVML | ❌ | 🔜 | 🔜 |
| 解码器占用率 | `UsageDecoder` | ✅ NVML | ❌ | 🔜 | 🔜 |

> \* 笔记本电脑上风扇通常由 OEM 固件统一管控，无法经由 NVML 获取，此时该字段返回 `nil`。

## 环境要求

- **Go 1.21+**
- **Windows**：NVIDIA 驱动已安装，且 `nvml.dll` 位于 `%PATH%` 或 `System32` 目录下（用于 NVML 指标）
- **Linux**：可访问 `/sys/class/drm`（无需额外驱动）
- **macOS**：`system_profiler` 无需额外配置；`powermetrics` 需要 `sudo` 权限

## 安装

```bash
go get github.com/Henry-GongZY/GoGPU
```

## 快速开始

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    gogpu "github.com/Henry-GongZY/GoGPU"
)

func main() {
    monitor := gogpu.NewMonitor()
    defer monitor.Close()

    // 列出所有物理 GPU
    gpus, err := monitor.GetGPUs()
    if err != nil {
        log.Fatal(err)
    }
    for i, gpu := range gpus {
        fmt.Printf("[%d] %s %s — %d MB 显存\n",
            i, gpu.Vendor, gpu.Model, gpu.TotalMemory/1024/1024)
    }

    // 查询第 0 张 GPU 的实时状态
    status, err := monitor.GetStatus(0)
    if err != nil {
        log.Fatal(err)
    }
    b, _ := json.MarshalIndent(status, "", "  ")
    fmt.Println(string(b))
}
```

**示例输出（Windows 上的 NVIDIA RTX 4060 Laptop GPU）：**
```json
{
  "Temperature": 47,
  "CoreClock": 420,
  "MemoryClock": 6001,
  "MemoryUsed": 2496270336,
  "PowerDraw": 13,
  "FanSpeed": null,
  "Usage3D": 22,
  "UsageCompute": null,
  "UsageEncoder": 0,
  "UsageDecoder": 1
}
```

## API 参考

### `gogpu.NewMonitor() GPUMonitor`

返回当前操作系统对应的 `GPUMonitor` 实例，内部自动选择最佳可用后端。

### `GPUMonitor` 接口

```go
type GPUMonitor interface {
    GetGPUs() ([]GPUInfo, error)
    GetStatus(index int) (GPUStatus, error)
    Close() error
}
```

### `GPUInfo`

| 字段 | 类型 | 说明 |
|---|---|---|
| `Vendor` | `string` | 厂商名称（如 `"NVIDIA"`、`"AMD"`） |
| `Model` | `string` | 显卡型号 |
| `TotalMemory` | `uint64` | 显存总量（字节） |

### `GPUStatus`

所有字段均为指针类型——`nil` 表示该设备或平台不支持此指标。

| 字段 | 类型 | 单位 |
|---|---|---|
| `Temperature` | `*uint` | °C |
| `CoreClock` | `*uint` | MHz |
| `MemoryClock` | `*uint` | MHz |
| `MemoryUsed` | `*uint64` | 字节 |
| `PowerDraw` | `*uint` | W |
| `FanSpeed` | `*uint` | % |
| `Usage3D` | `*uint` | % |
| `UsageCompute` | `*uint` | % |
| `UsageEncoder` | `*uint` | % |
| `UsageDecoder` | `*uint` | % |

## 架构设计

```
github.com/Henry-GongZY/GoGPU
├── monitor.go            # 公共 API：GPUMonitor 接口、GPUInfo、GPUStatus 类型定义
├── gogpu_windows.go      # Windows 入口（构建标签：windows）
│
└── internal/
    ├── platform/
    │   ├── windows/      # WMI 设备发现 + NVML 聚合编排
    │   ├── linux/        # sysfs / DRM（计划中）
    │   └── darwin/       # IOKit / powermetrics（计划中）
    └── vendors/
        ├── nvidia/       # NVML 绑定（NVGPUMonitor）
        ├── amd/          # ADL / ROCm（计划中）
        └── intel/        # IGSC（计划中）
```

### 设计亮点

- **防腐层（Anti-corruption Layer）**：`internal/` 下的所有代码对外部调用者完全不可见，彻底避免实现细节泄漏。
- **无循环依赖**：`vendors/nvidia` 定义了自己独立的 `NVGPUInfo` / `NVGPUStatus`，不反向依赖 `platform` 层，从根源杜绝 import cycle。
- **构建标签隔离**：`gogpu_windows.go` 带有 `//go:build windows` 标签，不同平台的编译完全互不干扰。

## 调试

运行内置的示例程序，验证本机环境是否一切正常：

```bash
go run ./cmd/example
```

## 路线图

- [ ] Linux sysfs 实现（温度、占用率、频率）
- [ ] macOS IOKit / powermetrics 实现
- [ ] Windows AMD 支持（ADL / ADLx）
- [ ] Intel 跨平台支持
- [ ] 非 NVIDIA 显卡的 Windows PDH API 降级方案

## 开源协议

[MIT](LICENSE)
