package maintenance

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"testing"
	"time"
)

func activeWindow(id string, scopeIds, scopeTypes []string) *evt.MaintenanceWindow {
	now := time.Now().Unix()
	return &evt.MaintenanceWindow{
		WindowId:   id,
		Name:       "Window " + id,
		Status:     evt.MaintenanceStatus_MAINTENANCE_STATUS_ACTIVE,
		StartTime:  now - 3600, // 1 hour ago
		EndTime:    now + 3600, // 1 hour from now
		ScopeIds:   scopeIds,
		ScopeTypes: scopeTypes,
	}
}

func TestIsSuppressed_NoWindows(t *testing.T) {
	eval := New()
	if eval.IsSuppressed("node-1", "Router") {
		t.Error("expected not suppressed with no windows loaded")
	}
}

func TestIsSuppressed_GlobalWindow(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", nil, nil), // no scope = applies to everything
	})

	if !eval.IsSuppressed("node-1", "Router") {
		t.Error("expected suppressed by global window")
	}
	if !eval.IsSuppressed("node-2", "Switch") {
		t.Error("expected suppressed by global window")
	}
}

func TestIsSuppressed_ByEntityID(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", []string{"node-1", "node-2"}, nil),
	})

	if !eval.IsSuppressed("node-1", "Router") {
		t.Error("expected node-1 suppressed")
	}
	if !eval.IsSuppressed("node-2", "Switch") {
		t.Error("expected node-2 suppressed")
	}
	if eval.IsSuppressed("node-3", "Router") {
		t.Error("expected node-3 not suppressed")
	}
}

func TestIsSuppressed_ByEntityType(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", nil, []string{"Router"}),
	})

	if !eval.IsSuppressed("node-1", "Router") {
		t.Error("expected Router suppressed")
	}
	if eval.IsSuppressed("node-1", "Switch") {
		t.Error("expected Switch not suppressed")
	}
}

func TestIsSuppressed_ExpiredWindow(t *testing.T) {
	now := time.Now().Unix()
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		{
			WindowId:  "w1",
			Status:    evt.MaintenanceStatus_MAINTENANCE_STATUS_ACTIVE,
			StartTime: now - 7200, // 2 hours ago
			EndTime:   now - 3600, // 1 hour ago (expired)
		},
	})

	if eval.IsSuppressed("node-1", "Router") {
		t.Error("expected not suppressed by expired window")
	}
}

func TestIsSuppressed_FutureWindow(t *testing.T) {
	now := time.Now().Unix()
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		{
			WindowId:  "w1",
			Status:    evt.MaintenanceStatus_MAINTENANCE_STATUS_SCHEDULED,
			StartTime: now + 3600, // 1 hour from now
			EndTime:   now + 7200, // 2 hours from now
		},
	})

	if eval.IsSuppressed("node-1", "Router") {
		t.Error("expected not suppressed by future window")
	}
}

func TestIsSuppressed_CompletedWindowIgnored(t *testing.T) {
	now := time.Now().Unix()
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		{
			WindowId:  "w1",
			Status:    evt.MaintenanceStatus_MAINTENANCE_STATUS_COMPLETED,
			StartTime: now - 3600,
			EndTime:   now + 3600,
		},
	})

	if eval.IsSuppressed("node-1", "Router") {
		t.Error("expected not suppressed by completed window")
	}
}

func TestIsSuppressed_CancelledWindowIgnored(t *testing.T) {
	now := time.Now().Unix()
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		{
			WindowId:  "w1",
			Status:    evt.MaintenanceStatus_MAINTENANCE_STATUS_CANCELLED,
			StartTime: now - 3600,
			EndTime:   now + 3600,
		},
	})

	if eval.IsSuppressed("node-1", "Router") {
		t.Error("expected not suppressed by cancelled window")
	}
}

func TestIsSuppressed_ScheduledWindowInRange(t *testing.T) {
	now := time.Now().Unix()
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		{
			WindowId:  "w1",
			Status:    evt.MaintenanceStatus_MAINTENANCE_STATUS_SCHEDULED,
			StartTime: now - 3600,
			EndTime:   now + 3600,
		},
	})

	if !eval.IsSuppressed("node-1", "Router") {
		t.Error("expected suppressed by scheduled window that is in time range")
	}
}

func TestGetActiveWindow_ReturnsWindow(t *testing.T) {
	eval := New()
	w := activeWindow("w1", []string{"node-1"}, nil)
	eval.LoadWindows([]*evt.MaintenanceWindow{w})

	got := eval.GetActiveWindow("node-1", "Router")
	if got == nil {
		t.Fatal("expected non-nil window")
	}
	if got.WindowId != "w1" {
		t.Errorf("expected window w1, got %s", got.WindowId)
	}
}

func TestGetActiveWindow_ReturnsNil(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", []string{"node-1"}, nil),
	})

	got := eval.GetActiveWindow("node-2", "Router")
	if got != nil {
		t.Error("expected nil for non-matching entity")
	}
}

func TestLoadWindows_ReplacesExisting(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", nil, nil),
	})
	if !eval.IsSuppressed("node-1", "") {
		t.Error("expected suppressed after first load")
	}

	eval.LoadWindows(nil)
	if eval.IsSuppressed("node-1", "") {
		t.Error("expected not suppressed after loading nil")
	}
}

func TestIsSuppressed_MultipleWindows(t *testing.T) {
	eval := New()
	eval.LoadWindows([]*evt.MaintenanceWindow{
		activeWindow("w1", []string{"node-1"}, nil),
		activeWindow("w2", nil, []string{"Switch"}),
	})

	if !eval.IsSuppressed("node-1", "Router") {
		t.Error("expected node-1 suppressed by w1")
	}
	if !eval.IsSuppressed("node-99", "Switch") {
		t.Error("expected Switch type suppressed by w2")
	}
	if eval.IsSuppressed("node-99", "Router") {
		t.Error("expected node-99 Router not suppressed")
	}
}
