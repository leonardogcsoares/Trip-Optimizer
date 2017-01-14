package skyscanner

import "time"

// CheapestPathRequest TODO
type CheapestPathRequest struct {
	Start  CPStart   `json:"start"`
	Places []CPPlace `json:"places"`
}

// CPStart TODO
type CPStart struct {
	Name string `json:"name"`
	Date string `json:"date"` // Format: YYYY-MM
}

// CPPlace TODO
type CPPlace struct {
	Name      string `json:"name"`
	Stay      int    `json:"stay_duration"`
	startDate time.Time
}
