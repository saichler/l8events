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
│   ├── maintenance/
│   │   └── maintenance.go        # Maintenance window evaluator
│   └── convert/
│       ├── convert.go            # Converter engine (Parser interface, dispatch)
│       ├── helpers.go            # Type conversion utilities
│       ├── parsers_ops.go        # 9 parsers: Audit, System, Monitoring, Security, Integration, Performance, Syslog, Trap, Automation
│       ├── parsers_infra.go      # 7 parsers: Network, Kubernetes, Compute, Storage, Power, GPU, Topology
│       └── builtins.go           # Built-in parser registration
└── l8ui/events/
    ├── l8events-enums.js              # Core enums (Severity, AlarmState, EventState, EventCategory, etc.)
    ├── l8events-category-enums.js     # Sub-category enums per EventCategory (15 enums + renderers)
    ├── l8events-alarm-table.js        # Alarm table columns + form definition
    ├── l8events-alarm-detail.js       # Alarm detail renderer (state history, notes)
    ├── l8events-event-viewer.js       # Event table columns + form definition
    ├── l8events-archive-viewer.js     # Archive table columns (alarm + event)
    ├── l8events-maintenance.js        # Maintenance window columns + form definition
    ├── l8events-state-actions.js      # Acknowledge/Clear/Suppress action buttons
    └── l8events.css                   # Shared event/alarm styles (--layer8d-* tokens)
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
| `EventCategory` | 0=UNSPECIFIED, 1=AUDIT, 2=SYSTEM, 3=MONITORING, 4=SECURITY, 5=INTEGRATION, 6=CUSTOM, 7=NETWORK, 8=KUBERNETES, 9=PERFORMANCE, 10=SYSLOG, 11=TRAP, 12=COMPUTE, 13=STORAGE, 14=POWER, 15=GPU, 16=TOPOLOGY, 17=AUTOMATION | `EventCategory_EVENT_CATEGORY_` |
| `MaintenanceStatus` | 0=UNSPECIFIED, 1=SCHEDULED, 2=ACTIVE, 3=COMPLETED, 4=CANCELLED | `MaintenanceStatus_MAINTENANCE_STATUS_` |
| `RecurrenceType` | 0=UNSPECIFIED, 1=NONE, 2=DAILY, 3=WEEKLY, 4=MONTHLY | `RecurrenceType_RECURRENCE_TYPE_` |

### Sub-Category Enums (per EventCategory)

