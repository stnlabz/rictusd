package blockchain

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Block struct {
	Index     int                    `json:"index"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	PrevHash  string                 `json:"prev_hash"`
	Hash      string                 `json:"hash"`
}

type Ledger struct {
	FilePath string
	LastHash string
}

func NewDumbChain() *Ledger {
	_ = os.MkdirAll("data/blockchain", 0755)
	return &Ledger{
		FilePath: "data/blockchain/local_chain.json",
		LastHash: "00000000000000000000000000000000", // Genesis placeholder
	}
}

// VerifyChain performs a cryptographic audit of the existing file
func (l *Ledger) VerifyChain() (bool, error) {
	file, err := os.Open(l.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // A missing chain is valid for a new setup
		}
		return false, err
	}
	defer file.Close()

	var lastHash string = "00000000000000000000000000000000"
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		var b Block
		if err := json.Unmarshal(scanner.Bytes(), &b); err != nil {
			return false, fmt.Errorf("corrupt line in ledger")
		}

		// 1. Validate Linkage
		if b.PrevHash != lastHash {
			return false, fmt.Errorf("chain broken at hash %s", b.Hash)
		}

		// 2. Validate Data Integrity
		record := fmt.Sprintf("%v%s%s", b.Data, b.Timestamp, b.PrevHash)
		calculatedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(record)))
		
		if b.Hash != calculatedHash {
			return false, fmt.Errorf("block %s data has been tampered with", b.Hash)
		}

		lastHash = b.Hash
	}

	l.LastHash = lastHash
	return true, nil
}

// CommitThreat hashes the data and appends it to our local "dumb" chain
func (l *Ledger) CommitThreat(threat map[string]interface{}) string {
	timestamp := time.Now().Format(time.RFC3339)
	
	record := fmt.Sprintf("%v%s%s", threat, timestamp, l.LastHash)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(record)))

	block := Block{
		Index:     0,
		Timestamp: timestamp,
		Data:      threat,
		PrevHash:  l.LastHash,
		Hash:      hash,
	}

	file, _ := os.OpenFile(l.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	entry, _ := json.Marshal(block)
	file.Write(append(entry, '\n'))
	file.Close()

	l.LastHash = hash
	return hash
}
