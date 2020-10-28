package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/loophole/cli/internal/pkg/db"
)

type CurrentSiteResponse struct {
	URL       string    `json:"url"`
	StartedAt time.Time `json:"startedAt"`
}

func CurrentSiteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	event, err := db.GetLastStartupEvent()
	current := CurrentSiteResponse{
		StartedAt: event.Timestamp,
		URL:       fmt.Sprintf("https://%s.loophole.site", event.SiteID),
	}
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		currentJSON, err := json.Marshal(current)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(currentJSON)
	}
}
