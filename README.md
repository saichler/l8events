# l8events

Shared event, alarm, and maintenance library for Layer 8 projects. Provides generic protobuf types, Go backend packages, and reusable l8ui components that any consumer project (l8alarms, l8erp, etc.) can import.

l8events depends only on `google.golang.org/protobuf`. It does **not** depend on l8orm, l8services, l8bus, l8web, l8notify, or any consumer project.

## Directory Structure

```
l8events/
├── proto/
│   ├── l8events.proto            # Shared protobuf types
│   └── make-bindings.sh          # Generates Go bindings
├── go/
│   ├── go.mod
│   ├── types/l8events/
│   │   └── l8events.pb.go        # Generated protobuf Go types
│   ├── state/
│   │   └── state.go              # Alarm state machine (transition validation)
│   ├── archive/
│   │   └── archive.go            # Generic archive engine
│   └── maintenance/
│       └── maintenance.go        # Maintenance window evaluator
└── l8ui/events/
    ├── l8events-enums.js          # Shared enums + renderers
    ├── l8events-alarm-table.js    # Alarm table columns + form definition
    ├── l8events-alarm-detail.js   # Alarm detail renderer (state history, notes)
    ├── l8events-event-viewer.js   # Event table columns + form definition
    ├── l8events-archive-viewer.js # Archive table columns (alarm + event)
    ├── l8events-maintenance.js    # Maintenance window columns + form definition
    ├── l8events-state-actions.js  # Acknowledge/Clear/Suppress action buttons
    └── l8events.css               # Shared event/alarm styles (--layer8d-* tokens)
```

---

## Protobuf Types

Defined in `proto/l8events.proto`. Generate with `cd proto && ./make-bindings.sh`.

### Enums

| Enum | Values | Go Constants Prefix |
|------|--------|---------------------|
| `Severity` | 0=UNSPECIFIED, 1=INFO, 2=WARNING, 3=MINOR, 4=MAJOR, 5=CRITICAL | `Severity_SEVERITY_` |
| `AlarmState` | 0=UNSPECIFIED, 1=ACTIVE, 2=ACKNOWLEDGED, 3=CLEARED, 4=SUPPRESSED | `AlarmState_ALARM_STATE_` |
| `EventState` | 0=UNSPECIFIED, 1=NEW, 2=PROCESSED, 3=DISCARDED, 4=ARCHIVED | `EventState_EVENT_STATE_` |
| `EventCategory` | 0=UNSPECIFIED, 1=AUDIT, 2=SYSTEM, 3=MONITORING, 4=SECURITY, 5=INTEGRATION, 6=CUSTOM | `EventCategory_EVENT_CATEGORY_` |
| `MaintenanceStatus` | 0=UNSPECIFIED, 1=SCHEDULED, 2=ACTIVE, 3=COMPLETED, 4=CANCELLED | `MaintenanceStatus_MAINTENANCE_STATUS_` |
| `RecurrenceType` | 0=UNSPECIFIED, 1=NONE, 2=DAILY, 3=WEEKLY, 4=MONTHLY | `RecurrenceType_RECURRENCE_TYPE_` |

### Messages and JSON Field Names

#### `EventRecord` (immutable after creation)

| Go Field | JSON Name | Type | Description |
|----------|-----------|------|-------------|
| `EventId` | `eventId` | string | Primary key |
| `Category` | `category` | EventCategory (int) | Event category |
| `EventType` | `eventType` | string | Consumer-defined type (e.g., "TRAP", "USER_LOGIN") |
| `State` | `state` | EventState (int) | Processing state |
| `Severity` | `severity` | Severity (int) | Severity level |
| `SourceId` | `sourceId` | string | ID of the entity that generated the event |
| `SourceName` | `sourceName` | string | Human-readable source name |
| `SourceType` | `sourceType` | string | Type of source (e.g., "Node", "User") |
| `Message` | `message` | string | Event message |
| `Attributes` | `attributes` | map[string]string | Extensible key-value pairs |
| `OccurredAt` | `occurredAt` | int64 | Unix timestamp when event happened |
| `ReceivedAt` | `receivedAt` | int64 | Unix timestamp when system received it |
| `ProcessedAt` | `processedAt` | int64 | Unix timestamp when processing completed |
| `GeneratedAlarmId` | `generatedAlarmId` | string | If this event generated an alarm |

