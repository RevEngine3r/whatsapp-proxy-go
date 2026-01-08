package main

import (
	"fmt"
	"os"
)

// Version information (set via ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	fmt.Printf("WhatsApp Proxy Go v%s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Println("\n🚧 In Development - Step 1 Complete")
	fmt.Println("\nProject structure initialized.")
	fmt.Println("Next: Configuration management implementation")
	os.Exit(0)
}
