# l8events

Shared event, alarm, and maintenance library for Layer 8 projects. Provides generic protobuf types
and Go backend packages that any consumer project (l8alarms, l8erp, etc.) can import. The
matching UI components (enums, tables, detail views) ship as part of the `l8ui` library at
`l8ui/events/` — see "l8ui Components" below.

l8events depends only on `google.golang.org/protobuf`. It does **not** depend on l8orm, l8services, l8bus, l8web, l8notify, or any consumer project.

## Directory Structure

```
l8events/
├── proto/
│   ├── l8events.proto                 # Core types (EventRecord, AlarmRecord, MaintenanceWindow, enums)
│   ├── l8events_categories.proto      # 16 category event messages + sub-category enums
│   └── make-bindings.sh               # Generates Go bindings
├── go/
│   ├── go.mod
│   ├── test.sh                        # Runs all tests with coverage
│   ├── types/l8events/
│   │   ├── l8events.pb.go             # Generated from l8events.proto
│   │   └── l8events_categories.pb.go  # Generated from l8events_categories.proto
│   ├── state/
│   │   └── state.go                   # Alarm state machine (transition validation)
│   ├── archive/
│   │   └── archive.go                 # Generic archive engine
│   ├── maintenance/
│   │   └── maintenance.go             # Maintenance window evaluator
│   ├── convert/
│   │   ├── convert.go                 # Converter engine (Parser interface, dispatch)
│   │   ├── helpers.go                 # Type conversion utilities
│   │   ├── parsers_ops.go             # 9 parsers: Audit, System, Monitoring, Security, Integration, Performance, Syslog, Trap, Automation
│   │   ├── parsers_infra.go           # 7 parsers: Network, Kubernetes, Compute, Storage, Power, GPU, Topology
│   │   └── builtins.go                # Built-in parser registration
│   └── tests/                         # All Go tests (black-box, package-external)
│       ├── state/state_test.go
│       ├── archive/archive_test.go
│       ├── maintenance/maintenance_test.go
│       └── convert/convert_core_test.go, convert_ops_test.go, convert_infra_test.go
└── (UI components live in the l8ui repo at l8ui/events/ — see "l8ui Components" below)
```

---

## Protobuf Types

Defined in `proto/l8events.proto` (core types) and `proto/l8events_categories.proto` (16 category event messages). Generate with `cd proto && ./make-bindings.sh`.

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

These components ship as part of the `l8ui` library at `l8ui/events/`, used identically on both
desktop and mobile — no per-platform variant, no separate mobile build. Add `l8ui` to your project
via `setup-l8ui-submodule.sh` (see `l8ui-copy-to-new-project.md`) and they're available
automatically; there is no separate copy step for l8events specifically.

**Prerequisites:** these components depend on l8ui shared library globals (all loaded on both
desktop and mobile pages): `Layer8DRenderers`, `Layer8EnumFactory`, `Layer8ColumnFactory`,
`Layer8FormFactory`.

Full API reference (script loading order, every exported global —
`L8EventsEnums`, `L8EventsCategoryEnums`, `L8EventsAlarmTable`, `L8EventsAlarmDetail`,
`L8EventsEventViewer`, `L8EventsArchiveViewer`, `L8EventsMaintenance`, `L8EventsStateActions`) now
lives in `l8ui`'s own docs: `l8ui/rules/l8events-ui.md`.

---

## Testing

All four Go packages have unit tests. Run the full suite with coverage:

```bash
cd go && ./test.sh
```

This script rebuilds dependencies from scratch, runs all tests with coverage across `state`, `archive`, `maintenance`, and `convert`, and opens the coverage report in a browser.

To run tests directly (after vendoring):
```bash
cd go && go test ./...
```

All tests live under `go/tests/`, one subdirectory per package, as black-box tests (`package xxx_test`) that exercise only the exported API.

| Package Under Test | Test Location | What It Covers |
|---------------------|----------------|----------------|
| `state` | `go/tests/state/state_test.go` | State transition validation, side effects (AcknowledgedBy, ClearedAt, etc.) |
| `archive` | `go/tests/archive/archive_test.go` | Cascade archival flow, Store interface mock |
| `maintenance` | `go/tests/maintenance/maintenance_test.go` | Window scope matching, time range evaluation |
| `convert` | `go/tests/convert/convert_core_test.go`, `convert_ops_test.go`, `convert_infra_test.go` | Parser dispatch, attribute mapping, error handling |

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

1. Add `l8ui` as a submodule under your project's web directory (`setup-l8ui-submodule.sh`, see
   `l8ui-copy-to-new-project.md`) — `l8ui/events/` comes with it, no separate copy step
2. Add the script includes to both `app.html` and `m/app.html` (same include list for both —
   see `l8ui/rules/l8events-ui.md` for the loading order)
3. Use the shared columns/forms/renderers in your module's definition files
