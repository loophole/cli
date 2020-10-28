package models

import "time"

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	SiteID    string    `json:"siteId"`
	Message   string    `json:"message"`
}
