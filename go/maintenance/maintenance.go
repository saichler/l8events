package maintenance

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"sync"
	"time"
)

// Evaluator checks entities against active maintenance windows.
type Evaluator struct {
	windows []*evt.MaintenanceWindow
	mtx     sync.RWMutex
}

// New creates a new maintenance window Evaluator.
func New() *Evaluator {
	return &Evaluator{}
}

// LoadWindows replaces the current set of active windows.
func (e *Evaluator) LoadWindows(windows []*evt.MaintenanceWindow) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	e.windows = windows
}

// IsSuppressed returns true if the given entity (by ID or type) is covered
// by an active maintenance window at the current time.
func (e *Evaluator) IsSuppressed(entityID, entityType string) bool {
	return e.GetActiveWindow(entityID, entityType) != nil
}

// GetActiveWindow returns the active maintenance window covering the entity,
// or nil if none.
func (e *Evaluator) GetActiveWindow(entityID, entityType string) *evt.MaintenanceWindow {
	e.mtx.RLock()
	defer e.mtx.RUnlock()

	now := time.Now().Unix()

	for _, w := range e.windows {
		if w.Status != evt.MaintenanceStatus_MAINTENANCE_STATUS_ACTIVE &&
			w.Status != evt.MaintenanceStatus_MAINTENANCE_STATUS_SCHEDULED {
			continue
		}

		if now < w.StartTime || now > w.EndTime {
			continue
		}

		if !matchesScope(w, entityID, entityType) {
			continue
		}

		return w
	}

	return nil
}

func matchesScope(w *evt.MaintenanceWindow, entityID, entityType string) bool {
	// If no scope defined, window applies to everything
	if len(w.ScopeIds) == 0 && len(w.ScopeTypes) == 0 {
		return true
	}

	// Check entity ID match
	if entityID != "" && len(w.ScopeIds) > 0 {
		for _, id := range w.ScopeIds {
			if id == entityID {
				return true
			}
		}
	}

	// Check entity type match
	if entityType != "" && len(w.ScopeTypes) > 0 {
		for _, t := range w.ScopeTypes {
			if t == entityType {
				return true
			}
		}
	}

	return false
}
