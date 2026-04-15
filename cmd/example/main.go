package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	gogpu "github.com/Henry-GongZY/GoGPU"
)

func main() {
	// We use the UnifiedMonitor directly on Windows via the standard interface
	var monitor gogpu.GPUMonitor = gogpu.NewMonitor()

	// 1. Get static Info
	infos, err := monitor.GetGPUs()
	if err != nil {
		log.Fatalf("Failed to get GPU list: %v", err)
	}

	fmt.Println("=== System GPU List ===")
	for i, info := range infos {
		fmt.Printf("[%d] Vendor: %s, Model: %s, Mem: %d MB\n", i, info.Vendor, info.Model, info.TotalMemory/1024/1024)
	}

	if len(infos) == 0 {
		fmt.Println("No GPUs detected.")
		return
	}

	fmt.Println("\n=== GPU Real-time Status ===")
	for i := range infos {
		fmt.Printf("\n--- Status for GPU [%d] ---\n", i)
		status, err := monitor.GetStatus(i)
		if err != nil {
			fmt.Printf("Failed to get status: %v\n", err)
			continue
		}
		bytes, _ := json.MarshalIndent(status, "", "  ")
		fmt.Printf("%s\n", string(bytes))
	}

	// Clean up resources
	monitor.Close()
	time.Sleep(1 * time.Second)
}