#### `AlarmRecord` (state machine lifecycle)

| Go Field | JSON Name | Type | Description |
|----------|-----------|------|-------------|
| `AlarmId` | `alarmId` | string | Primary key |
| `DefinitionId` | `definitionId` | string | What rule/definition triggered this alarm |
| `Name` | `name` | string | Alarm name |
| `Description` | `description` | string | Alarm description |
| `State` | `state` | AlarmState (int) | Current lifecycle state |
| `Severity` | `severity` | Severity (int) | Current severity |
| `OriginalSeverity` | `originalSeverity` | Severity (int) | Severity when first raised |
| `SourceId` | `sourceId` | string | Entity the alarm is about |
| `SourceName` | `sourceName` | string | Human-readable source name |
| `SourceType` | `sourceType` | string | Source type |
| `FirstOccurrence` | `firstOccurrence` | int64 | Unix timestamp of first occurrence |
| `LastOccurrence` | `lastOccurrence` | int64 | Unix timestamp of last occurrence |
| `OccurrenceCount` | `occurrenceCount` | int32 | Number of occurrences (dedup) |
| `DedupKey` | `dedupKey` | string | Deduplication key |
| `EventId` | `eventId` | string | Originating event ID |
| `AcknowledgedBy` | `acknowledgedBy` | string | Who acknowledged |
| `AcknowledgedAt` | `acknowledgedAt` | int64 | When acknowledged |
| `ClearedBy` | `clearedBy` | string | Who cleared |
| `ClearedAt` | `clearedAt` | int64 | When cleared |
| `IsSuppressed` | `isSuppressed` | bool | Whether suppressed |
| `SuppressedBy` | `suppressedBy` | string | Who/what suppressed |
| `Attributes` | `attributes` | map[string]string | Extensible key-value pairs |
| `Notes` | `notes` | []*AlarmNote | Attached notes |
| `StateHistory` | `stateHistory` | []*AlarmStateChange | State transition audit trail |

#### `AlarmNote` (child of AlarmRecord)

| Go Field | JSON Name | Type |
|----------|-----------|------|
| `NoteId` | `noteId` | string |
| `Author` | `author` | string |
| `Text` | `text` | string |
| `CreatedAt` | `createdAt` | int64 |

#### `AlarmStateChange` (child of AlarmRecord)

| Go Field | JSON Name | Type |
|----------|-----------|------|
| `FromState` | `fromState` | AlarmState (int) |
| `ToState` | `toState` | AlarmState (int) |
| `ChangedBy` | `changedBy` | string |
| `Reason` | `reason` | string |
| `ChangedAt` | `changedAt` | int64 |

#### `ArchiveInfo`

| Go Field | JSON Name | Type |
|----------|-----------|------|
| `ArchivedAt` | `archivedAt` | int64 |
| `ArchivedBy` | `archivedBy` | string |
| `ArchiveReason` | `archiveReason` | string |

#### `MaintenanceWindow`

| Go Field | JSON Name | Type | Description |
|----------|-----------|------|-------------|
| `WindowId` | `windowId` | string | Primary key |
| `Name` | `name` | string | Window name |
| `Description` | `description` | string | Description |
| `Status` | `status` | MaintenanceStatus (int) | Current status |
| `StartTime` | `startTime` | int64 | Unix timestamp |
| `EndTime` | `endTime` | int64 | Unix timestamp |
| `CreatedBy` | `createdBy` | string | Creator |
| `CreatedAt` | `createdAt` | int64 | Creation timestamp |
| `Recurrence` | `recurrence` | RecurrenceType (int) | Recurrence pattern |
| `RecurrenceInterval` | `recurrenceInterval` | int32 | Interval value |
| `ScopeIds` | `scopeIds` | []string | Entity IDs this window applies to |
| `ScopeTypes` | `scopeTypes` | []string | Entity types this window applies to |

