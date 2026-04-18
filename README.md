<div align="center">
  <img src="gogpu.png" alt="GoGPU Logo" width="600">
</div>

> [中文文档](README_zh.md)

# GoGPU

A lightweight, cross-platform Go library for GPU performance monitoring. GoGPU provides a unified interface to query real-time GPU metrics across Windows, Linux, and macOS — automatically leveraging the best available backend for each platform and vendor.

## Features

- **Unified API** — one interface, all platforms
- **Windows** — WMI for device discovery + NVML for high-fidelity NVIDIA metrics
- **Linux** *(planned)* — sysfs / DRM for vendor-agnostic monitoring
- **macOS** *(planned)* — IOKit / `powermetrics` for Apple Silicon and discrete GPUs
- **Virtual GPU filtering** — automatically skips non-physical adapters (e.g. Indirect Display Drivers)
- **Vendor SDK fallback** — NVML is preferred over WMI when available; graceful degradation otherwise
- **Nil-safe metrics** — every optional field uses a pointer; `nil` means "not supported on this hardware", never "zero"

## Supported Metrics

| Metric | Field | Windows (NVIDIA) | Windows (other) | Linux | macOS |
|---|---|---|---|---|---|
| GPU Core Temperature | `Temperature` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Core Clock | `CoreClock` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Memory Clock | `MemoryClock` | ✅ NVML | ❌ | 🔜 | 🔜 |
| VRAM Used | `MemoryUsed` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Power Draw | `PowerDraw` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Fan Speed | `FanSpeed` | ✅ NVML* | ❌ | 🔜 | 🔜 |
| 3D / Render Utilization | `Usage3D` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Compute / CUDA Utilization | `UsageCompute` | 🔜 | ❌ | 🔜 | 🔜 |
| Encoder Utilization | `UsageEncoder` | ✅ NVML | ❌ | 🔜 | 🔜 |
| Decoder Utilization | `UsageDecoder` | ✅ NVML | ❌ | 🔜 | 🔜 |

> \* Fan speed may be unavailable on laptops where the fan is controlled by OEM firmware (returns `nil`).

## Requirements

- **Go 1.21+**
- **Windows**: NVIDIA driver with `nvml.dll` in `%PATH%` or `System32` (for NVML metrics)
- **Linux**: `/sys/class/drm` access (no additional drivers required for sysfs)
- **macOS**: No additional requirements for `system_profiler`; `powermetrics` requires `sudo`

## Installation

```bash
go get github.com/Henry-GongZY/GoGPU
```

## Quick Start

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

    // List all physical GPUs
    gpus, err := monitor.GetGPUs()
    if err != nil {
        log.Fatal(err)
    }
    for i, gpu := range gpus {
        fmt.Printf("[%d] %s %s — %d MB VRAM\n",
            i, gpu.Vendor, gpu.Model, gpu.TotalMemory/1024/1024)
    }

    // Query real-time status for GPU 0
    status, err := monitor.GetStatus(0)
    if err != nil {
        log.Fatal(err)
    }
    b, _ := json.MarshalIndent(status, "", "  ")
    fmt.Println(string(b))
}
```

**Example output (NVIDIA RTX 4060 Laptop GPU on Windows):**
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

## API Reference

### `gogpu.NewMonitor() GPUMonitor`

Returns a `GPUMonitor` for the current OS. Internally selects the best available backend.

### `GPUMonitor` interface

```go
type GPUMonitor interface {
    GetGPUs() ([]GPUInfo, error)
    GetStatus(index int) (GPUStatus, error)
    Close() error
}
```

### `GPUInfo`

| Field | Type | Description |
|---|---|---|
| `Vendor` | `string` | GPU vendor name (e.g. `"NVIDIA"`, `"AMD"`) |
| `Model` | `string` | GPU model name |
| `TotalMemory` | `uint64` | Total VRAM in bytes |

### `GPUStatus`

All pointer fields — `nil` indicates the metric is not supported on this device or platform.

| Field | Type | Unit |
|---|---|---|
| `Temperature` | `*uint` | °C |
| `CoreClock` | `*uint` | MHz |
| `MemoryClock` | `*uint` | MHz |
| `MemoryUsed` | `*uint64` | Bytes |
| `PowerDraw` | `*uint` | W |
| `FanSpeed` | `*uint` | % |
| `Usage3D` | `*uint` | % |
| `UsageCompute` | `*uint` | % |
| `UsageEncoder` | `*uint` | % |
| `UsageDecoder` | `*uint` | % |

## Architecture

```
github.com/Henry-GongZY/GoGPU
├── monitor.go            # Public API: GPUMonitor interface, GPUInfo, GPUStatus
├── gogpu_windows.go      # Windows entry point (build tag: windows)
│
└── internal/
    ├── platform/
    │   ├── windows/      # WMI device discovery + NVML orchestration
    │   ├── linux/        # sysfs / DRM (planned)
    │   └── darwin/       # IOKit / powermetrics (planned)
    └── vendors/
        ├── nvidia/       # NVML bindings (NVGPUMonitor)
        ├── amd/          # ADL / ROCm (planned)
        └── intel/        # IGSC (planned)
```

## Debugging

Run the bundled example to verify everything works on your machine:

```bash
go run ./cmd/example
```

## Roadmap

- [ ] Linux sysfs implementation (temperature, utilization, clocks)
- [ ] macOS IOKit / powermetrics implementation
- [ ] AMD Windows support via ADL/ADLx
- [ ] Intel Windows/Linux support
- [ ] PDH-based Windows fallback for non-NVIDIA GPUs

## License

[MIT](LICENSE)
