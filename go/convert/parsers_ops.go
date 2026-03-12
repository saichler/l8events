package convert

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"google.golang.org/protobuf/proto"
)

// auditParser converts EventRecord → AuditEvent.
type auditParser struct{}

func (p *auditParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	sa, err := i32(a, "serviceArea")
	if err != nil {
		return nil, err
	}
	e := &evt.AuditEvent{SubCategory: evt.AuditEventType(sc), ServiceArea: sa}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.UserId = str(a, "userId")
	e.UserName = str(a, "userName")
	e.UserIp = str(a, "userIp")
	e.Action = str(a, "action")
	e.ServiceName = str(a, "serviceName")
	e.EntityName = str(a, "entityName")
	e.PreviousValue = str(a, "previousValue")
	e.NewValue = str(a, "newValue")
	return e, nil
}

// systemParser converts EventRecord → SystemEvent.
type systemParser struct{}

func (p *systemParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	e := &evt.SystemEvent{SubCategory: evt.SystemEventType(sc)}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.ServiceName = str(a, "serviceName")
	e.NodeId = str(a, "nodeId")
	e.NodeIp = str(a, "nodeIp")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	e.Version = str(a, "version")
	e.ErrorCode = str(a, "errorCode")
	e.ErrorDetail = str(a, "errorDetail")
	return e, nil
}

// monitoringParser converts EventRecord → MonitoringEvent.
type monitoringParser struct{}

