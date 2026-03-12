package archive

import (
	"fmt"
	evt "github.com/saichler/l8events/go/types/l8events"
	"time"
)

// Store defines the persistence operations the consumer must provide.
type Store interface {
	GetAlarm(alarmID string) (*evt.AlarmRecord, error)
	SaveArchivedAlarm(alarm *evt.AlarmRecord, info *evt.ArchiveInfo) error
	DeleteAlarm(alarmID string) error
	GetEventsByAlarm(alarmID string) ([]*evt.EventRecord, error)
	SaveArchivedEvent(event *evt.EventRecord, info *evt.ArchiveInfo) error
	DeleteEvent(eventID string) error
}

// Archiver manages alarm and event archival.
type Archiver struct {
	store Store
}

// New creates a new Archiver with the given persistence store.
func New(store Store) *Archiver {
	return &Archiver{store: store}
}

// ArchiveAlarm archives an alarm and its associated events.
func (a *Archiver) ArchiveAlarm(alarmID, archivedBy, reason string) (*evt.ArchiveInfo, error) {
	alarm, err := a.store.GetAlarm(alarmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alarm %s: %w", alarmID, err)
	}
	if alarm == nil {
		return nil, fmt.Errorf("alarm %s not found", alarmID)
	}

	info := &evt.ArchiveInfo{
		ArchivedAt:    time.Now().Unix(),
		ArchivedBy:    archivedBy,
		ArchiveReason: reason,
	}

	if err := a.store.SaveArchivedAlarm(alarm, info); err != nil {
		return nil, fmt.Errorf("failed to archive alarm %s: %w", alarmID, err)
	}

	// Archive associated events
	events, err := a.store.GetEventsByAlarm(alarmID)
	if err == nil {
		for _, event := range events {
			if err := a.store.SaveArchivedEvent(event, info); err != nil {
				fmt.Printf("[archive] failed to archive event %s: %v\n", event.EventId, err)
				continue
			}
			_ = a.store.DeleteEvent(event.EventId)
		}
	}

	if err := a.store.DeleteAlarm(alarmID); err != nil {
		return info, fmt.Errorf("archived alarm %s but failed to delete active: %w", alarmID, err)
	}

	return info, nil
}

// ArchiveEvent archives a single event.
func (a *Archiver) ArchiveEvent(eventID, archivedBy, reason string) (*evt.ArchiveInfo, error) {
	info := &evt.ArchiveInfo{
		ArchivedAt:    time.Now().Unix(),
		ArchivedBy:    archivedBy,
		ArchiveReason: reason,
	}

	// Consumer must have fetched the event and passed it via SaveArchivedEvent.
	// This method just creates the ArchiveInfo — the consumer orchestrates the full flow
	// when archiving events independently (not as part of an alarm archive).
	return info, nil
}