---

## Go Packages

### `state` — Alarm State Machine

Validates and enforces alarm state transitions. CLEARED is terminal (no transitions out).

```go
import "github.com/saichler/l8events/go/state"
import evt "github.com/saichler/l8events/go/types/l8events"

// Check if a transition is valid
state.ValidTransition(evt.AlarmState_ALARM_STATE_ACTIVE, evt.AlarmState_ALARM_STATE_ACKNOWLEDGED) // true
state.ValidTransition(evt.AlarmState_ALARM_STATE_CLEARED, evt.AlarmState_ALARM_STATE_ACTIVE)       // false

// Apply a state transition (updates alarm.State, appends to alarm.StateHistory,
// and sets fields like AcknowledgedBy/AcknowledgedAt based on target state)
err := state.Transition(alarm, evt.AlarmState_ALARM_STATE_ACKNOWLEDGED, "admin", "investigating")

// Convenience functions (call Transition internally)
err := state.Acknowledge(alarm, "admin")
err := state.Clear(alarm, "system")
err := state.Suppress(alarm, "maintenance-window-1")
```

**Valid transitions:**
| From | To |
|------|-----|
| ACTIVE | ACKNOWLEDGED, CLEARED, SUPPRESSED |
| ACKNOWLEDGED | ACTIVE, CLEARED, SUPPRESSED |
| SUPPRESSED | ACTIVE, ACKNOWLEDGED, CLEARED |
| CLEARED | (terminal — no transitions) |

**Side effects of `Transition()`:**
- Appends an `AlarmStateChange` entry to `alarm.StateHistory`
- ACKNOWLEDGED: sets `alarm.AcknowledgedBy` and `alarm.AcknowledgedAt`
- CLEARED: sets `alarm.ClearedBy` and `alarm.ClearedAt`
- SUPPRESSED: sets `alarm.IsSuppressed = true` and `alarm.SuppressedBy`
- ACTIVE (reactivate): sets `alarm.IsSuppressed = false` and clears `alarm.SuppressedBy`

### `archive` — Archive Engine

Generic archive engine. Consumers provide a `Store` interface implementation for persistence.

```go
import "github.com/saichler/l8events/go/archive"

// Implement the Store interface
type myStore struct { /* ... */ }
func (s *myStore) GetAlarm(alarmID string) (*evt.AlarmRecord, error) { /* ... */ }
func (s *myStore) SaveArchivedAlarm(alarm *evt.AlarmRecord, info *evt.ArchiveInfo) error { /* ... */ }
func (s *myStore) DeleteAlarm(alarmID string) error { /* ... */ }
func (s *myStore) GetEventsByAlarm(alarmID string) ([]*evt.EventRecord, error) { /* ... */ }
func (s *myStore) SaveArchivedEvent(event *evt.EventRecord, info *evt.ArchiveInfo) error { /* ... */ }
func (s *myStore) DeleteEvent(eventID string) error { /* ... */ }

// Create archiver and archive
archiver := archive.New(&myStore{})
info, err := archiver.ArchiveAlarm("alarm-123", "admin", "resolved")
// ArchiveAlarm: fetches alarm, saves archived copy, archives all associated events, deletes originals

info, err := archiver.ArchiveEvent("event-456", "admin", "cleanup")
// ArchiveEvent: creates ArchiveInfo only — consumer orchestrates the full flow for standalone events
```

### `maintenance` — Maintenance Window Evaluator

Thread-safe evaluator that checks whether an entity is covered by an active maintenance window. Uses `sync.RWMutex` for concurrent access.

```go
import "github.com/saichler/l8events/go/maintenance"

eval := maintenance.New()

// Load/replace active windows (call on startup and when windows change)
eval.LoadWindows(activeWindows)

// Check if an entity is suppressed (by ID or type match)
if eval.IsSuppressed("node-456", "Router") {
    // Skip alarm creation or suppress notification
}

// Get the specific active window covering an entity
window := eval.GetActiveWindow("node-456", "Router")
```

