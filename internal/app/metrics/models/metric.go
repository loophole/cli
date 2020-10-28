package models

import "time"

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	SiteID    string    `json:"siteId"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
}