| Enum | Category | Values |
|------|----------|--------|
| `AuditEventType` | Audit (1) | 0=UNSPECIFIED, 1=CREATE, 2=UPDATE, 3=DELETE, 4=LOGIN, 5=LOGOUT, 6=CONFIG_CHANGE, 7=PERMISSION_CHANGE, 8=EXPORT, 9=IMPORT |
| `SystemEventType` | System (2) | 0=UNSPECIFIED, 1=SERVICE_START, 2=SERVICE_STOP, 3=HEALTH_CHECK, 4=CONFIG_RELOAD, 5=LICENSE, 6=ERROR, 7=UPGRADE, 8=BACKUP, 9=RESTORE |
| `MonitoringEventType` | Monitoring (3) | 0=UNSPECIFIED, 1=POLL_SUCCESS, 2=POLL_FAILURE, 3=TARGET_UNREACHABLE, 4=TARGET_RECOVERED, 5=DATA_STALE, 6=COLLECTION_START, 7=COLLECTION_COMPLETE, 8=PARSE_ERROR |
| `SecurityEventType` | Security (4) | 0=UNSPECIFIED, 1=AUTH_SUCCESS, 2=AUTH_FAILURE, 3=ACCESS_DENIED, 4=PRIVILEGE_ESCALATION, 5=CERT_EXPIRY, 6=CERT_RENEWED, 7=POLICY_VIOLATION, 8=BRUTE_FORCE, 9=TOKEN_REVOKED |
| `IntegrationEventType` | Integration (5) | 0=UNSPECIFIED, 1=API_CALL_SUCCESS, 2=API_CALL_FAILURE, 3=WEBHOOK_RECEIVED, 4=WEBHOOK_FAILED, 5=SYNC_START, 6=SYNC_COMPLETE, 7=SYNC_FAILED, 8=CONNECTOR_UP, 9=CONNECTOR_DOWN |
| `NetworkEventType` | Network (7) | 0=UNSPECIFIED, 1=DEVICE_STATUS, 2=INTERFACE, 3=BGP, 4=OSPF, 5=MPLS, 6=LDP, 7=SR, 8=TE, 9=VRF, 10=QOS, 11=HARDWARE |
| `KubernetesEventType` | Kubernetes (8) | 0=UNSPECIFIED, 1=POD, 2=NODE, 3=DEPLOYMENT, 4=STATEFULSET, 5=DAEMONSET, 6=SERVICE, 7=NAMESPACE, 8=NETWORK_POLICY |
| `PerformanceMetric` | Performance (9) | 0=UNSPECIFIED, 1=CPU, 2=MEMORY, 3=TEMPERATURE, 4=TRAFFIC, 5=DISK, 6=FAN_SPEED, 7=POWER_LOAD, 8=VOLTAGE, 9=LATENCY, 10=PACKET_LOSS |
| `ThresholdType` | Performance (9) | 0=UNSPECIFIED, 1=UPPER, 2=LOWER |
| `ComputeEventType` | Compute (12) | 0=UNSPECIFIED, 1=HYPERVISOR_STATUS, 2=VM_STATUS, 3=VM_MIGRATION, 4=VM_RESOURCE, 5=HOST_RESOURCE |
| `StorageEventType` | Storage (13) | 0=UNSPECIFIED, 1=ARRAY_STATUS, 2=VOLUME_STATUS, 3=CAPACITY, 4=REPLICATION, 5=DISK, 6=CONTROLLER |
| `PowerEventType` | Power (14) | 0=UNSPECIFIED, 1=PSU_STATUS, 2=PDU_STATUS, 3=UPS_STATUS, 4=BATTERY, 5=LOAD, 6=VOLTAGE, 7=TEMPERATURE |
| `GpuEventType` | GPU (15) | 0=UNSPECIFIED, 1=STATUS, 2=TEMPERATURE, 3=MEMORY, 4=UTILIZATION, 5=ERROR, 6=POWER |
| `TopologyEventType` | Topology (16) | 0=UNSPECIFIED, 1=LINK_DISCOVERED, 2=LINK_LOST, 3=NEIGHBOR_CHANGE, 4=TOPOLOGY_CHANGE |
| `AutomationEventType` | Automation (17) | 0=UNSPECIFIED, 1=RULE_TRIGGERED, 2=RULE_COMPLETED, 3=RULE_FAILED, 4=POLICY_VIOLATION, 5=REMEDIATION |

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

### Category Event Messages

All category event messages share common fields: `eventId`, `propertyId`, `sourceId`, `sourceType`. The `propertyId` is a string from `l8reflect/go/reflect/properties` that references the exact attribute in a source model. Each message adds domain-specific parsed fields.

#### `AuditEvent` (Category 1)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | AuditEventType (int) | CREATE, UPDATE, DELETE, LOGIN, etc. |
| `userId` | string | Who performed the action |
| `userName` | string | Display name |
| `userIp` | string | Client IP |
| `action` | string | Action description |
| `serviceName` | string | Service that processed the action |
| `serviceArea` | int32 | Service area number |
| `entityName` | string | Affected entity name |
| `previousValue` | string | Value before change |
| `newValue` | string | Value after change |
| `message` | string | Summary message |

#### `SystemEvent` (Category 2)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | SystemEventType (int) | SERVICE_START, HEALTH_CHECK, ERROR, etc. |
| `serviceName` | string | Affected service |
| `nodeId` | string | Node identifier |
| `nodeIp` | string | Node IP address |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `version` | string | Version (for upgrades) |
| `errorCode` | string | Error code |
| `errorDetail` | string | Error details |
| `message` | string | Summary message |

