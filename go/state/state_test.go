package state

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"testing"
)

func TestValidTransition_AllowedFromActive(t *testing.T) {
	allowed := []evt.AlarmState{
		evt.AlarmState_ALARM_STATE_ACKNOWLEDGED,
		evt.AlarmState_ALARM_STATE_CLEARED,
		evt.AlarmState_ALARM_STATE_SUPPRESSED,
	}
	for _, to := range allowed {
		if !ValidTransition(evt.AlarmState_ALARM_STATE_ACTIVE, to) {
			t.Errorf("expected ACTIVE -> %s to be valid", to)
		}
	}
}

func TestValidTransition_AllowedFromAcknowledged(t *testing.T) {
	allowed := []evt.AlarmState{
		evt.AlarmState_ALARM_STATE_ACTIVE,
		evt.AlarmState_ALARM_STATE_CLEARED,
		evt.AlarmState_ALARM_STATE_SUPPRESSED,
	}
	for _, to := range allowed {
		if !ValidTransition(evt.AlarmState_ALARM_STATE_ACKNOWLEDGED, to) {
			t.Errorf("expected ACKNOWLEDGED -> %s to be valid", to)
		}
	}
}

func TestValidTransition_AllowedFromSuppressed(t *testing.T) {
	allowed := []evt.AlarmState{
		evt.AlarmState_ALARM_STATE_ACTIVE,
		evt.AlarmState_ALARM_STATE_ACKNOWLEDGED,
		evt.AlarmState_ALARM_STATE_CLEARED,
	}
	for _, to := range allowed {
		if !ValidTransition(evt.AlarmState_ALARM_STATE_SUPPRESSED, to) {
			t.Errorf("expected SUPPRESSED -> %s to be valid", to)
		}
	}
}

func TestValidTransition_ClearedIsTerminal(t *testing.T) {
	all := []evt.AlarmState{
		evt.AlarmState_ALARM_STATE_ACTIVE,
		evt.AlarmState_ALARM_STATE_ACKNOWLEDGED,
		evt.AlarmState_ALARM_STATE_SUPPRESSED,
		evt.AlarmState_ALARM_STATE_CLEARED,
	}
	for _, to := range all {
		if ValidTransition(evt.AlarmState_ALARM_STATE_CLEARED, to) {
			t.Errorf("expected CLEARED -> %s to be invalid (terminal)", to)
		}
	}
}

func TestValidTransition_UnspecifiedIsInvalid(t *testing.T) {
	if ValidTransition(evt.AlarmState_ALARM_STATE_UNSPECIFIED, evt.AlarmState_ALARM_STATE_ACTIVE) {
		t.Error("expected UNSPECIFIED -> ACTIVE to be invalid")
	}
}

func TestTransition_NilAlarm(t *testing.T) {
	err := Transition(nil, evt.AlarmState_ALARM_STATE_ACKNOWLEDGED, "admin", "")
	if err == nil {
		t.Error("expected error for nil alarm")
	}
}

func TestTransition_InvalidTransition(t *testing.T) {
	alarm := &evt.AlarmRecord{State: evt.AlarmState_ALARM_STATE_CLEARED}
	err := Transition(alarm, evt.AlarmState_ALARM_STATE_ACTIVE, "admin", "")
	if err == nil {
		t.Error("expected error for CLEARED -> ACTIVE")
	}
}

func TestTransition_Acknowledge(t *testing.T) {
	alarm := &evt.AlarmRecord{State: evt.AlarmState_ALARM_STATE_ACTIVE}
	err := Acknowledge(alarm, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alarm.State != evt.AlarmState_ALARM_STATE_ACKNOWLEDGED {
		t.Errorf("expected ACKNOWLEDGED, got %s", alarm.State)
	}
	if alarm.AcknowledgedBy != "admin" {
		t.Errorf("expected AcknowledgedBy=admin, got %s", alarm.AcknowledgedBy)
	}
	if alarm.AcknowledgedAt == 0 {
		t.Error("expected AcknowledgedAt to be set")
	}
	if len(alarm.StateHistory) != 1 {
		t.Fatalf("expected 1 state history entry, got %d", len(alarm.StateHistory))
	}
	entry := alarm.StateHistory[0]
	if entry.FromState != evt.AlarmState_ALARM_STATE_ACTIVE {
		t.Errorf("expected FromState=ACTIVE, got %s", entry.FromState)
	}
	if entry.ToState != evt.AlarmState_ALARM_STATE_ACKNOWLEDGED {
		t.Errorf("expected ToState=ACKNOWLEDGED, got %s", entry.ToState)
	}
	if entry.ChangedBy != "admin" {
		t.Errorf("expected ChangedBy=admin, got %s", entry.ChangedBy)
	}
}

func TestTransition_Clear(t *testing.T) {
	alarm := &evt.AlarmRecord{State: evt.AlarmState_ALARM_STATE_ACTIVE}
	err := Clear(alarm, "system")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alarm.State != evt.AlarmState_ALARM_STATE_CLEARED {
		t.Errorf("expected CLEARED, got %s", alarm.State)
	}
	if alarm.ClearedBy != "system" {
		t.Errorf("expected ClearedBy=system, got %s", alarm.ClearedBy)
	}
	if alarm.ClearedAt == 0 {
		t.Error("expected ClearedAt to be set")
	}
}

func TestTransition_Suppress(t *testing.T) {
	alarm := &evt.AlarmRecord{State: evt.AlarmState_ALARM_STATE_ACTIVE}
	err := Suppress(alarm, "maint-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alarm.State != evt.AlarmState_ALARM_STATE_SUPPRESSED {
		t.Errorf("expected SUPPRESSED, got %s", alarm.State)
	}
	if !alarm.IsSuppressed {
		t.Error("expected IsSuppressed=true")
	}
	if alarm.SuppressedBy != "maint-1" {
		t.Errorf("expected SuppressedBy=maint-1, got %s", alarm.SuppressedBy)
	}
}

func TestTransition_Reactivate(t *testing.T) {
	alarm := &evt.AlarmRecord{
		State:        evt.AlarmState_ALARM_STATE_SUPPRESSED,
		IsSuppressed: true,
		SuppressedBy: "maint-1",
	}
	err := Transition(alarm, evt.AlarmState_ALARM_STATE_ACTIVE, "admin", "maintenance over")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alarm.State != evt.AlarmState_ALARM_STATE_ACTIVE {
		t.Errorf("expected ACTIVE, got %s", alarm.State)
	}
	if alarm.IsSuppressed {
		t.Error("expected IsSuppressed=false after reactivation")
	}
	if alarm.SuppressedBy != "" {
		t.Errorf("expected SuppressedBy cleared, got %s", alarm.SuppressedBy)
	}
	entry := alarm.StateHistory[0]
	if entry.Reason != "maintenance over" {
		t.Errorf("expected reason='maintenance over', got %s", entry.Reason)
	}
}

func TestTransition_MultipleTransitions(t *testing.T) {
	alarm := &evt.AlarmRecord{State: evt.AlarmState_ALARM_STATE_ACTIVE}
	if err := Acknowledge(alarm, "admin"); err != nil {
		t.Fatal(err)
	}
	if err := Clear(alarm, "admin"); err != nil {
		t.Fatal(err)
	}
	if len(alarm.StateHistory) != 2 {
		t.Errorf("expected 2 state history entries, got %d", len(alarm.StateHistory))
	}
}
