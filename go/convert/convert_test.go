package convert

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

// makeRecord creates an EventRecord with the given category, common fields, and attributes.
func makeRecord(cat evt.EventCategory, attrs map[string]string) *evt.EventRecord {
	if attrs == nil {
		attrs = make(map[string]string)
	}
	attrs["propertyId"] = "prop-1"
	return &evt.EventRecord{
		EventId:    "evt-001",
		Category:   cat,
		SourceId:   "src-1",
		SourceType: "device",
		Message:    "test message",
		Attributes: attrs,
	}
}

func TestConvert_NilRecord(t *testing.T) {
	c := New()
	_, err := c.Convert(nil)
	if err == nil {
		t.Fatal("expected error for nil record")
	}
}

func TestConvert_Unspecified(t *testing.T) {
	c := New()
	_, err := c.Convert(&evt.EventRecord{})
	if err == nil {
		t.Fatal("expected error for UNSPECIFIED category")
	}
}

func TestConvert_Custom(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_CUSTOM, nil)
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != nil {
		t.Fatal("expected nil message for CUSTOM category")
	}
}

func TestConvert_UnregisteredCategory(t *testing.T) {
	c := &Converter{parsers: make(map[evt.EventCategory]Parser)}
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUDIT, nil)
	_, err := c.Convert(r)
	if err == nil {
		t.Fatal("expected error for unregistered category")
	}
}

// verifyCommon checks the common fields present on all category events.
func verifyCommon(t *testing.T, msg proto.Message, name string) {
	t.Helper()
	// Use protobuf reflection to check common fields
	refl := msg.ProtoReflect()
	fields := refl.Descriptor().Fields()

	check := func(fieldName protoreflect.Name, expected string) {
		fd := fields.ByName(fieldName)
		if fd == nil {
			t.Errorf("%s: field %s not found", name, fieldName)
			return
		}
		got := refl.Get(fd).String()
		if got != expected {
			t.Errorf("%s.%s = %q, want %q", name, fieldName, got, expected)
		}
	}

	check("event_id", "evt-001")
	check("property_id", "prop-1")
	check("source_id", "src-1")
	check("source_type", "device")
}

func TestConvert_AuditEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUDIT, map[string]string{
		"subCategory":   "2",
		"userId":        "user-1",
		"userName":      "admin",
		"userIp":        "10.0.0.1",
		"action":        "UPDATE",
		"serviceName":   "Employee",
		"serviceArea":   "10",
		"entityName":    "emp-001",
		"previousValue": "old",
		"newValue":      "new",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.AuditEvent)
	verifyCommon(t, msg, "AuditEvent")
	if e.SubCategory != 2 {
		t.Errorf("SubCategory = %d, want 2", e.SubCategory)
	}
	if e.UserId != "user-1" {
		t.Errorf("UserId = %q, want %q", e.UserId, "user-1")
	}
	if e.ServiceArea != 10 {
		t.Errorf("ServiceArea = %d, want 10", e.ServiceArea)
	}
	if e.Message != "test message" {
		t.Errorf("Message = %q, want %q", e.Message, "test message")
	}
}

func TestConvert_SystemEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_SYSTEM, map[string]string{
		"subCategory":   "1",
		"serviceName":   "vnet",
		"nodeId":        "node-1",
		"nodeIp":        "10.0.0.2",
		"previousState": "running",
		"currentState":  "stopped",
		"version":       "1.2.3",
		"errorCode":     "E001",
		"errorDetail":   "timeout",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.SystemEvent)
	verifyCommon(t, msg, "SystemEvent")
	if e.SubCategory != 1 {
		t.Errorf("SubCategory = %d, want 1", e.SubCategory)
	}
	if e.NodeId != "node-1" {
		t.Errorf("NodeId = %q, want %q", e.NodeId, "node-1")
	}
	if e.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", e.Version, "1.2.3")
	}
}

func TestConvert_MonitoringEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_MONITORING, map[string]string{
		"subCategory":     "3",
		"targetId":        "tgt-1",
		"targetName":      "switch-01",
		"targetType":      "switch",
		"protocol":        "SNMP",
		"pollDurationMs":  "1500",
		"itemsCollected":  "42",
		"errorCode":       "",
		"lastSuccessAt":   "1700000000",
		"staleDurationSec": "300",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.MonitoringEvent)
	verifyCommon(t, msg, "MonitoringEvent")
	if e.PollDurationMs != 1500 {
		t.Errorf("PollDurationMs = %d, want 1500", e.PollDurationMs)
	}
	if e.ItemsCollected != 42 {
		t.Errorf("ItemsCollected = %d, want 42", e.ItemsCollected)
	}
	if e.TargetId != "tgt-1" {
		t.Errorf("TargetId = %q, want %q", e.TargetId, "tgt-1")
	}
}