#### `MonitoringEvent` (Category 3)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | MonitoringEventType (int) | POLL_SUCCESS, TARGET_UNREACHABLE, etc. |
| `targetId` | string | Monitored target ID |
| `targetName` | string | Target display name |
| `targetType` | string | Target type |
| `protocol` | string | Monitoring protocol used |
| `pollDurationMs` | int64 | Poll duration in ms |
| `itemsCollected` | int32 | Number of items collected |
| `errorCode` | string | Error code |
| `errorDetail` | string | Error details |
| `lastSuccessAt` | int64 | Unix timestamp of last success |
| `staleDurationSec` | int64 | Seconds since last fresh data |
| `message` | string | Summary message |

#### `SecurityEvent` (Category 4)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | SecurityEventType (int) | AUTH_FAILURE, CERT_EXPIRY, BRUTE_FORCE, etc. |
| `userId` | string | User involved |
| `userName` | string | User display name |
| `userIp` | string | Client IP |
| `targetResource` | string | Targeted resource |
| `authMethod` | string | Authentication method |
| `failureReason` | string | Failure reason |
| `attemptCount` | int32 | Number of attempts |
| `certSubject` | string | Certificate subject |
| `certExpiry` | int64 | Certificate expiry timestamp |
| `policyName` | string | Violated policy name |
| `message` | string | Summary message |

#### `IntegrationEvent` (Category 5)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | IntegrationEventType (int) | API_CALL_SUCCESS, WEBHOOK_RECEIVED, etc. |
| `integrationName` | string | Integration identifier |
| `remoteSystem` | string | Remote system name |
| `remoteUrl` | string | Remote URL |
| `httpMethod` | string | HTTP method |
| `httpStatus` | int32 | HTTP status code |
| `requestDurationMs` | int64 | Request duration in ms |
| `itemsSynced` | int32 | Items synchronized |
| `errorCode` | string | Error code |
| `errorDetail` | string | Error details |
| `retryCount` | int32 | Retry attempts |
| `message` | string | Summary message |

#### `NetworkEvent` (Category 7)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `deviceName` | string | Device hostname |
| `deviceIp` | string | Device IP address |
| `deviceType` | int32 | Device type enum |
| `subCategory` | NetworkEventType (int) | DEVICE_STATUS, INTERFACE, BGP, OSPF, etc. |
| `componentId` | string | Component identifier (interface name, peer IP, etc.) |
| `componentName` | string | Component display name |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `message` | string | Summary message |

#### `KubernetesEvent` (Category 8)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `clusterId` | string | Cluster identifier |
| `namespace` | string | K8s namespace |
| `subCategory` | KubernetesEventType (int) | POD, NODE, DEPLOYMENT, etc. |
| `resourceName` | string | Resource name |
| `resourceKind` | string | Resource kind |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `reason` | string | K8s event reason |
| `message` | string | Summary message |
| `containerName` | string | Container name (for pod events) |
| `readyReplicas` | int32 | Ready replica count |
| `desiredReplicas` | int32 | Desired replica count |

#### `PerformanceEvent` (Category 9)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | PerformanceMetric (int) | CPU, MEMORY, TEMPERATURE, TRAFFIC, etc. |
| `metricName` | string | Metric name |
| `metricUnit` | string | Unit (%, bytes, Celsius, etc.) |
| `currentValue` | double | Current metric value |
| `thresholdValue` | double | Threshold that was crossed |
| `thresholdType` | ThresholdType (int) | UPPER or LOWER |
| `baselineValue` | double | Normal baseline value |
| `durationSeconds` | int64 | How long threshold exceeded |
| `componentId` | string | Component ID (interface, disk, etc.) |
| `componentName` | string | Component display name |
| `message` | string | Summary message |

