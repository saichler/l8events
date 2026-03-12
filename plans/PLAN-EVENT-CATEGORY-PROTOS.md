# Event Category Proto Definitions

## Intent

Each `EventCategory` will eventually become a Layer 8 Service. Each category gets its own proto message representing the **parsed, structured understanding** of the raw event — not the raw text.

The key concept is `property_id` (a string from `l8reflect/go/reflect/properties`), which links the event to the **exact attribute or nested attribute** in the source model that the event refers to. For example, an interface-down syslog event will have a `property_id` pointing to the specific `Interface.status` property within the `NetworkDevice` model hierarchy.

All category event messages embed a reference to the base `EventRecord` via `event_id`, inheriting the common fields (severity, timestamps, source, attributes) while adding category-specific parsed fields.

---

## Common Fields (all category events share these)

Every category event message includes:

| Field | Type | Description |
|-------|------|-------------|
| `event_id` | string | References the base `EventRecord` |
| `property_id` | string | l8reflect property ID pointing to the specific model attribute this event relates to |
| `source_id` | string | ID of the source entity (e.g., device ID, pod name, VM ID) |
| `source_type` | string | Type of source entity (e.g., "NetworkDevice", "K8sPod") |

---

## Category 1: AuditEvent

Events from user actions — CRUD operations, configuration changes, administrative actions. Captures who did what to which entity and when.

**Source**: Any Layer 8 service handling user-initiated operations

```
AuditEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute that was changed (e.g., NetworkDevice.equipmentinfo.location)
  source_id           string              → entity ID that was acted upon
  source_type         string              → entity type (e.g., "NetworkDevice", "K8sPod", "L8PTarget")
  sub_category        AuditEventType      → enum: CREATE, UPDATE, DELETE, LOGIN, LOGOUT, CONFIG_CHANGE, PERMISSION_CHANGE, EXPORT, IMPORT
  user_id             string              → user who performed the action
  user_name           string              → human-readable user name
  user_ip             string              → IP address of the user's session
  action              string              → HTTP method or action label (e.g., "POST", "PUT", "DELETE")
  service_name        string              → Layer 8 service that handled the request
  service_area        int32               → service area number
  entity_name         string              → human-readable name of the affected entity
  previous_value      string              → value before the change (serialized)
  new_value           string              → value after the change (serialized)
  message             string              → human-readable description
```

**AuditEventType enum**:
- AUDIT_EVENT_TYPE_UNSPECIFIED = 0
- AUDIT_EVENT_TYPE_CREATE = 1 (entity created via POST)
- AUDIT_EVENT_TYPE_UPDATE = 2 (entity modified via PUT)
- AUDIT_EVENT_TYPE_DELETE = 3 (entity deleted via DELETE)
- AUDIT_EVENT_TYPE_LOGIN = 4 (user login)
- AUDIT_EVENT_TYPE_LOGOUT = 5 (user logout)
- AUDIT_EVENT_TYPE_CONFIG_CHANGE = 6 (system configuration changed)
- AUDIT_EVENT_TYPE_PERMISSION_CHANGE = 7 (user/role permissions modified)
- AUDIT_EVENT_TYPE_EXPORT = 8 (data exported)
- AUDIT_EVENT_TYPE_IMPORT = 9 (data imported)

---

## Category 2: SystemEvent

Internal system events — service lifecycle, health checks, configuration reloads, license status, internal errors.

**Source**: Layer 8 platform services, vnet, UI server, log agent

```
SystemEvent
  event_id            string              → base EventRecord
  property_id         string              → specific system attribute (e.g., service health status)
  source_id           string              → service or component instance ID
  source_type         string              → "Service", "VNet", "UIServer", "LogAgent", "Database"
  sub_category        SystemEventType     → enum: SERVICE_START, SERVICE_STOP, HEALTH_CHECK, CONFIG_RELOAD, LICENSE, ERROR, UPGRADE, BACKUP, RESTORE
  service_name        string              → name of the service/component
  node_id             string              → node/host where the service runs
  node_ip             string              → IP of the node
  previous_state      string              → state before change
  current_state       string              → state after change
  version             string              → software version (for upgrade events)
  error_code          string              → error code (for error events)
  error_detail        string              → detailed error information
  message             string              → human-readable description
```