**Scope matching logic:**
- If a window has no `ScopeIds` and no `ScopeTypes`, it applies to everything
- If `entityID` matches any entry in `window.ScopeIds`, the entity is covered
- If `entityType` matches any entry in `window.ScopeTypes`, the entity is covered
- Only windows with status ACTIVE or SCHEDULED are evaluated
- Only windows where `now` is between `StartTime` and `EndTime` are evaluated

---

## l8ui Components

Reusable UI components in `l8ui/events/`. Consumer projects copy these into their own `l8ui/events/` directory.

**Prerequisites:** These components depend on the l8ui shared library globals:
- `Layer8DRenderers` (provides `createStatusRenderer`, `renderEnum`)
- `Layer8EnumFactory` (provides `create()`)
- `Layer8ColumnFactory` (provides `col`, `status`, `enum`, `date`, `number`)
- `Layer8FormFactory` (provides `form`, `section`, `text`, `textarea`, `select`, `date`, `number`)

### Script Loading Order

```html
<!-- In app.html — after l8ui shared scripts, before module scripts -->
<link rel="stylesheet" href="l8ui/events/l8events.css">
<script src="l8ui/events/l8events-enums.js"></script>         <!-- Must be first (others depend on it) -->
<script src="l8ui/events/l8events-state-actions.js"></script>  <!-- Before alarm-detail (detail uses it) -->
<script src="l8ui/events/l8events-alarm-table.js"></script>
<script src="l8ui/events/l8events-alarm-detail.js"></script>
<script src="l8ui/events/l8events-event-viewer.js"></script>
<script src="l8ui/events/l8events-archive-viewer.js"></script>
<script src="l8ui/events/l8events-maintenance.js"></script>
```

### `window.L8EventsEnums`

Shared enum maps and renderers. All other l8events UI components depend on this.

**Enum maps** (for `f.select()` and column definitions):
- `L8EventsEnums.SEVERITY` — `{ enum: { 0: 'Unspecified', 1: 'Info', ... 5: 'Critical' } }`
- `L8EventsEnums.ALARM_STATE` — `{ enum: { 0: 'Unspecified', 1: 'Active', ... 4: 'Suppressed' } }`
- `L8EventsEnums.EVENT_STATE` — `{ enum: { 0: 'Unspecified', 1: 'New', ... 4: 'Archived' } }`
- `L8EventsEnums.EVENT_CATEGORY` — `{ enum: { 0: 'Unspecified', 1: 'Audit', ... 6: 'Custom' } }`
- `L8EventsEnums.MAINTENANCE_STATUS` — `{ enum: { 0: 'Unspecified', 1: 'Scheduled', ... 4: 'Cancelled' } }`
- `L8EventsEnums.RECURRENCE_TYPE` — `{ enum: { 0: 'Unspecified', 1: 'None', ... 4: 'Monthly' } }`

**Renderers** (for `col.status()` / `col.enum()` 4th argument):
- `L8EventsEnums.render.severity` — colored status badge (Critical=red, Major/Minor=warning, Warning=blue, Info=muted)
- `L8EventsEnums.render.alarmState` — colored status badge (Active=red, Acknowledged=blue, Cleared=green, Suppressed=muted)
- `L8EventsEnums.render.eventState` — colored status badge (New=blue, Processed=green, Discarded/Archived=muted)
- `L8EventsEnums.render.eventCategory` — plain enum text label
- `L8EventsEnums.render.maintenanceStatus` — colored status badge (Scheduled=blue, Active=warning, Completed=green, Cancelled=muted)
- `L8EventsEnums.render.recurrenceType` — plain enum text label

**Usage in consumer column definitions:**
```javascript
// Use shared renderers in your module's *-columns.js
...col.status('severity', 'Severity', null, L8EventsEnums.render.severity),
...col.status('state', 'State', null, L8EventsEnums.render.alarmState),

// Use shared enums in your module's *-forms.js
...f.select('severity', 'Severity', L8EventsEnums.SEVERITY),
...f.select('state', 'State', L8EventsEnums.ALARM_STATE),
```

### `window.L8EventsAlarmTable`