#### `SyslogEvent` (Category 10)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `deviceName` | string | Device hostname |
| `deviceIp` | string | Device IP |
| `facility` | int32 | Syslog facility code |
| `facilityName` | string | Facility name (kern, user, etc.) |
| `syslogSeverity` | int32 | Syslog severity (0-7) |
| `syslogSeverityName` | string | Severity name (emerg, alert, etc.) |
| `mnemonic` | string | Message mnemonic |
| `processName` | string | Originating process |
| `rawMessage` | string | Original message text |
| `parsedMessage` | string | Parsed/normalized message |
| `timestamp` | int64 | Syslog timestamp |

#### `TrapEvent` (Category 11)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `deviceName` | string | Device hostname |
| `deviceIp` | string | Device IP |
| `trapOid` | string | Trap OID |
| `trapName` | string | Trap name (from MIB) |
| `genericTrap` | int32 | Generic trap number (SNMPv1) |
| `specificTrap` | int32 | Specific trap number (SNMPv1) |
| `enterpriseOid` | string | Enterprise OID |
| `snmpVersion` | string | SNMP version (v1/v2c/v3) |
| `community` | string | Community string |
| `varbinds` | map[string]string | Variable bindings |
| `uptime` | int64 | Device uptime |
| `message` | string | Summary message |

#### `ComputeEvent` (Category 12)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | ComputeEventType (int) | HYPERVISOR_STATUS, VM_STATUS, VM_MIGRATION, etc. |
| `hostName` | string | Host name |
| `hostIp` | string | Host IP |
| `vmName` | string | VM name |
| `vmId` | string | VM identifier |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `cpuCount` | int32 | CPU count |
| `memoryMb` | int64 | Memory in MB |
| `message` | string | Summary message |

#### `StorageEvent` (Category 13)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | StorageEventType (int) | ARRAY_STATUS, VOLUME_STATUS, CAPACITY, etc. |
| `arrayName` | string | Storage array name |
| `volumeName` | string | Volume name |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `capacityBytes` | int64 | Total capacity |
| `usedBytes` | int64 | Used capacity |
| `usagePercent` | double | Usage percentage |
| `message` | string | Summary message |