func TestConvert_SecurityEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_SECURITY, map[string]string{
		"subCategory":    "1",
		"userId":         "user-2",
		"userName":       "bob",
		"userIp":         "192.168.1.1",
		"targetResource": "/admin",
		"authMethod":     "password",
		"failureReason":  "bad password",
		"attemptCount":   "3",
		"certSubject":    "CN=test",
		"certExpiry":     "1800000000",
		"policyName":     "lockout",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.SecurityEvent)
	verifyCommon(t, msg, "SecurityEvent")
	if e.AttemptCount != 3 {
		t.Errorf("AttemptCount = %d, want 3", e.AttemptCount)
	}
	if e.CertExpiry != 1800000000 {
		t.Errorf("CertExpiry = %d, want 1800000000", e.CertExpiry)
	}
}

func TestConvert_IntegrationEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_INTEGRATION, map[string]string{
		"subCategory":       "2",
		"integrationName":   "salesforce",
		"remoteSystem":      "SFDC",
		"remoteUrl":         "https://sf.example.com",
		"httpMethod":        "POST",
		"httpStatus":        "201",
		"requestDurationMs": "450",
		"itemsSynced":       "10",
		"retryCount":        "0",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.IntegrationEvent)
	verifyCommon(t, msg, "IntegrationEvent")
	if e.HttpStatus != 201 {
		t.Errorf("HttpStatus = %d, want 201", e.HttpStatus)
	}
	if e.RequestDurationMs != 450 {
		t.Errorf("RequestDurationMs = %d, want 450", e.RequestDurationMs)
	}
}

func TestConvert_NetworkEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_NETWORK, map[string]string{
		"subCategory":   "1",
		"deviceName":    "core-sw-01",
		"deviceIp":      "10.1.1.1",
		"deviceType":    "2",
		"componentId":   "eth0",
		"componentName": "Ethernet0",
		"previousState": "up",
		"currentState":  "down",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.NetworkEvent)
	verifyCommon(t, msg, "NetworkEvent")
	if e.DeviceType != 2 {
		t.Errorf("DeviceType = %d, want 2", e.DeviceType)
	}
	if e.DeviceName != "core-sw-01" {
		t.Errorf("DeviceName = %q, want %q", e.DeviceName, "core-sw-01")
	}
	if e.CurrentState != "down" {
		t.Errorf("CurrentState = %q, want %q", e.CurrentState, "down")
	}
}

func TestConvert_KubernetesEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_KUBERNETES, map[string]string{
		"subCategory":     "3",
		"clusterId":       "cluster-1",
		"namespace":       "production",
		"resourceName":    "api-deploy",
		"resourceKind":    "Deployment",
		"previousState":   "Available",
		"currentState":    "Progressing",
		"reason":          "NewReplicaSet",
		"containerName":   "api",
		"readyReplicas":   "2",
		"desiredReplicas": "3",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.KubernetesEvent)
	verifyCommon(t, msg, "KubernetesEvent")
	if e.ReadyReplicas != 2 {
		t.Errorf("ReadyReplicas = %d, want 2", e.ReadyReplicas)
	}
	if e.DesiredReplicas != 3 {
		t.Errorf("DesiredReplicas = %d, want 3", e.DesiredReplicas)
	}
	if e.Namespace != "production" {
		t.Errorf("Namespace = %q, want %q", e.Namespace, "production")
	}
}

func TestConvert_PerformanceEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_PERFORMANCE, map[string]string{
		"subCategory":     "1",
		"metricName":      "cpu_util",
		"metricUnit":      "percent",
		"currentValue":    "95.5",
		"thresholdValue":  "90.0",
		"thresholdType":   "1",
		"baselineValue":   "60.0",
		"durationSeconds": "120",
		"componentId":     "cpu-0",
		"componentName":   "CPU 0",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.PerformanceEvent)
	verifyCommon(t, msg, "PerformanceEvent")
	if e.CurrentValue != 95.5 {
		t.Errorf("CurrentValue = %f, want 95.5", e.CurrentValue)
	}
	if e.ThresholdValue != 90.0 {
		t.Errorf("ThresholdValue = %f, want 90.0", e.ThresholdValue)
	}
	if e.ThresholdType != 1 {
		t.Errorf("ThresholdType = %d, want 1", e.ThresholdType)
	}
}

