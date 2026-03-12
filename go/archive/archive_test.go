package archive

import (
	"fmt"
	evt "github.com/saichler/l8events/go/types/l8events"
	"testing"
)

// mockStore implements Store for testing.
type mockStore struct {
	alarms          map[string]*evt.AlarmRecord
	events          map[string]*evt.EventRecord
	eventsByAlarm   map[string][]string // alarmID -> []eventID
	archivedAlarms  map[string]*evt.AlarmRecord
	archivedEvents  map[string]*evt.EventRecord
	archiveInfos    map[string]*evt.ArchiveInfo // keyed by alarm/event ID
	failGetAlarm    bool
	failSaveAlarm   bool
	failDeleteAlarm bool
}

func newMockStore() *mockStore {
	return &mockStore{
		alarms:         make(map[string]*evt.AlarmRecord),
		events:         make(map[string]*evt.EventRecord),
		eventsByAlarm:  make(map[string][]string),
		archivedAlarms: make(map[string]*evt.AlarmRecord),
		archivedEvents: make(map[string]*evt.EventRecord),
		archiveInfos:   make(map[string]*evt.ArchiveInfo),
	}
}

func (s *mockStore) GetAlarm(alarmID string) (*evt.AlarmRecord, error) {
	if s.failGetAlarm {
		return nil, fmt.Errorf("store error")
	}
	return s.alarms[alarmID], nil
}

func (s *mockStore) SaveArchivedAlarm(alarm *evt.AlarmRecord, info *evt.ArchiveInfo) error {
	if s.failSaveAlarm {
		return fmt.Errorf("save error")
	}
	s.archivedAlarms[alarm.AlarmId] = alarm
	s.archiveInfos[alarm.AlarmId] = info
	return nil
}

func (s *mockStore) DeleteAlarm(alarmID string) error {
	if s.failDeleteAlarm {
		return fmt.Errorf("delete error")
	}
	delete(s.alarms, alarmID)
	return nil
}

func (s *mockStore) GetEventsByAlarm(alarmID string) ([]*evt.EventRecord, error) {
	var result []*evt.EventRecord
	for _, eid := range s.eventsByAlarm[alarmID] {
		if e, ok := s.events[eid]; ok {
			result = append(result, e)
		}
	}
	return result, nil
}

func (s *mockStore) SaveArchivedEvent(event *evt.EventRecord, info *evt.ArchiveInfo) error {
	s.archivedEvents[event.EventId] = event
	s.archiveInfos[event.EventId] = info
	return nil
}

func (s *mockStore) DeleteEvent(eventID string) error {
	delete(s.events, eventID)
	return nil
}

func TestArchiveAlarm_Success(t *testing.T) {
	store := newMockStore()
	store.alarms["a1"] = &evt.AlarmRecord{AlarmId: "a1", Name: "test alarm"}
	store.events["e1"] = &evt.EventRecord{EventId: "e1", GeneratedAlarmId: "a1"}
	store.events["e2"] = &evt.EventRecord{EventId: "e2", GeneratedAlarmId: "a1"}
	store.eventsByAlarm["a1"] = []string{"e1", "e2"}

	archiver := New(store)
	info, err := archiver.ArchiveAlarm("a1", "admin", "resolved")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected non-nil ArchiveInfo")
	}
	if info.ArchivedBy != "admin" {
		t.Errorf("expected ArchivedBy=admin, got %s", info.ArchivedBy)
	}
	if info.ArchiveReason != "resolved" {
		t.Errorf("expected ArchiveReason=resolved, got %s", info.ArchiveReason)
	}
	if info.ArchivedAt == 0 {
		t.Error("expected ArchivedAt to be set")
	}

	// Alarm should be archived and deleted from active
	if _, ok := store.archivedAlarms["a1"]; !ok {
		t.Error("alarm not found in archived store")
	}
	if _, ok := store.alarms["a1"]; ok {
		t.Error("alarm should have been deleted from active store")
	}

	// Events should be archived and deleted
	if _, ok := store.archivedEvents["e1"]; !ok {
		t.Error("event e1 not found in archived store")
	}
	if _, ok := store.archivedEvents["e2"]; !ok {
		t.Error("event e2 not found in archived store")
	}
	if _, ok := store.events["e1"]; ok {
		t.Error("event e1 should have been deleted from active store")
	}
}

func TestArchiveAlarm_NotFound(t *testing.T) {
	store := newMockStore()
	archiver := New(store)
	_, err := archiver.ArchiveAlarm("nonexistent", "admin", "")
	if err == nil {
		t.Error("expected error for nonexistent alarm")
	}
}

func TestArchiveAlarm_GetAlarmError(t *testing.T) {
	store := newMockStore()
	store.failGetAlarm = true
	archiver := New(store)
	_, err := archiver.ArchiveAlarm("a1", "admin", "")
	if err == nil {
		t.Error("expected error when GetAlarm fails")
	}
}

func TestArchiveAlarm_SaveError(t *testing.T) {
	store := newMockStore()
	store.alarms["a1"] = &evt.AlarmRecord{AlarmId: "a1"}
	store.failSaveAlarm = true
	archiver := New(store)
	_, err := archiver.ArchiveAlarm("a1", "admin", "")
	if err == nil {
		t.Error("expected error when SaveArchivedAlarm fails")
	}
}

func TestArchiveAlarm_DeleteError(t *testing.T) {
	store := newMockStore()
	store.alarms["a1"] = &evt.AlarmRecord{AlarmId: "a1"}
	store.failDeleteAlarm = true
	archiver := New(store)
	info, err := archiver.ArchiveAlarm("a1", "admin", "resolved")
	// Should return info but also an error about delete failure
	if info == nil {
		t.Error("expected ArchiveInfo even on delete failure")
	}
	if err == nil {
		t.Error("expected error when DeleteAlarm fails")
	}
}

func TestArchiveAlarm_NoEvents(t *testing.T) {
	store := newMockStore()
	store.alarms["a1"] = &evt.AlarmRecord{AlarmId: "a1"}
	archiver := New(store)
	info, err := archiver.ArchiveAlarm("a1", "admin", "cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected non-nil ArchiveInfo")
	}
	if len(store.archivedEvents) != 0 {
		t.Errorf("expected no archived events, got %d", len(store.archivedEvents))
	}
}

func TestArchiveEvent(t *testing.T) {
	store := newMockStore()
	archiver := New(store)
	info, err := archiver.ArchiveEvent("e1", "admin", "cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected non-nil ArchiveInfo")
	}
	if info.ArchivedBy != "admin" {
		t.Errorf("expected ArchivedBy=admin, got %s", info.ArchivedBy)
	}
	if info.ArchiveReason != "cleanup" {
		t.Errorf("expected ArchiveReason=cleanup, got %s", info.ArchiveReason)
	}
}
