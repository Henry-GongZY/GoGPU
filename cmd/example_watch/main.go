package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gogpu "github.com/Henry-GongZY/GoGPU"
)

func main() {
	// 1. Initialize the monitor
	monitor := gogpu.NewMonitor()
	defer monitor.Close() // Ensure resources are freed

	fmt.Println("Starting continuous monitoring...")
	fmt.Println("Press Ctrl+C or wait 5 seconds to stop.")
	fmt.Println("---------------------------------------")

	// 2. Create a context that times out after 5 seconds
	// In a real server/agent app, this might be context.Background() 
	// cancelled when the service stops.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Start watching all GPUs. This runs entirely in the background.
	// It will push an event to the eventsChan every 1 second.
	eventsChan := monitor.GetAllStatus(ctx, 1*time.Second)

	// 4. Consume the stream
	for events := range eventsChan {
		fmt.Printf("\n[Timestamp: %s]\n", time.Now().Format("15:04:05"))
		
		for _, ev := range events {
			// Check if this specific GPU threw an error
			if ev.Err != nil {
				fmt.Printf("GPU [%d] failed to read: %v\n", ev.Index, ev.Err)
				continue
			}

			// Print complete status via JSON
			bytes, _ := json.MarshalIndent(ev.Status, "", "  ")
			fmt.Printf("GPU [%d] Status:\n%s\n", ev.Index, string(bytes))
		}
	}

	fmt.Println("\n---------------------------------------")
	fmt.Println("Context cancelled, monitoring gracefully stopped.")
}
