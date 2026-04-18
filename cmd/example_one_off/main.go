package main

import (
	"encoding/json"
	"fmt"
	"log"

	gogpu "github.com/Henry-GongZY/GoGPU"
)

func main() {
	// 1. Initialize the monitor
	monitor := gogpu.NewMonitor()
	defer monitor.Close() // Ensure resources are freed

	// 2. Discover available GPUs
	infos, err := monitor.GetGPUs()
	if err != nil {
		log.Fatalf("Failed to get GPU list: %v", err)
	}

	if len(infos) == 0 {
		fmt.Println("No GPUs detected on this system.")
		return
	}

	fmt.Println("=== System GPU List ===")
	for i, info := range infos {
		fmt.Printf("[%d] Vendor: %s, Model: %s, Mem: %d MB\n", 
			i, info.Vendor, info.Model, info.TotalMemory/1024/1024)
	}

	fmt.Println("\n=== One-off Query (GetStatusOnce) ===")
	// 3. Perform a single fetch for each GPU
	for i := range infos {
		fmt.Printf("--- Status for GPU [%d] ---\n", i)
		
		status, err := monitor.GetStatusOnce(i)
		if err != nil {
			fmt.Printf("Failed to get status for GPU %d: %v\n", i, err)
			continue
		}
		
		// Print the status beautifully as JSON
		bytes, _ := json.MarshalIndent(status, "", "  ")
		fmt.Printf("%s\n", string(bytes))
	}
}
