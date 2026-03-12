package convert

import evt "github.com/saichler/l8events/go/types/l8events"

// registerBuiltins wires all 16 built-in category parsers into the converter.
func registerBuiltins(c *Converter) {
	c.Register(evt.EventCategory_EVENT_CATEGORY_AUDIT, &auditParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_SYSTEM, &systemParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_MONITORING, &monitoringParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_SECURITY, &securityParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_INTEGRATION, &integrationParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_NETWORK, &networkParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_KUBERNETES, &kubernetesParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_PERFORMANCE, &performanceParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_SYSLOG, &syslogParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_TRAP, &trapParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_COMPUTE, &computeParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_STORAGE, &storageParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_POWER, &powerParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_GPU, &gpuParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_TOPOLOGY, &topologyParser{})
	c.Register(evt.EventCategory_EVENT_CATEGORY_AUTOMATION, &automationParser{})
}