#### `PowerEvent` (Category 14)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | PowerEventType (int) | PSU_STATUS, PDU_STATUS, UPS_STATUS, BATTERY, etc. |
| `deviceName` | string | Device name |
| `componentName` | string | Component name (PSU #1, etc.) |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `voltage` | double | Voltage reading |
| `currentAmps` | double | Current in amps |
| `loadPercent` | double | Load percentage |
| `wattage` | double | Power in watts |
| `batteryPercent` | double | Battery level |
| `runtimeMinutes` | int32 | Remaining runtime |
| `message` | string | Summary message |

#### `GpuEvent` (Category 15)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | GpuEventType (int) | STATUS, TEMPERATURE, MEMORY, UTILIZATION, ERROR, POWER |
| `deviceName` | string | Device name |
| `hostName` | string | Host name |
| `gpuIndex` | int32 | GPU index on host |
| `gpuModel` | string | GPU model name |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `temperatureCelsius` | double | Temperature |
| `utilizationPercent` | double | GPU utilization |
| `memoryUsedBytes` | int64 | Memory used |
| `memoryTotalBytes` | int64 | Total memory |
| `powerDrawWatts` | double | Power draw |
| `eccErrors` | int64 | ECC error count |
| `message` | string | Summary message |

#### `TopologyEvent` (Category 16)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | TopologyEventType (int) | LINK_DISCOVERED, LINK_LOST, NEIGHBOR_CHANGE, TOPOLOGY_CHANGE |
| `localDeviceId` | string | Local device ID |
| `localDeviceName` | string | Local device name |
| `localInterface` | string | Local interface |
| `remoteDeviceId` | string | Remote device ID |
| `remoteDeviceName` | string | Remote device name |
| `remoteInterface` | string | Remote interface |
| `discoveryProtocol` | string | Discovery protocol (LLDP, CDP, etc.) |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `message` | string | Summary message |

#### `AutomationEvent` (Category 17)

| JSON Name | Type | Description |
|-----------|------|-------------|
| `eventId` | string | Primary key |
| `propertyId` | string | l8reflect property reference |
| `sourceId` | string | Source entity ID |
| `sourceType` | string | Source entity type |
| `subCategory` | AutomationEventType (int) | RULE_TRIGGERED, RULE_COMPLETED, RULE_FAILED, etc. |
| `ruleId` | string | Rule identifier |
| `ruleName` | string | Rule display name |
| `workflowId` | string | Workflow identifier |
| `triggerEventId` | string | Event that triggered the rule |
| `actionTaken` | string | Action performed |
| `previousState` | string | State before |
| `currentState` | string | State after |
| `success` | bool | Whether action succeeded |
| `errorMessage` | string | Error message if failed |
| `durationMs` | int64 | Execution duration in ms |
| `message` | string | Summary message |

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

### `convert` — Event Record Conversion Engine

Converts generic `EventRecord` instances (with data in `Attributes` map) into typed category-specific protobuf structs. Pre-loaded with all 16 built-in parsers. Supports custom parser registration.

```go
import "github.com/saichler/l8events/go/convert"
import evt "github.com/saichler/l8events/go/types/l8events"

// Create a converter (pre-loaded with all 16 category parsers)
conv := convert.New()

// Convert an EventRecord to its typed category struct
record := &evt.EventRecord{
    EventId:    "evt-001",
    Category:   evt.EventCategory_EVENT_CATEGORY_NETWORK,
    SourceId:   "switch-01",
    SourceType: "switch",
    Message:    "Interface down",
    Attributes: map[string]string{
        "propertyId":    "prop-1",
        "subCategory":   "2",
        "deviceName":    "core-sw-01",
        "deviceIp":      "10.1.1.1",
        "componentId":   "Gi0/1",
        "previousState": "up",
        "currentState":  "down",
    },
}

msg, err := conv.Convert(record)
// msg is a proto.Message — type-assert to the expected category struct
netEvent := msg.(*evt.NetworkEvent)
// netEvent.DeviceName == "core-sw-01"
// netEvent.SubCategory == NetworkEventType(2)
```

**Convert() behavior:**
| Input | Result |
|-------|--------|
| `nil` record | error |
| UNSPECIFIED category | error |
| CUSTOM category | `(nil, nil)` — no struct for custom events |
| Unregistered category | error |
| Valid category | typed `proto.Message` |

**Field mapping:**
- Common fields (`EventId`, `PropertyId`, `SourceId`, `SourceType`, `Message`) are copied from the record's top-level fields and `Attributes["propertyId"]`
- `SubCategory` is parsed from `Attributes["subCategory"]` as int32 (15 of 16 parsers — SyslogEvent has no SubCategory)
- Domain fields are parsed from `Attributes[camelCaseFieldName]` with type conversion (string, int32, int64, float64, bool)
- `TrapEvent.Varbinds` collects all attributes with prefix `varbinds.` into a map (e.g., `varbinds.1.3.6.1` → key `1.3.6.1`)

**Error strategy:** Lenient — missing attributes yield zero values (no error). Malformed numeric/bool strings return an error.

**Custom parser registration:**
```go
// Replace a built-in parser or register a new one
conv.Register(evt.EventCategory_EVENT_CATEGORY_AUDIT, &myCustomAuditParser{})
```

The `Parser` interface:
```go
type Parser interface {
    Parse(record *evt.EventRecord) (proto.Message, error)
}
```

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
<script src="l8ui/events/l8events-enums.js"></script>              <!-- Must be first (others depend on it) -->
<script src="l8ui/events/l8events-category-enums.js"></script>     <!-- Sub-category enums (depends on enums) -->
<script src="l8ui/events/l8events-state-actions.js"></script>      <!-- Before alarm-detail (detail uses it) -->
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
- `L8EventsEnums.EVENT_CATEGORY` — `{ enum: { 0: 'Unspecified', 1: 'Audit', ... 17: 'Automation' } }`
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

### `window.L8EventsCategoryEnums`

Sub-category enum maps and renderers for each `EventCategory`. Loaded from `l8events-category-enums.js`.

**Enum maps** (for `f.select()` in category-specific forms):
- `L8EventsCategoryEnums.AUDIT_EVENT_TYPE` — Create, Update, Delete, Login, Logout, Config Change, Permission Change, Export, Import
- `L8EventsCategoryEnums.SYSTEM_EVENT_TYPE` — Service Start/Stop, Health Check, Config Reload, License, Error, Upgrade, Backup, Restore
- `L8EventsCategoryEnums.MONITORING_EVENT_TYPE` — Poll Success/Failure, Target Unreachable/Recovered, Data Stale, Collection Start/Complete, Parse Error
- `L8EventsCategoryEnums.SECURITY_EVENT_TYPE` — Auth Success/Failure, Access Denied, Privilege Escalation, Cert Expiry/Renewed, Policy Violation, Brute Force, Token Revoked
- `L8EventsCategoryEnums.INTEGRATION_EVENT_TYPE` — API Call Success/Failure, Webhook Received/Failed, Sync Start/Complete/Failed, Connector Up/Down
- `L8EventsCategoryEnums.NETWORK_EVENT_TYPE` — Device Status, Interface, BGP, OSPF, MPLS, LDP, Segment Routing, Traffic Engineering, VRF, QoS, Hardware
- `L8EventsCategoryEnums.KUBERNETES_EVENT_TYPE` — Pod, Node, Deployment, StatefulSet, DaemonSet, Service, Namespace, Network Policy
- `L8EventsCategoryEnums.PERFORMANCE_METRIC` — CPU, Memory, Temperature, Traffic, Disk, Fan Speed, Power Load, Voltage, Latency, Packet Loss
- `L8EventsCategoryEnums.THRESHOLD_TYPE` — Upper, Lower
- `L8EventsCategoryEnums.COMPUTE_EVENT_TYPE` — Hypervisor Status, VM Status/Migration/Resource, Host Resource
- `L8EventsCategoryEnums.STORAGE_EVENT_TYPE` — Array Status, Volume Status, Capacity, Replication, Disk, Controller
- `L8EventsCategoryEnums.POWER_EVENT_TYPE` — PSU/PDU/UPS Status, Battery, Load, Voltage, Temperature
- `L8EventsCategoryEnums.GPU_EVENT_TYPE` — Status, Temperature, Memory, Utilization, Error, Power
- `L8EventsCategoryEnums.TOPOLOGY_EVENT_TYPE` — Link Discovered/Lost, Neighbor Change, Topology Change
- `L8EventsCategoryEnums.AUTOMATION_EVENT_TYPE` — Rule Triggered/Completed/Failed, Policy Violation, Remediation

**Renderers** (colored status badges for enums with status semantics, plain text for others):
- `L8EventsCategoryEnums.render.systemEventType` — colored (Start=green, Stop=red, Error=red, etc.)
- `L8EventsCategoryEnums.render.monitoringEventType` — colored (Poll Success=green, Failure=red, etc.)
- `L8EventsCategoryEnums.render.securityEventType` — colored (Auth Success=green, Auth Failure=red, etc.)
- `L8EventsCategoryEnums.render.integrationEventType` — colored (API Success=green, Failure=red, etc.)
- `L8EventsCategoryEnums.render.automationEventType` — colored (Completed=green, Failed=red, etc.)
- All others: plain enum text renderers (no color classes)

**Usage:**
```javascript
// In category-specific column definitions
...col.status('subCategory', 'Sub-type', null, L8EventsCategoryEnums.render.systemEventType),
...col.enum('subCategory', 'Sub-type', null, L8EventsCategoryEnums.render.networkEventType),

// In category-specific form definitions
...f.select('subCategory', 'Event Type', L8EventsCategoryEnums.NETWORK_EVENT_TYPE),
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
    "github.com/saichler/l8events/go/convert"
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
