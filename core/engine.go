package core

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"rictusd/modules/blockchain"
	"rictusd/modules/digitd"
)

type Engine struct {
	mu              sync.Mutex
	DigitDAddress   string
	Running         bool
	LastInteraction time.Time
	Ledger          *blockchain.Ledger
	Voice           *digitd.Client
}

func NewEngine(digitDIP string) *Engine {
	e := &Engine{
		DigitDAddress:   digitDIP,
		Running:         true,
		LastInteraction: time.Now(),
		Ledger:          blockchain.NewDumbChain(),
		Voice:           digitd.NewClient(digitDIP),
	}

	// Immediate Integrity Verification
	valid, err := e.Ledger.VerifyChain()
	if !valid {
		e.Log(fmt.Sprintf("🚨 LEDGER CORRUPTION: %v", err))
		e.ReportToDigitD("System alert: Local ledger integrity failed.")
	} else {
		e.Log("Ledger integrity verified successfully.")
	}

	go e.HandleSignals()
	e.Log("RictusD Engine online. [Hailo-8 Monitoring Active]")
	
	return e
}

// ProcessVisualThreat pipes NPU results to DigitD
func (e *Engine) ProcessVisualThreat(label string, confidence float64) {
	threat := map[string]interface{}{
		"source":     "Hailo-8_NPU",
		"event":      "Visual_Intrusion",
		"target":     label,
		"confidence": confidence,
	}

	// 1. Commit to the immutable ledger
	txHash := e.Ledger.CommitThreat(threat)
	e.Log(fmt.Sprintf("Visual threat confirmed and hashed: %s", txHash[:8]))
	
	// 2. High-priority vocalization via DigitD
	vocalAlert := fmt.Sprintf("Security alert. %s detected by NPU with %d percent certainty. Incident recorded in transaction %s.", 
		label, int(confidence*100), txHash[:8])
	
	e.ReportToDigitD(vocalAlert)
}

func (e *Engine) ReportToDigitD(alert string) {
	_ = e.Voice.Speak("Enforcer_Alert", alert)
}

func (e *Engine) HandleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	e.Log(fmt.Sprintf("Signal %v received. Standing down.", sig))
	os.Exit(0)
}
