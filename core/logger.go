package core

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logMu sync.Mutex
)

// InitLogger sets up the directory and redirects standard output streams
func InitLogger(path string) {
	_ = os.MkdirAll("logs", 0755)

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Logger Initialization Failed: %v\n", err)
		return
	}

	// Redirecting standard output so fmt.Printf everywhere goes to the log
	os.Stdout = logFile
	os.Stderr = logFile
	log.SetOutput(logFile)
}

// Log is a package-level helper for standardized Rictus formatting
func (e *Engine) Log(message string) {
	logMu.Lock()
	defer logMu.Unlock()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] 💀 RICTUS: %s\n", timestamp, message)
}
