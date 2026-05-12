package storage

import (
	"errors"
	"time"

	"github.com/roma-glushko/bp/internal/domain"
)

var ErrSessionNotFound = errors.New("session not found")

type Store interface {
	CreateSession(session *domain.MeasurementSession) error
	GetSession(id string) (*domain.MeasurementSession, error)
	UpdateSession(session *domain.MeasurementSession) error
	DeleteSession(id string) error
	ListSessions(from, to time.Time) ([]domain.MeasurementSession, error)

	GetSettings() (*domain.Settings, error)
	SaveSettings(settings *domain.Settings) error
}