**SystemEventType enum**:
- SYSTEM_EVENT_TYPE_UNSPECIFIED = 0
- SYSTEM_EVENT_TYPE_SERVICE_START = 1 (service started)
- SYSTEM_EVENT_TYPE_SERVICE_STOP = 2 (service stopped)
- SYSTEM_EVENT_TYPE_HEALTH_CHECK = 3 (health check pass/fail)
- SYSTEM_EVENT_TYPE_CONFIG_RELOAD = 4 (configuration reloaded)
- SYSTEM_EVENT_TYPE_LICENSE = 5 (license expiry/renewal/violation)
- SYSTEM_EVENT_TYPE_ERROR = 6 (internal system error)
- SYSTEM_EVENT_TYPE_UPGRADE = 7 (software upgrade applied)
- SYSTEM_EVENT_TYPE_BACKUP = 8 (backup completed/failed)
- SYSTEM_EVENT_TYPE_RESTORE = 9 (restore completed/failed)

---

## Category 3: MonitoringEvent

Polling and collection lifecycle events — poll success/failure, data collection status, target reachability, data staleness detection.

**Source**: Collector service, parser service, cache services

```
MonitoringEvent
  event_id            string              → base EventRecord
  property_id         string              → specific monitoring attribute (e.g., target reachability)
  source_id           string              → target ID or collector ID
  source_type         string              → "L8PTarget", "Collector", "Parser", "Cache"
  sub_category        MonitoringEventType → enum: POLL_SUCCESS, POLL_FAILURE, TARGET_UNREACHABLE, TARGET_RECOVERED, DATA_STALE, COLLECTION_START, COLLECTION_COMPLETE, PARSE_ERROR
  target_id           string              → polling target ID
  target_name         string              → target name/IP
  target_type         string              → target type (e.g., "Network Device", "K8s Cluster")
  protocol            string              → polling protocol used (e.g., "SNMPv2", "SSH", "Kubectl")
  poll_duration_ms    int64               → how long the poll took
  items_collected     int32               → number of items/metrics collected
  error_code          string              → error code (for failure events)
  error_detail        string              → error details
  last_success_at     int64               → timestamp of last successful poll
  stale_duration_sec  int64               → how long data has been stale (for staleness events)
  message             string              → human-readable description
```

