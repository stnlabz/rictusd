package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	// Using the existing Block struct for decoding
	"rictusd/modules/blockchain" 
)

const (
	pidFile = "data/rictusd.pid"
	binPath = "./bin/rictusd"
	ledgerPath = "data/blockchain/local_chain.json"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "start":
		startDaemon()
	case "stop":
		stopDaemon()
	case "restart":
		stopDaemon()
		startDaemon()
	case "status":
		checkPID()
	case "inspect":
		inspectLedger() // New inspection handler
	default:
		printHelp()
	}
}

func inspectLedger() {
	file, err := os.Open(ledgerPath)
	if err != nil {
		fmt.Printf("❌ Unable to open ledger: %v. Does the data directory exist?\n", err)
		return
	}
	defer file.Close()

	fmt.Println("📜 --- RICTUSD VERIFIED THREAT HISTORY ---")
	fmt.Printf("%-20s | %-10s | %s\n", "TIMESTAMP", "HASH (8)", "THREAT DATA")
	fmt.Println("--------------------------------------------------------------------------------")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var b blockchain.Block
		if err := json.Unmarshal(scanner.Bytes(), &b); err != nil {
			continue
		}
		
		// Truncate timestamp for readability and show short hash
		ts := b.Timestamp
		if len(ts) > 19 {
			ts = ts[:19]
		}
		
		shortHash := b.Hash
		if len(shortHash) > 8 {
			shortHash = shortHash[:8]
		}

		fmt.Printf("%-20s | %-10s | %v\n", ts, shortHash, b.Data)
	}
}

func startDaemon() {
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		fmt.Printf("❌ Error: %s not found. Run 'make build' first.\n", binPath)
		return
	}

	cmd := exec.Command(binPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	
	err := cmd.Start()
	if err != nil {
		fmt.Printf("❌ Failed to start RictusD: %v\n", err)
		return
	}

	_ = os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	fmt.Printf("💀 RictusD Enforcer started (PID: %d)\n", cmd.Process.Pid)
}

func stopDaemon() {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("⚠️ RictusD PID file not found. Process may not be running.")
		return
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)
	
	process, err := os.FindProcess(pid)
	if err == nil {
		_ = process.Signal(syscall.SIGTERM)
		_ = os.Remove(pidFile)
		fmt.Printf("🛑 RictusD (PID %d) termination signal sent.\n", pid)
	}
}

func checkPID() {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("⚪ RictusD is idle.")
		return
	}
	fmt.Printf("🟢 RictusD is active (PID: %s)\n", string(data))
}

func printHelp() {
	fmt.Println("RictusD Enforcer Control Daemon")
	fmt.Println("Usage:")
	fmt.Println("  rictusctl start   - Launch enforcer background process")
	fmt.Println("  rictusctl stop    - Terminate active enforcer")
	fmt.Println("  rictusctl restart - Cycle the enforcer process")
	fmt.Println("  rictusctl status  - Check process state")
	fmt.Println("  rictusctl inspect - View verified threat ledger") // Added to help
}
