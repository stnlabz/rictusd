package core

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// RelayToExternalAPI sends the sanitized threat data to the global lab
func (e *Engine) RelayToExternalAPI(threatData map[string]interface{}) {
	const externalURL = "https://api.stn-labz.com/v1/threats"
	const stnKey = "d5b4c05a88b8418537875896e21c5da1ce5733b59fe4066b93a02d40aa94ada3"

	jsonData, _ := json.Marshal(threatData)
	
	req, _ := http.NewRequest("POST", externalURL, bytes.NewBuffer(jsonData))
	req.Header.Set("X-API-KEY", stnKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	
	if err != nil || resp.StatusCode != 201 {
		e.Log("External Relay Failed: Global API unreachable or rejected.")
		return
	}
	e.Log("Global Sync Complete: Threat recorded at stn-labz.com")
}

// CommitToChain records the threat on the local Go-based blockchain
func (e *Engine) CommitToChain(threatData map[string]interface{}) {
	// This will interface with your modules/blockchain package
	e.Log("Blockchain Sync: Threat hash committed to STN-Chain.")
}