Provides reusable column and form definitions for alarm tables.

**`L8EventsAlarmTable.getColumns()`** — Returns column array:

| Key | Label | Type |
|-----|-------|------|
| `severity` | Severity | status badge |
| `name` | Name | text |
| `sourceName` | Source | text |
| `state` | State | status badge |
| `firstOccurrence` | First Occurrence | date |
| `lastOccurrence` | Last Occurrence | date |
| `occurrenceCount` | Count | number |
| `acknowledgedBy` | Acknowledged By | text |

**`L8EventsAlarmTable.getFormDefinition()`** — Returns form definition with sections:
- **Alarm Information**: alarmId, name, description, severity (select), state (select), definitionId
- **Source**: sourceId, sourceName, sourceType
- **Timing**: firstOccurrence, lastOccurrence, occurrenceCount
- **Acknowledgement**: acknowledgedBy, acknowledgedAt, clearedBy, clearedAt

**Usage in consumer:**
```javascript
// In your module's *-columns.js — use directly or merge with domain columns
MyModule.columns = {
    MyAlarmModel: L8EventsAlarmTable.getColumns()
};

// Or merge with domain-specific columns
MyModule.columns = {
    MyAlarmModel: [
        ...L8EventsAlarmTable.getColumns(),
        ...col.col('nodeId', 'Node'),  // domain-specific
    ]
};

// In your module's *-forms.js
MyModule.forms = {
    MyAlarmModel: L8EventsAlarmTable.getFormDefinition()
};
```

### `window.L8EventsAlarmDetail`

Renders a complete alarm detail view inside a container element. Includes fields display, state history timeline, notes, and optional state action buttons.

**`L8EventsAlarmDetail.render(container, alarm, options)`**

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | DOM Element | Target element to render into |
| `alarm` | object | AlarmRecord object (JSON from server) |
| `options.showStateHistory` | boolean | Show state history timeline (default: true) |
| `options.showNotes` | boolean | Show notes section (default: true) |
| `options.onStateChange` | function | Callback `(alarmId, newState, reason)` — if provided, renders action buttons |

**What it renders:**
1. Alarm fields grid: severity badge, state badge, name, source, occurrence count
2. State history timeline (reverse chronological): each entry shows `FromState -> ToState`, who changed it, when, and reason
3. Notes list (reverse chronological): author, date, text
4. State action buttons (if `onStateChange` callback provided): delegates to `L8EventsStateActions`

**Usage:**
```javascript
// In a popup's onShow callback
L8EventsAlarmDetail.render(popupBody, alarm, {
    showStateHistory: true,
    showNotes: true,
    onStateChange: (alarmId, newState, reason) => {
        // POST state change to your alarm service endpoint
        fetch(`/myprefix/${serviceArea}/MyAlarmSvc`, {
            method: 'PUT',
            headers: getHeaders(),
            body: JSON.stringify({ alarmId, state: newState })
        });
    }
});
```

### `window.L8EventsEventViewer`

Provides column and form definitions for read-only event log tables.

**`L8EventsEventViewer.getColumns()`** — Returns column array:

| Key | Label | Type |
|-----|-------|------|
| `occurredAt` | Timestamp | date |
| `category` | Category | enum (plain text) |
| `eventType` | Type | text |
| `severity` | Severity | status badge |
| `sourceName` | Source | text |
| `message` | Message | text |
| `state` | State | status badge |

**`L8EventsEventViewer.getFormDefinition()`** — Returns form definition with sections:
- **Event Information**: eventId, category (select), eventType, severity (select), state (select)
- **Source**: sourceId, sourceName, sourceType
- **Content**: message (textarea)
- **Timing**: occurredAt, receivedAt, processedAt

### `window.L8EventsArchiveViewer`

Provides column definitions for archived alarm and event tables. Extends the base columns with archive-specific fields.

**`L8EventsArchiveViewer.getArchivedAlarmColumns()`** — Returns column array:

| Key | Label | Type |
|-----|-------|------|
| `severity` | Severity | status badge |
| `name` | Name | text |
| `sourceName` | Source | text |
| `state` | State | status badge |
| `firstOccurrence` | First Occurrence | date |
| `clearedAt` | Cleared At | date |
| `archivedAt` | Archived At | date |
| `archivedBy` | Archived By | text |
| `archiveReason` | Reason | text |

**`L8EventsArchiveViewer.getArchivedEventColumns()`** — Returns column array:

| Key | Label | Type |
|-----|-------|------|
| `occurredAt` | Timestamp | date |
| `category` | Category | enum |
| `eventType` | Type | text |
| `severity` | Severity | status badge |
| `sourceName` | Source | text |
| `message` | Message | text |
| `archivedAt` | Archived At | date |
| `archivedBy` | Archived By | text |

### `window.L8EventsMaintenance`

Provides column and form definitions for maintenance window tables and create/edit forms.

**`L8EventsMaintenance.getColumns()`** — Returns column array:

| Key | Label | Type |
|-----|-------|------|
| `name` | Name | text |
| `status` | Status | status badge |
| `startTime` | Start Time | date |
| `endTime` | End Time | date |
| `recurrence` | Recurrence | enum (plain text) |
| `createdBy` | Created By | text |
| `createdAt` | Created At | date |

**`L8EventsMaintenance.getFormDefinition()`** — Returns form definition with sections:
- **Details**: name (required), description, status (select)
- **Schedule**: startTime (required), endTime (required), recurrence (select), recurrenceInterval
- **Scope**: scopeIds (comma-separated text), scopeTypes (comma-separated text)

### `window.L8EventsStateActions`

Renders Acknowledge/Clear/Suppress/Reactivate action buttons based on the alarm's current state.

**`L8EventsStateActions.render(container, alarm, onAction)`**

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | DOM Element | Target element to render buttons into |
| `alarm` | object | AlarmRecord object — reads `alarm.state` and `alarm.alarmId` |
| `onAction` | function | Callback `(alarmId, newState, reason)` called when a button is clicked |

**`L8EventsStateActions.getAvailableActions(currentState)`** — Returns array of action objects:

| Current State | Available Actions |
|---------------|-------------------|
| 1 (ACTIVE) | Acknowledge (->2), Clear (->3), Suppress (->4) |
| 2 (ACKNOWLEDGED) | Clear (->3), Suppress (->4) |
| 4 (SUPPRESSED) | Reactivate (->1), Acknowledge (->2), Clear (->3) |
| 3 (CLEARED) | (none — terminal state) |

Each action object: `{ state: <int>, label: <string>, className: <string> }`

**Button CSS classes:**
- `.l8events-action-acknowledge` — primary color
- `.l8events-action-clear` — success/green
- `.l8events-action-suppress` — light/muted
- `.l8events-action-reactivate` — warning/orange

### `l8events.css`

All styles use `--layer8d-*` theme tokens exclusively. No hardcoded colors, no `[data-theme="dark"]` blocks. Dark mode works automatically through the l8ui theme system.

**CSS class prefixes:**
- `.l8events-severity-*` — severity badge colors
- `.l8events-detail-*` — detail popup layout (section, grid, field)
- `.l8events-timeline-*` — state history timeline (entry, dot, content, transition, meta, reason)
- `.l8events-note-*` — notes section (header, author, date, text)
- `.l8events-state-actions` — action button container
- `.l8events-action-*` — individual action button colors
- `.l8events-maintenance-active` — maintenance window highlight

---

## Consumer Integration Pattern

### Go Backend

```go
import (
    "github.com/saichler/l8events/go/state"
    "github.com/saichler/l8events/go/archive"
    "github.com/saichler/l8events/go/maintenance"
    evt "github.com/saichler/l8events/go/types/l8events"
)
```

Add to `go.mod`:
```
require github.com/saichler/l8events/go v0.0.0-<latest>
```

### l8ui Components

1. Copy `l8ui/events/` into your project's web directory: `cp -r l8events/l8ui/events/ <project>/go/<app>/ui/web/l8ui/events/`
2. Add script includes to `app.html` (see loading order above)
3. Use the shared columns/forms/renderers in your module's definition files
