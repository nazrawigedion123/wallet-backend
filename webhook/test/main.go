package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func readPayloadFromFile(filename string) ([]byte, error) {
	// Read the entire file into a byte slice using os.ReadFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return data, nil
}

func extractEventID(body []byte) (string, error) {
	var payload struct {
		EventID string `json:"event_id"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}
	if payload.EventID == "" {
		return "", fmt.Errorf("event_id is required")
	}
	return payload.EventID, nil
}

func main() {
	key := []byte("your_very_secret_key")

	// Read payload from file
	message, err := readPayloadFromFile("payload.json")
	if err != nil {
		fmt.Printf("Error reading payload: %v\n", err)
		os.Exit(1)
	}

	// Extract just the event_id from the payload
	eventID, err := extractEventID(message)
	if err != nil {
		fmt.Printf("Error extracting event_id: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using event_id for HMAC: %s\n", eventID)

	// Create a new HMAC by defining the hash type and key
	h := hmac.New(sha256.New, key)

	// Write only the event_id to it
	h.Write([]byte(eventID))

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(h.Sum(nil))
	fmt.Printf("HMAC signature: %s\n", signature)

	// Verify the signature
	receivedSig, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		os.Exit(1)
	}

	isValid := ValidMAC([]byte(eventID), receivedSig, key)
	fmt.Printf("Is valid HMAC: %t\n", isValid)

	// Tamper with the event_id and check again
	tamperedEventID := eventID + " "
	isValid = ValidMAC([]byte(tamperedEventID), receivedSig, key)
	fmt.Printf("Is valid HMAC after tampering: %t\n", isValid)
}