func (p *monitoringParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	pollMs, err := i64(a, "pollDurationMs")
	if err != nil {
		return nil, err
	}
	items, err := i32(a, "itemsCollected")
	if err != nil {
		return nil, err
	}
	lastSuccess, err := i64(a, "lastSuccessAt")
	if err != nil {
		return nil, err
	}
	staleSec, err := i64(a, "staleDurationSec")
	if err != nil {
		return nil, err
	}
	e := &evt.MonitoringEvent{
		SubCategory:      evt.MonitoringEventType(sc),
		PollDurationMs:   pollMs,
		ItemsCollected:   items,
		LastSuccessAt:    lastSuccess,
		StaleDurationSec: staleSec,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.TargetId = str(a, "targetId")
	e.TargetName = str(a, "targetName")
	e.TargetType = str(a, "targetType")
	e.Protocol = str(a, "protocol")
	e.ErrorCode = str(a, "errorCode")
	e.ErrorDetail = str(a, "errorDetail")
	return e, nil
}

// securityParser converts EventRecord → SecurityEvent.
type securityParser struct{}

func (p *securityParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	attempts, err := i32(a, "attemptCount")
	if err != nil {
		return nil, err
	}
	certExp, err := i64(a, "certExpiry")
	if err != nil {
		return nil, err
	}
	e := &evt.SecurityEvent{
		SubCategory:  evt.SecurityEventType(sc),
		AttemptCount: attempts,
		CertExpiry:   certExp,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.UserId = str(a, "userId")
	e.UserName = str(a, "userName")
	e.UserIp = str(a, "userIp")
	e.TargetResource = str(a, "targetResource")
	e.AuthMethod = str(a, "authMethod")
	e.FailureReason = str(a, "failureReason")
	e.CertSubject = str(a, "certSubject")
	e.PolicyName = str(a, "policyName")
	return e, nil
}

// integrationParser converts EventRecord → IntegrationEvent.
type integrationParser struct{}

func (p *integrationParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	httpStatus, err := i32(a, "httpStatus")
	if err != nil {
		return nil, err
	}
	reqDur, err := i64(a, "requestDurationMs")
	if err != nil {
		return nil, err
	}
	itemsSynced, err := i32(a, "itemsSynced")
	if err != nil {
		return nil, err
	}
	retryCount, err := i32(a, "retryCount")
	if err != nil {
		return nil, err
	}
	e := &evt.IntegrationEvent{
		SubCategory:       evt.IntegrationEventType(sc),
		HttpStatus:        httpStatus,
		RequestDurationMs: reqDur,
		ItemsSynced:       itemsSynced,
		RetryCount:        retryCount,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.IntegrationName = str(a, "integrationName")
	e.RemoteSystem = str(a, "remoteSystem")
	e.RemoteUrl = str(a, "remoteUrl")
	e.HttpMethod = str(a, "httpMethod")
	e.ErrorCode = str(a, "errorCode")
	e.ErrorDetail = str(a, "errorDetail")
	return e, nil
}

// performanceParser converts EventRecord → PerformanceEvent.
type performanceParser struct{}

func (p *performanceParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	curVal, err := f64(a, "currentValue")
	if err != nil {
		return nil, err
	}
	threshVal, err := f64(a, "thresholdValue")
	if err != nil {
		return nil, err
	}
	threshType, err := i32(a, "thresholdType")
	if err != nil {
		return nil, err
	}
	baseline, err := f64(a, "baselineValue")
	if err != nil {
		return nil, err
	}
	durSec, err := i64(a, "durationSeconds")
	if err != nil {
		return nil, err
	}
	e := &evt.PerformanceEvent{
		SubCategory:     evt.PerformanceMetric(sc),
		CurrentValue:    curVal,
		ThresholdValue:  threshVal,
		ThresholdType:   evt.ThresholdType(threshType),
		BaselineValue:   baseline,
		DurationSeconds: durSec,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.MetricName = str(a, "metricName")
	e.MetricUnit = str(a, "metricUnit")
	e.ComponentId = str(a, "componentId")
	e.ComponentName = str(a, "componentName")
	return e, nil
}

// syslogParser converts EventRecord → SyslogEvent. No SubCategory field.
type syslogParser struct{}

func (p *syslogParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	facility, err := i32(a, "facility")
	if err != nil {
		return nil, err
	}
	syslogSev, err := i32(a, "syslogSeverity")
	if err != nil {
		return nil, err
	}
	ts, err := i64(a, "timestamp")
	if err != nil {
		return nil, err
	}
	e := &evt.SyslogEvent{
		Facility:       facility,
		SyslogSeverity: syslogSev,
		Timestamp:      ts,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.ParsedMessage)
	e.DeviceName = str(a, "deviceName")
	e.DeviceIp = str(a, "deviceIp")
	e.FacilityName = str(a, "facilityName")
	e.SyslogSeverityName = str(a, "syslogSeverityName")
	e.Mnemonic = str(a, "mnemonic")
	e.ProcessName = str(a, "processName")
	e.RawMessage = str(a, "rawMessage")
	return e, nil
}

// trapParser converts EventRecord → TrapEvent.
type trapParser struct{}

func (p *trapParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	genTrap, err := i32(a, "genericTrap")
	if err != nil {
		return nil, err
	}
	specTrap, err := i32(a, "specificTrap")
	if err != nil {
		return nil, err
	}
	uptime, err := i64(a, "uptime")
	if err != nil {
		return nil, err
	}
	e := &evt.TrapEvent{
		GenericTrap:  genTrap,
		SpecificTrap: specTrap,
		Uptime:       uptime,
		Varbinds:     extractMap(a, "varbinds."),
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.DeviceName = str(a, "deviceName")
	e.DeviceIp = str(a, "deviceIp")
	e.TrapOid = str(a, "trapOid")
	e.TrapName = str(a, "trapName")
	e.EnterpriseOid = str(a, "enterpriseOid")
	e.SnmpVersion = str(a, "snmpVersion")
	e.Community = str(a, "community")
	return e, nil
}

// automationParser converts EventRecord → AutomationEvent.
type automationParser struct{}

func (p *automationParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	success, err := boolean(a, "success")
	if err != nil {
		return nil, err
	}
	durMs, err := i64(a, "durationMs")
	if err != nil {
		return nil, err
	}
	e := &evt.AutomationEvent{
		SubCategory: evt.AutomationEventType(sc),
		Success:     success,
		DurationMs:  durMs,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.RuleId = str(a, "ruleId")
	e.RuleName = str(a, "ruleName")
	e.WorkflowId = str(a, "workflowId")
	e.TriggerEventId = str(a, "triggerEventId")
	e.ActionTaken = str(a, "actionTaken")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	e.ErrorMessage = str(a, "errorMessage")
	return e, nil
}
