package state

import (
	"fmt"
	evt "github.com/saichler/l8events/go/types/l8events"
	"time"
)

// validTransitions defines allowed state transitions.
// CLEARED is terminal — no transitions out.
var validTransitions = map[evt.AlarmState][]evt.AlarmState{
	evt.AlarmState_ALARM_STATE_ACTIVE: {
		evt.AlarmState_ALARM_STATE_ACKNOWLEDGED,
		evt.AlarmState_ALARM_STATE_CLEARED,
		evt.AlarmState_ALARM_STATE_SUPPRESSED,
	},
	evt.AlarmState_ALARM_STATE_ACKNOWLEDGED: {
		evt.AlarmState_ALARM_STATE_ACTIVE,
		evt.AlarmState_ALARM_STATE_CLEARED,
		evt.AlarmState_ALARM_STATE_SUPPRESSED,
	},
	evt.AlarmState_ALARM_STATE_SUPPRESSED: {
		evt.AlarmState_ALARM_STATE_ACTIVE,
		evt.AlarmState_ALARM_STATE_ACKNOWLEDGED,
		evt.AlarmState_ALARM_STATE_CLEARED,
	},
}

// ValidTransition returns true if transitioning from → to is allowed.
func ValidTransition(from, to evt.AlarmState) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition updates the alarm's state and appends a state history entry.
// Returns error if the transition is invalid.
func Transition(alarm *evt.AlarmRecord, newState evt.AlarmState, changedBy, reason string) error {
	if alarm == nil {
		return fmt.Errorf("alarm is nil")
	}

	if !ValidTransition(alarm.State, newState) {
		return fmt.Errorf("invalid transition from %s to %s",
			alarm.State.String(), newState.String())
	}

	now := time.Now().Unix()
	oldState := alarm.State

	alarm.StateHistory = append(alarm.StateHistory, &evt.AlarmStateChange{
		FromState: oldState,
		ToState:   newState,
		ChangedBy: changedBy,
		Reason:    reason,
		ChangedAt: now,
	})

	alarm.State = newState

	switch newState {
	case evt.AlarmState_ALARM_STATE_ACKNOWLEDGED:
		alarm.AcknowledgedBy = changedBy
		alarm.AcknowledgedAt = now
	case evt.AlarmState_ALARM_STATE_CLEARED:
		alarm.ClearedBy = changedBy
		alarm.ClearedAt = now
	case evt.AlarmState_ALARM_STATE_SUPPRESSED:
		alarm.IsSuppressed = true
		alarm.SuppressedBy = changedBy
	case evt.AlarmState_ALARM_STATE_ACTIVE:
		alarm.IsSuppressed = false
		alarm.SuppressedBy = ""
	}

	return nil
}

// Acknowledge is a convenience for transitioning to ACKNOWLEDGED.
func Acknowledge(alarm *evt.AlarmRecord, acknowledgedBy string) error {
	return Transition(alarm, evt.AlarmState_ALARM_STATE_ACKNOWLEDGED, acknowledgedBy, "")
}

// Clear is a convenience for transitioning to CLEARED.
func Clear(alarm *evt.AlarmRecord, clearedBy string) error {
	return Transition(alarm, evt.AlarmState_ALARM_STATE_CLEARED, clearedBy, "")
}

// Suppress is a convenience for transitioning to SUPPRESSED.
func Suppress(alarm *evt.AlarmRecord, suppressedBy string) error {
	return Transition(alarm, evt.AlarmState_ALARM_STATE_SUPPRESSED, suppressedBy, "")
}
