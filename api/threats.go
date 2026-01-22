package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rictusd/core"
)

// StartServer initializes the listener for remote Sentinel reports
func StartServer(engine *core.Engine, port string) {
	http.HandleFunc("/v1/threats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var threat map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&threat); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// 1. Process via Engine (Blockchain + Voice)
		engine.Log(fmt.Sprintf("Threat received from %s", r.RemoteAddr))
		
		// Commit to our "Dumb" Ledger
		txHash := engine.Ledger.CommitThreat(threat)
		
		// Relay to DigitD Voice
		alertMsg := fmt.Sprintf("External threat neutralized. Block Hash: %s", txHash[:8])
		engine.ReportToDigitD(alertMsg)

		// 2. Respond to Sentinel
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"captured", "hash":"` + txHash + `"}`))
	})

	engine.Log("API Server listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		engine.Log("API Server Failure: " + err.Error())
	}
}