func TestConvert_SyslogEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_SYSLOG, map[string]string{
		"deviceName":         "router-01",
		"deviceIp":           "10.0.0.5",
		"facility":           "4",
		"facilityName":       "auth",
		"syslogSeverity":     "3",
		"syslogSeverityName": "error",
		"mnemonic":           "LINK-3-UPDOWN",
		"processName":        "sshd",
		"rawMessage":         "<36>router-01 sshd: login failure",
		"timestamp":          "1700000000",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.SyslogEvent)
	verifyCommon(t, msg, "SyslogEvent")
	if e.Facility != 4 {
		t.Errorf("Facility = %d, want 4", e.Facility)
	}
	if e.SyslogSeverity != 3 {
		t.Errorf("SyslogSeverity = %d, want 3", e.SyslogSeverity)
	}
	if e.Timestamp != 1700000000 {
		t.Errorf("Timestamp = %d, want 1700000000", e.Timestamp)
	}
	// SyslogEvent maps record.Message to ParsedMessage
	if e.ParsedMessage != "test message" {
		t.Errorf("ParsedMessage = %q, want %q", e.ParsedMessage, "test message")
	}
}

func TestConvert_TrapEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_TRAP, map[string]string{
		"deviceName":    "switch-01",
		"deviceIp":      "10.0.0.6",
		"trapOid":       "1.3.6.1.4.1.9.9.1",
		"trapName":      "linkDown",
		"genericTrap":   "2",
		"specificTrap":  "0",
		"enterpriseOid": "1.3.6.1.4.1.9",
		"snmpVersion":   "2c",
		"community":     "public",
		"uptime":        "86400",
		"varbinds.1.3.6.1.2.1.1.5.0": "switch-01",
		"varbinds.1.3.6.1.2.1.2.2.1.1": "1",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.TrapEvent)
	verifyCommon(t, msg, "TrapEvent")
	if e.GenericTrap != 2 {
		t.Errorf("GenericTrap = %d, want 2", e.GenericTrap)
	}
	if e.Uptime != 86400 {
		t.Errorf("Uptime = %d, want 86400", e.Uptime)
	}
	if len(e.Varbinds) != 2 {
		t.Errorf("Varbinds count = %d, want 2", len(e.Varbinds))
	}
	if e.Varbinds["1.3.6.1.2.1.1.5.0"] != "switch-01" {
		t.Errorf("Varbinds OID wrong: got %q", e.Varbinds["1.3.6.1.2.1.1.5.0"])
	}
}

func TestConvert_ComputeEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_COMPUTE, map[string]string{
		"subCategory":   "1",
		"hostName":      "esxi-01",
		"hostIp":        "10.0.0.10",
		"vmName":        "web-01",
		"vmId":          "vm-123",
		"previousState": "poweredOn",
		"currentState":  "suspended",
		"cpuCount":      "8",
		"memoryMb":      "16384",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.ComputeEvent)
	verifyCommon(t, msg, "ComputeEvent")
	if e.CpuCount != 8 {
		t.Errorf("CpuCount = %d, want 8", e.CpuCount)
	}
	if e.MemoryMb != 16384 {
		t.Errorf("MemoryMb = %d, want 16384", e.MemoryMb)
	}
}

func TestConvert_StorageEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_STORAGE, map[string]string{
		"subCategory":   "2",
		"arrayName":     "netapp-01",
		"volumeName":    "vol0",
		"previousState": "optimal",
		"currentState":  "degraded",
		"capacityBytes": "1099511627776",
		"usedBytes":     "879609302220",
		"usagePercent":  "80.0",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.StorageEvent)
	verifyCommon(t, msg, "StorageEvent")
	if e.CapacityBytes != 1099511627776 {
		t.Errorf("CapacityBytes = %d, want 1099511627776", e.CapacityBytes)
	}
	if e.UsagePercent != 80.0 {
		t.Errorf("UsagePercent = %f, want 80.0", e.UsagePercent)
	}
}

func TestConvert_PowerEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_POWER, map[string]string{
		"subCategory":    "1",
		"deviceName":     "ups-01",
		"componentName":  "battery-1",
		"previousState":  "normal",
		"currentState":   "on-battery",
		"voltage":        "220.5",
		"currentAmps":    "15.2",
		"loadPercent":    "72.3",
		"wattage":        "3345.6",
		"batteryPercent": "95.0",
		"runtimeMinutes": "45",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.PowerEvent)
	verifyCommon(t, msg, "PowerEvent")
	if e.Voltage != 220.5 {
		t.Errorf("Voltage = %f, want 220.5", e.Voltage)
	}
	if e.RuntimeMinutes != 45 {
		t.Errorf("RuntimeMinutes = %d, want 45", e.RuntimeMinutes)
	}
}