**MonitoringEventType enum**:
- MONITORING_EVENT_TYPE_UNSPECIFIED = 0
- MONITORING_EVENT_TYPE_POLL_SUCCESS = 1 (poll completed successfully)
- MONITORING_EVENT_TYPE_POLL_FAILURE = 2 (poll failed)
- MONITORING_EVENT_TYPE_TARGET_UNREACHABLE = 3 (target cannot be reached)
- MONITORING_EVENT_TYPE_TARGET_RECOVERED = 4 (previously unreachable target is back)
- MONITORING_EVENT_TYPE_DATA_STALE = 5 (data hasn't been updated beyond threshold)
- MONITORING_EVENT_TYPE_COLLECTION_START = 6 (collection cycle started)
- MONITORING_EVENT_TYPE_COLLECTION_COMPLETE = 7 (collection cycle completed)
- MONITORING_EVENT_TYPE_PARSE_ERROR = 8 (collected data failed to parse)

---

## Category 4: SecurityEvent

Security-related events — authentication, authorization, access control, certificate management, threat detection.

**Source**: Auth layer, API gateway, security policies

```
SecurityEvent
  event_id            string              → base EventRecord
  property_id         string              → specific security attribute
  source_id           string              → user ID, service ID, or policy ID
  source_type         string              → "User", "Service", "Policy", "Certificate"
  sub_category        SecurityEventType   → enum: AUTH_SUCCESS, AUTH_FAILURE, ACCESS_DENIED, PRIVILEGE_ESCALATION, CERT_EXPIRY, CERT_RENEWED, POLICY_VIOLATION, BRUTE_FORCE, TOKEN_REVOKED
  user_id             string              → user involved (if applicable)
  user_name           string              → human-readable user name
  user_ip             string              → IP address of the client
  target_resource     string              → resource being accessed (e.g., URL path, service name)
  auth_method         string              → authentication method (e.g., "password", "token", "certificate", "SSO")
  failure_reason      string              → why auth/access failed
  attempt_count       int32               → number of attempts (for brute force detection)
  cert_subject        string              → certificate subject (for cert events)
  cert_expiry         int64               → certificate expiration timestamp
  policy_name         string              → security policy name (for policy events)
  message             string              → human-readable description
```

**SecurityEventType enum**:
- SECURITY_EVENT_TYPE_UNSPECIFIED = 0
- SECURITY_EVENT_TYPE_AUTH_SUCCESS = 1 (successful authentication)
- SECURITY_EVENT_TYPE_AUTH_FAILURE = 2 (failed authentication)
- SECURITY_EVENT_TYPE_ACCESS_DENIED = 3 (authorized user denied access to resource)
- SECURITY_EVENT_TYPE_PRIVILEGE_ESCALATION = 4 (role/permission elevation)
- SECURITY_EVENT_TYPE_CERT_EXPIRY = 5 (certificate approaching/past expiration)
- SECURITY_EVENT_TYPE_CERT_RENEWED = 6 (certificate renewed)
- SECURITY_EVENT_TYPE_POLICY_VIOLATION = 7 (security policy violated)
- SECURITY_EVENT_TYPE_BRUTE_FORCE = 8 (multiple failed auth attempts detected)
- SECURITY_EVENT_TYPE_TOKEN_REVOKED = 9 (auth token invalidated)

---

## Category 5: IntegrationEvent

Events from external system integrations — API calls to/from third-party systems, webhook deliveries, data synchronization, connector health.

**Source**: Integration connectors, webhook receivers, external API clients

```
IntegrationEvent
  event_id            string              → base EventRecord
  property_id         string              → specific integration attribute
  source_id           string              → integration/connector ID
  source_type         string              → "Connector", "Webhook", "API", "SyncJob"
  sub_category        IntegrationEventType → enum: API_CALL_SUCCESS, API_CALL_FAILURE, WEBHOOK_RECEIVED, WEBHOOK_FAILED, SYNC_START, SYNC_COMPLETE, SYNC_FAILED, CONNECTOR_UP, CONNECTOR_DOWN
  integration_name    string              → name of the integration/connector
  remote_system       string              → external system name (e.g., "ServiceNow", "Jira", "Slack")
  remote_url          string              → endpoint URL (sanitized — no credentials)
  http_method         string              → HTTP method (GET, POST, PUT, etc.)
  http_status         int32               → HTTP response status code
  request_duration_ms int64               → round-trip time
  items_synced        int32               → number of items synchronized (for sync events)
  error_code          string              → error code
  error_detail        string              → error details
  retry_count         int32               → number of retries attempted
  message             string              → human-readable description
```

**IntegrationEventType enum**:
- INTEGRATION_EVENT_TYPE_UNSPECIFIED = 0
- INTEGRATION_EVENT_TYPE_API_CALL_SUCCESS = 1 (outbound API call succeeded)
- INTEGRATION_EVENT_TYPE_API_CALL_FAILURE = 2 (outbound API call failed)
- INTEGRATION_EVENT_TYPE_WEBHOOK_RECEIVED = 3 (inbound webhook received and processed)
- INTEGRATION_EVENT_TYPE_WEBHOOK_FAILED = 4 (inbound webhook processing failed)
- INTEGRATION_EVENT_TYPE_SYNC_START = 5 (data sync job started)
- INTEGRATION_EVENT_TYPE_SYNC_COMPLETE = 6 (data sync job completed)
- INTEGRATION_EVENT_TYPE_SYNC_FAILED = 7 (data sync job failed)
- INTEGRATION_EVENT_TYPE_CONNECTOR_UP = 8 (integration connector healthy)
- INTEGRATION_EVENT_TYPE_CONNECTOR_DOWN = 9 (integration connector unreachable)

---

## Category 6: Custom — No Proto

Consumer-defined events. By definition, the structure cannot be predefined. Consumers use the base `EventRecord` with the `attributes` map for custom key-value pairs.

---

## Category 7: NetworkEvent

Events from network infrastructure — device status changes, interface state transitions, routing protocol events, hardware component failures.

**Source models**: `NetworkDevice`, `Interface`, `BgpPeer`, `OspfNeighbor`, `LdpSession`, `SrPolicy`, `TeTunnel`, `VrfInstance`, `Chassis`, `Module`, `Fan`, `PowerSupply`

```
NetworkEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., Interface.status, BgpPeer.state)
  source_id           string              → device ID
  source_type         string              → "NetworkDevice"
  device_name         string              → sysName of the device
  device_ip           string              → management IP
  device_type         int32               → DeviceType enum value
  sub_category        NetworkEventType    → enum: DEVICE_STATUS, INTERFACE, BGP, OSPF, MPLS, LDP, SR, TE, VRF, QOS, HARDWARE
  component_id        string              → ID of the specific component (interface name, peer IP, etc.)
  component_name      string              → human-readable component name
  previous_state      string              → state before change (enum label or status string)
  current_state       string              → state after change
  message             string              → human-readable description
```

**NetworkEventType enum**:
- NETWORK_EVENT_TYPE_UNSPECIFIED = 0
- NETWORK_EVENT_TYPE_DEVICE_STATUS = 1 (device online/offline/warning/critical)
- NETWORK_EVENT_TYPE_INTERFACE = 2 (interface up/down/admin_down)
- NETWORK_EVENT_TYPE_BGP = 3 (peer state change)
- NETWORK_EVENT_TYPE_OSPF = 4 (neighbor state change)
- NETWORK_EVENT_TYPE_MPLS = 5 (label/FEC change)
- NETWORK_EVENT_TYPE_LDP = 6 (session state change)
- NETWORK_EVENT_TYPE_SR = 7 (policy/path status change)
- NETWORK_EVENT_TYPE_TE = 8 (tunnel/LSP status change)
- NETWORK_EVENT_TYPE_VRF = 9 (VRF status change)
- NETWORK_EVENT_TYPE_QOS = 10 (QoS policy violation)
- NETWORK_EVENT_TYPE_HARDWARE = 11 (chassis/module/fan/PSU status change)

---

## Category 8: KubernetesEvent

Events from Kubernetes clusters — pod status, node health, workload changes, resource conditions.

**Source models**: `K8sPod`, `K8sNode`, `K8sDeployment`, `K8sStatefulSet`, `K8sDaemonSet`, `K8sService`, `K8sNamespace`, `K8sNetworkPolicy`

```
KubernetesEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., K8sPod.status, K8sNode.ready)
  source_id           string              → resource name (pod name, node name)
  source_type         string              → "K8sPod", "K8sNode", "K8sDeployment", etc.
  cluster_id          string              → cluster identifier
  namespace           string              → Kubernetes namespace
  sub_category        KubernetesEventType → enum: POD, NODE, DEPLOYMENT, STATEFULSET, DAEMONSET, SERVICE, NAMESPACE, NETWORK_POLICY
  resource_name       string              → name of the K8s resource
  resource_kind       string              → K8s kind (Pod, Node, Deployment, etc.)
  previous_state      string              → state before change
  current_state       string              → state after change
  reason              string              → K8s event reason (e.g., "CrashLoopBackOff", "ImagePullBackOff")
  message             string              → human-readable description
  container_name      string              → container name (for pod-level events)
  ready_replicas      int32               → current ready count (for workload events)
  desired_replicas    int32               → desired count (for workload events)
```

**KubernetesEventType enum**:
- KUBERNETES_EVENT_TYPE_UNSPECIFIED = 0
- KUBERNETES_EVENT_TYPE_POD = 1 (pod status change, restart, crash)
- KUBERNETES_EVENT_TYPE_NODE = 2 (node ready/not-ready, cordon/drain)
- KUBERNETES_EVENT_TYPE_DEPLOYMENT = 3 (replica change, rollout)
- KUBERNETES_EVENT_TYPE_STATEFULSET = 4 (replica ordering, PVC binding)
- KUBERNETES_EVENT_TYPE_DAEMONSET = 5 (scheduling, node coverage)
- KUBERNETES_EVENT_TYPE_SERVICE = 6 (endpoint change, LB IP assignment)
- KUBERNETES_EVENT_TYPE_NAMESPACE = 7 (lifecycle change)
- KUBERNETES_EVENT_TYPE_NETWORK_POLICY = 8 (policy enforcement change)

---

## Category 9: PerformanceEvent

Threshold violations on any monitored metric — CPU, memory, temperature, traffic, fan speed, etc. These are generated when a metric crosses a defined threshold.

**Source models**: Any model with time series metrics — `PerformanceMetrics`, `Cpu`, `Memory`, `Fan`, `PowerSupply`, `InterfaceStatistics`

```
PerformanceEvent
  event_id            string              → base EventRecord
  property_id         string              → specific metric attribute (e.g., Cpu.utilization_percent)
  source_id           string              → entity ID (device, pod, VM)
  source_type         string              → entity type
  sub_category        PerformanceMetric   → enum: CPU, MEMORY, TEMPERATURE, TRAFFIC, DISK, FAN_SPEED, POWER_LOAD, VOLTAGE, LATENCY, PACKET_LOSS
  metric_name         string              → human-readable metric name (e.g., "CPU Utilization")
  metric_unit         string              → unit of measurement (e.g., "%", "°C", "RPM", "bps")
  current_value       double              → value that triggered the event
  threshold_value     double              → threshold that was crossed
  threshold_type      ThresholdType       → enum: UPPER, LOWER
  baseline_value      double              → normal/expected value (if available)
  duration_seconds    int64               → how long the metric has been above/below threshold
  component_id        string              → specific component (CPU ID, interface name, fan ID)
  component_name      string              → human-readable component name
  message             string              → human-readable description
```

**PerformanceMetric enum**:
- PERFORMANCE_METRIC_UNSPECIFIED = 0
- PERFORMANCE_METRIC_CPU = 1
- PERFORMANCE_METRIC_MEMORY = 2
- PERFORMANCE_METRIC_TEMPERATURE = 3
- PERFORMANCE_METRIC_TRAFFIC = 4
- PERFORMANCE_METRIC_DISK = 5
- PERFORMANCE_METRIC_FAN_SPEED = 6
- PERFORMANCE_METRIC_POWER_LOAD = 7
- PERFORMANCE_METRIC_VOLTAGE = 8
- PERFORMANCE_METRIC_LATENCY = 9
- PERFORMANCE_METRIC_PACKET_LOSS = 10

**ThresholdType enum**:
- THRESHOLD_TYPE_UNSPECIFIED = 0
- THRESHOLD_TYPE_UPPER = 1 (value exceeded maximum — e.g., CPU > 90%)
- THRESHOLD_TYPE_LOWER = 2 (value dropped below minimum — e.g., fan RPM < 1000)

---

## Category 10: SyslogEvent

Parsed syslog messages from network devices. The raw syslog text is parsed into structured fields with the `property_id` linking to the specific model attribute the message refers to.

**Source**: Syslog messages (RFC 5424) from network devices

```
SyslogEvent
  event_id            string              → base EventRecord
  property_id         string              → model attribute this syslog refers to (e.g., Interface.status for "%LINK-3-UPDOWN")
  source_id           string              → device ID
  source_type         string              → "NetworkDevice"
  device_name         string              → sysName of originating device
  device_ip           string              → IP that sent the syslog
  facility            int32               → syslog facility code (0-23)
  facility_name       string              → human-readable facility (e.g., "kern", "local7")
  syslog_severity     int32               → syslog severity (0-7, distinct from l8events Severity)
  syslog_severity_name string             → human-readable syslog severity (e.g., "emergency", "warning")
  mnemonic            string              → vendor mnemonic (e.g., "LINK-3-UPDOWN", "BGP-5-ADJCHANGE")
  process_name        string              → originating process (e.g., "bgpd", "ospfd", "lineproto")
  raw_message         string              → original syslog message text
  parsed_message      string              → cleaned/normalized message
  timestamp           int64               → syslog timestamp (device time, may differ from received_at)
```

---

## Category 11: TrapEvent

Parsed SNMP traps from network devices. The raw trap OID and varbinds are parsed into structured fields with the `property_id` linking to the specific model attribute.

**Source**: SNMP traps (v1/v2c/v3) from network devices

```
TrapEvent
  event_id            string              → base EventRecord
  property_id         string              → model attribute this trap refers to (e.g., Interface.status for linkDown)
  source_id           string              → device ID
  source_type         string              → "NetworkDevice"
  device_name         string              → sysName of originating device
  device_ip           string              → IP that sent the trap
  trap_oid            string              → trap OID (e.g., "1.3.6.1.6.3.1.1.5.3" for linkDown)
  trap_name           string              → human-readable trap name (e.g., "linkDown", "bgpEstablished")
  generic_trap        int32               → SNMPv1 generic trap number (0-6)
  specific_trap       int32               → SNMPv1 specific trap number
  enterprise_oid      string              → enterprise OID for vendor-specific traps
  snmp_version        string              → "v1", "v2c", or "v3"
  community           string              → SNMP community (v1/v2c) — may be masked for security
  varbinds            map<string, string>  → OID → value pairs from the trap
  uptime              int64               → sysUpTime from the trap
  message             string              → human-readable description
```

---

## Category 12: ComputeEvent

Events from hypervisors and virtual machines — VM status changes, host health, resource allocation changes.

**Source models**: `Hypervisor`, `VirtualMachine` (from probler's host/VM monitoring)

```
ComputeEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., VirtualMachine.status)
  source_id           string              → hypervisor or VM ID
  source_type         string              → "Hypervisor" or "VirtualMachine"
  sub_category        ComputeEventType    → enum: HYPERVISOR_STATUS, VM_STATUS, VM_MIGRATION, VM_RESOURCE, HOST_RESOURCE
  host_name           string              → hypervisor hostname
  host_ip             string              → hypervisor IP
  vm_name             string              → VM name (for VM events)
  vm_id               string              → VM identifier
  previous_state      string              → state before change
  current_state       string              → state after change
  cpu_count           int32               → allocated CPU count (for resource events)
  memory_mb           int64               → allocated memory in MB (for resource events)
  message             string              → human-readable description
```

**ComputeEventType enum**:
- COMPUTE_EVENT_TYPE_UNSPECIFIED = 0
- COMPUTE_EVENT_TYPE_HYPERVISOR_STATUS = 1 (host online/offline/maintenance)
- COMPUTE_EVENT_TYPE_VM_STATUS = 2 (VM power on/off/suspended/error)
- COMPUTE_EVENT_TYPE_VM_MIGRATION = 3 (vMotion, live migration)
- COMPUTE_EVENT_TYPE_VM_RESOURCE = 4 (CPU/memory allocation change)
- COMPUTE_EVENT_TYPE_HOST_RESOURCE = 5 (host capacity/utilization change)

---

## Category 13: StorageEvent

Events from storage systems — array health, volume status, capacity thresholds, replication state.

**Source models**: Storage arrays, volumes, LUNs (future probler models)

```
StorageEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute
  source_id           string              → storage array or volume ID
  source_type         string              → "StorageArray", "Volume", "LUN", "PersistentVolume"
  sub_category        StorageEventType    → enum: ARRAY_STATUS, VOLUME_STATUS, CAPACITY, REPLICATION, DISK, CONTROLLER
  array_name          string              → storage array name
  volume_name         string              → volume/LUN name (for volume events)
  previous_state      string              → state before change
  current_state       string              → state after change
  capacity_bytes      int64               → total capacity
  used_bytes          int64               → used capacity
  usage_percent       double              → usage percentage (for capacity events)
  message             string              → human-readable description
```

**StorageEventType enum**:
- STORAGE_EVENT_TYPE_UNSPECIFIED = 0
- STORAGE_EVENT_TYPE_ARRAY_STATUS = 1 (array online/degraded/offline)
- STORAGE_EVENT_TYPE_VOLUME_STATUS = 2 (volume online/degraded/offline)
- STORAGE_EVENT_TYPE_CAPACITY = 3 (capacity threshold crossed)
- STORAGE_EVENT_TYPE_REPLICATION = 4 (replication state change)
- STORAGE_EVENT_TYPE_DISK = 5 (physical disk failure/predictive)
- STORAGE_EVENT_TYPE_CONTROLLER = 6 (storage controller status)

---

## Category 14: PowerEvent

Events from power infrastructure — PDU status, UPS state, battery health, power consumption anomalies.

**Source models**: `PowerSupply` (existing in probler), PDU/UPS (future models)

```
PowerEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., PowerSupply.status, PowerSupply.load_percent)
  source_id           string              → PDU, UPS, or device ID
  source_type         string              → "PDU", "UPS", "PowerSupply"
  sub_category        PowerEventType      → enum: PSU_STATUS, PDU_STATUS, UPS_STATUS, BATTERY, LOAD, VOLTAGE, TEMPERATURE
  device_name         string              → parent device name (for PSU events on network devices)
  component_name      string              → PSU/PDU/UPS name
  previous_state      string              → state before change
  current_state       string              → state after change
  voltage             double              → current voltage reading
  current_amps        double              → current amperage reading
  load_percent        double              → current load percentage
  wattage             double              → current power draw in watts
  battery_percent     double              → battery charge percentage (UPS events)
  runtime_minutes     int32               → estimated runtime on battery (UPS events)
  message             string              → human-readable description
```

**PowerEventType enum**:
- POWER_EVENT_TYPE_UNSPECIFIED = 0
- POWER_EVENT_TYPE_PSU_STATUS = 1 (power supply OK/warning/failed)
- POWER_EVENT_TYPE_PDU_STATUS = 2 (PDU outlet on/off/overload)
- POWER_EVENT_TYPE_UPS_STATUS = 3 (UPS online/on-battery/bypass)
- POWER_EVENT_TYPE_BATTERY = 4 (battery low/replace/charging)
- POWER_EVENT_TYPE_LOAD = 5 (power load threshold crossed)
- POWER_EVENT_TYPE_VOLTAGE = 6 (voltage anomaly)
- POWER_EVENT_TYPE_TEMPERATURE = 7 (thermal event in power system)

---

## Category 15: GpuEvent

Events from GPU devices — GPU status, utilization, temperature, memory, error conditions.

**Source models**: `GpuDevice` (from probler's GPU monitoring)

```
GpuEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., GpuDevice.temperature, GpuDevice.utilization)
  source_id           string              → GPU device ID
  source_type         string              → "GpuDevice"
  sub_category        GpuEventType        → enum: GPU_STATUS, GPU_TEMPERATURE, GPU_MEMORY, GPU_UTILIZATION, GPU_ERROR, GPU_POWER
  device_name         string              → GPU device name
  host_name           string              → host machine name
  gpu_index           int32               → GPU index on the host
  gpu_model           string              → GPU model (e.g., "NVIDIA A100")
  previous_state      string              → state before change
  current_state       string              → state after change
  temperature_celsius double              → current temperature
  utilization_percent double              → current utilization
  memory_used_bytes   int64               → GPU memory used
  memory_total_bytes  int64               → GPU memory total
  power_draw_watts    double              → current power draw
  ecc_errors          int64               → ECC error count (for error events)
  message             string              → human-readable description
```

**GpuEventType enum**:
- GPU_EVENT_TYPE_UNSPECIFIED = 0
- GPU_EVENT_TYPE_STATUS = 1 (GPU online/offline/error)
- GPU_EVENT_TYPE_TEMPERATURE = 2 (thermal threshold)
- GPU_EVENT_TYPE_MEMORY = 3 (memory usage threshold)
- GPU_EVENT_TYPE_UTILIZATION = 4 (compute utilization threshold)
- GPU_EVENT_TYPE_ERROR = 5 (ECC errors, Xid errors, driver faults)
- GPU_EVENT_TYPE_POWER = 6 (power limit throttling)

---

## Category 16: TopologyEvent

Events from topology changes — link discovery, link loss, connectivity changes, neighbor relationships.

**Source models**: Topology graph derived from `NetworkDevice` interface/neighbor relationships

```
TopologyEvent
  event_id            string              → base EventRecord
  property_id         string              → specific attribute (e.g., Interface link)
  source_id           string              → device ID (one end of the link)
  source_type         string              → "NetworkDevice"
  sub_category        TopologyEventType   → enum: LINK_DISCOVERED, LINK_LOST, NEIGHBOR_CHANGE, TOPOLOGY_CHANGE
  local_device_id     string              → device on one end
  local_device_name   string              → device name
  local_interface     string              → interface name on local device
  remote_device_id    string              → device on other end
  remote_device_name  string              → remote device name
  remote_interface    string              → interface name on remote device
  discovery_protocol  string              → how the link was discovered (e.g., "LLDP", "CDP", "BGP", "OSPF")
  previous_state      string              → link state before change
  current_state       string              → link state after change
  message             string              → human-readable description
```

**TopologyEventType enum**:
- TOPOLOGY_EVENT_TYPE_UNSPECIFIED = 0
- TOPOLOGY_EVENT_TYPE_LINK_DISCOVERED = 1 (new link detected)
- TOPOLOGY_EVENT_TYPE_LINK_LOST = 2 (link no longer detected)
- TOPOLOGY_EVENT_TYPE_NEIGHBOR_CHANGE = 3 (neighbor relationship changed)
- TOPOLOGY_EVENT_TYPE_TOPOLOGY_CHANGE = 4 (STP/topology recalculation)

---

## Category 17: AutomationEvent

Events from automation workflows — rule execution, policy compliance, remediation actions.

**Source models**: Automation rules and workflows (future probler capability)

```
AutomationEvent
  event_id            string              → base EventRecord
  property_id         string              → model attribute the automation acted on
  source_id           string              → entity the automation targeted
  source_type         string              → entity type
  sub_category        AutomationEventType → enum: RULE_TRIGGERED, RULE_COMPLETED, RULE_FAILED, POLICY_VIOLATION, REMEDIATION
  rule_id             string              → automation rule ID
  rule_name           string              → human-readable rule name
  workflow_id         string              → workflow execution ID
  trigger_event_id    string              → event that triggered the automation
  action_taken        string              → what the automation did
  previous_state      string              → state before automation
  current_state       string              → state after automation
  success             bool                → whether the automation succeeded
  error_message       string              → error details (if failed)
  duration_ms         int64               → execution duration in milliseconds
  message             string              → human-readable description
```

**AutomationEventType enum**:
- AUTOMATION_EVENT_TYPE_UNSPECIFIED = 0
- AUTOMATION_EVENT_TYPE_RULE_TRIGGERED = 1 (rule condition met, execution started)
- AUTOMATION_EVENT_TYPE_RULE_COMPLETED = 2 (rule execution completed successfully)
- AUTOMATION_EVENT_TYPE_RULE_FAILED = 3 (rule execution failed)
- AUTOMATION_EVENT_TYPE_POLICY_VIOLATION = 4 (compliance policy violated)
- AUTOMATION_EVENT_TYPE_REMEDIATION = 5 (automated remediation action taken)

---

## Summary

| Category | Proto Message | Sub-category Enum | New Enums |
|----------|--------------|-------------------|-----------|
| 1 - Audit | `AuditEvent` | `AuditEventType` (10 values) | — |
| 2 - System | `SystemEvent` | `SystemEventType` (10 values) | — |
| 3 - Monitoring | `MonitoringEvent` | `MonitoringEventType` (9 values) | — |
| 4 - Security | `SecurityEvent` | `SecurityEventType` (10 values) | — |
| 5 - Integration | `IntegrationEvent` | `IntegrationEventType` (10 values) | — |
| 6 - Custom | — (uses base `EventRecord`) | — | — |
| 7 - Network | `NetworkEvent` | `NetworkEventType` (12 values) | — |
| 8 - Kubernetes | `KubernetesEvent` | `KubernetesEventType` (9 values) | — |
| 9 - Performance | `PerformanceEvent` | `PerformanceMetric` (11 values) | `ThresholdType` (3 values) |
| 10 - Syslog | `SyslogEvent` | — | — |
| 11 - Trap | `TrapEvent` | — | — |
| 12 - Compute | `ComputeEvent` | `ComputeEventType` (6 values) | — |
| 13 - Storage | `StorageEvent` | `StorageEventType` (7 values) | — |
| 14 - Power | `PowerEvent` | `PowerEventType` (8 values) | — |
| 15 - GPU | `GpuEvent` | `GpuEventType` (7 values) | — |
| 16 - Topology | `TopologyEvent` | `TopologyEventType` (5 values) | — |
| 17 - Automation | `AutomationEvent` | `AutomationEventType` (6 values) | — |

**Total**: 16 new proto messages, 15 new sub-category enums, 1 additional enum (`ThresholdType`), all in `proto/l8events.proto`
