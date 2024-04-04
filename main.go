package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const FILENAME = "task-sync.data"
const TIME_BETWEEN_SYNCS = 300 // In seconds

func mustFilename() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	file := fmt.Sprintf("%s/.task/%s", dir, FILENAME)
	return file
}

// SaveLastExecutedTime saves the lastExecutedTime to a file
func saveLastExecutedTime(lastExecutedTime time.Time) error {
	// Convert time to string
	timeStr := lastExecutedTime.Format(time.RFC3339)

	// Write time string to file
	return os.WriteFile(mustFilename(), []byte(timeStr), 0644)
}

func lastExecTime() (time.Time, error) {
	// Read time string from file
	timeStr, err := os.ReadFile(mustFilename())
	if err != nil {
		return time.Time{}, err
	}

	// Parse time string into time.Time object
	lastExecutedTime, err := time.Parse(time.RFC3339, string(timeStr))
	if err != nil {
		return time.Time{}, err
	}

	return lastExecutedTime, nil
}

func main() {
	// Read the last executed time
	lastTime, err := lastExecTime()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if we already saved recently so that we don't run into a deadlock
	if time.Since(lastTime) < TIME_BETWEEN_SYNCS*time.Second {
		// fmt.Println("Not Syncing")
		os.Exit(0)
	}

	// Save the current time
	err = saveLastExecutedTime(time.Now())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run the sync command
	cmd := exec.Command("task", "sync")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output the done command
	fmt.Println("Syncing done")
	os.Exit(0)
}
