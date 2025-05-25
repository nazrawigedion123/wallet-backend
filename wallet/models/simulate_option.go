package models

// type SimulationOptions struct {
// 	Count       int  `json:"count"`             // e.g., up to 1 million
// 	TierDist    bool `json:"tier_distribution"` // assign tiers randomly
// 	OutputToCSV bool `json:"output_to_csv"`
// 	TransactionTypes []string `json:"transaction_types"`

// }
type SimulationOptions struct {
	Count            int      `json:"count"`
	OutputToCSV      bool     `json:"output_to_csv"`
	IncludeFailed    bool     `json:"include_failed"`    // optional
	TransactionTypes []string `json:"transaction_types"` // ["payment", "transfer", etc.]
}
