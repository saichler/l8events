package convert_test

import (
	"github.com/saichler/l8events/go/convert"
	evt "github.com/saichler/l8types/go/types/l8events"
	"testing"
)

func TestConvert_AuditEvent(t *testing.T) {
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_MONITORING, map[string]string{
		"subCategory":      "3",
		"targetId":         "tgt-1",
		"targetName":       "switch-01",
		"targetType":       "switch",
		"protocol":         "SNMP",
		"pollDurationMs":   "1500",
		"itemsCollected":   "42",
		"errorCode":        "",
		"lastSuccessAt":    "1700000000",
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
	c := convert.New()
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
	c := convert.New()
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

func TestConvert_PerformanceEvent(t *testing.T) {
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_TRAP, map[string]string{
		"deviceName":                   "switch-01",
		"deviceIp":                     "10.0.0.6",
		"trapOid":                      "1.3.6.1.4.1.9.9.1",
		"trapName":                     "linkDown",
		"genericTrap":                  "2",
		"specificTrap":                 "0",
		"enterpriseOid":                "1.3.6.1.4.1.9",
		"snmpVersion":                  "2c",
		"community":                    "public",
		"uptime":                       "86400",
		"varbinds.1.3.6.1.2.1.1.5.0":   "switch-01",
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

func TestConvert_AutomationEvent(t *testing.T) {
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUTOMATION, map[string]string{
		"success": "not-a-bool",
	})
	_, err := c.Convert(r)
	if err == nil {
		t.Fatal("expected error for malformed boolean attribute")
	}
}

// Edge case: TrapEvent with no varbinds prefix
func TestConvert_TrapEvent_NoVarbinds(t *testing.T) {
	c := convert.New()
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
