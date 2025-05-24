package models

type SimulationOptions struct {
	Count       int  `json:"count"`             // e.g., up to 1 million
	TierDist    bool `json:"tier_distribution"` // assign tiers randomly
	OutputToCSV bool `json:"output_to_csv"`
	TransactionTypes []string `json:"transaction_types"`
}