func TestConvert_GpuEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_GPU, map[string]string{
		"subCategory":        "1",
		"deviceName":         "gpu-node-01",
		"hostName":           "ml-server-01",
		"gpuIndex":           "0",
		"gpuModel":           "A100",
		"previousState":      "idle",
		"currentState":       "active",
		"temperatureCelsius": "72.5",
		"utilizationPercent": "98.0",
		"memoryUsedBytes":    "34359738368",
		"memoryTotalBytes":   "42949672960",
		"powerDrawWatts":     "350.0",
		"eccErrors":          "0",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.GpuEvent)
	verifyCommon(t, msg, "GpuEvent")
	if e.TemperatureCelsius != 72.5 {
		t.Errorf("TemperatureCelsius = %f, want 72.5", e.TemperatureCelsius)
	}
	if e.MemoryUsedBytes != 34359738368 {
		t.Errorf("MemoryUsedBytes = %d, want 34359738368", e.MemoryUsedBytes)
	}
	if e.GpuModel != "A100" {
		t.Errorf("GpuModel = %q, want %q", e.GpuModel, "A100")
	}
}

func TestConvert_TopologyEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_TOPOLOGY, map[string]string{
		"subCategory":       "1",
		"localDeviceId":     "sw-01",
		"localDeviceName":   "core-switch-01",
		"localInterface":    "Gi0/1",
		"remoteDeviceId":    "sw-02",
		"remoteDeviceName":  "access-switch-02",
		"remoteInterface":   "Gi0/24",
		"discoveryProtocol": "LLDP",
		"previousState":     "connected",
		"currentState":      "disconnected",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.TopologyEvent)
	verifyCommon(t, msg, "TopologyEvent")
	if e.LocalDeviceId != "sw-01" {
		t.Errorf("LocalDeviceId = %q, want %q", e.LocalDeviceId, "sw-01")
	}
	if e.DiscoveryProtocol != "LLDP" {
		t.Errorf("DiscoveryProtocol = %q, want %q", e.DiscoveryProtocol, "LLDP")
	}
}

func TestConvert_AutomationEvent(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUTOMATION, map[string]string{
		"subCategory":    "2",
		"ruleId":         "rule-1",
		"ruleName":       "auto-remediate",
		"workflowId":     "wf-001",
		"triggerEventId": "evt-000",
		"actionTaken":    "restart-service",
		"previousState":  "failed",
		"currentState":   "running",
		"success":        "true",
		"errorMessage":   "",
		"durationMs":     "2500",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.AutomationEvent)
	verifyCommon(t, msg, "AutomationEvent")
	if !e.Success {
		t.Error("Success = false, want true")
	}
	if e.DurationMs != 2500 {
		t.Errorf("DurationMs = %d, want 2500", e.DurationMs)
	}
	if e.RuleName != "auto-remediate" {
		t.Errorf("RuleName = %q, want %q", e.RuleName, "auto-remediate")
	}
}

// Edge case: bad numeric string returns error
func TestConvert_BadNumericAttribute(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_MONITORING, map[string]string{
		"pollDurationMs": "not-a-number",
	})
	_, err := c.Convert(r)
	if err == nil {
		t.Fatal("expected error for malformed numeric attribute")
	}
}

// Edge case: bad boolean string returns error
func TestConvert_BadBooleanAttribute(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUTOMATION, map[string]string{
		"success": "not-a-bool",
	})
	_, err := c.Convert(r)
	if err == nil {
		t.Fatal("expected error for malformed boolean attribute")
	}
}

// Edge case: missing attributes yield zero values (no error)
func TestConvert_MissingAttributes(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_NETWORK, map[string]string{})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.NetworkEvent)
	if e.DeviceName != "" {
		t.Errorf("DeviceName = %q, want empty", e.DeviceName)
	}
	if e.DeviceType != 0 {
		t.Errorf("DeviceType = %d, want 0", e.DeviceType)
	}
}

// Edge case: TrapEvent with no varbinds prefix
func TestConvert_TrapEvent_NoVarbinds(t *testing.T) {
	c := New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_TRAP, map[string]string{
		"trapOid": "1.3.6.1",
	})
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := msg.(*evt.TrapEvent)
	if e.Varbinds != nil {
		t.Errorf("Varbinds = %v, want nil", e.Varbinds)
	}
}

// Custom parser registration
type mockParser struct {
	called bool
}

func (m *mockParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	m.called = true
	return &evt.AuditEvent{EventId: "custom"}, nil
}

func TestConvert_CustomParserRegistration(t *testing.T) {
	c := New()
	mp := &mockParser{}
	c.Register(evt.EventCategory_EVENT_CATEGORY_AUDIT, mp)
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUDIT, nil)
	msg, err := c.Convert(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mp.called {
		t.Fatal("custom parser was not called")
	}
	e := msg.(*evt.AuditEvent)
	if e.EventId != "custom" {
		t.Errorf("EventId = %q, want %q", e.EventId, "custom")
	}
}
