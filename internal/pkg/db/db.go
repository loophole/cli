package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/loophole/cli/internal/app/metrics/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dbFile := cache.GetLocalStorageFile("metrics.db", "")
	var err error
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Metric{})
}

func AddEvent(sessionId uuid.UUID, siteId string, eventMessage string) (*models.Event, error) {
	event := &models.Event{
		Timestamp: time.Now(),
		SessionID: sessionId.String(),
		SiteID:    siteId,
		Message:   eventMessage,
	}
	result := db.Create(event)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}

	return event, nil
}

func AddMetric(sessionId uuid.UUID, siteId string, metricName string, metricValue float64) (*models.Metric, error) {
	metric := &models.Metric{
		Timestamp: time.Now(),
		SessionID: sessionId.String(),
		SiteID:    siteId,
		Name:      metricName,
		Value:     metricValue,
	}
	result := db.Create(metric)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}

	return metric, nil
}

func GetLastStartupEvent() (*models.Event, error) {
	event := &models.Event{}
	result := db.Where("message = ?", "Site registered").Last(event)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}
	return event, nil
}

func GetAllEvents() ([]models.Event, error) {
	events := []models.Event{}
	result := db.Find(&events)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}
	return events, nil
}

func GetEvents(siteId string) ([]models.Event, error) {
	events := []models.Event{}
	result := db.Where(&models.Event{SiteID: siteId}).Find(&events)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}
	return events, nil
}

func GetAllMetrics() ([]models.Metric, error) {
	metrics := []models.Metric{}
	result := db.Find(&metrics)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}
	return metrics, nil
}

func GetMetrics(siteId string) ([]models.Metric, error) {
	metrics := []models.Metric{}
	result := db.Where(&models.Metric{SiteID: siteId}).Find(&metrics)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Problem fetching metrics")
		return nil, result.Error
	}
	return metrics, nil
}
