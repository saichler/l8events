package convert_test

import (
	"github.com/saichler/l8events/go/convert"
	evt "github.com/saichler/l8types/go/types/l8events"
	"testing"
)

func TestConvert_NetworkEvent(t *testing.T) {
	c := convert.New()
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
	c := convert.New()
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

func TestConvert_ComputeEvent(t *testing.T) {
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
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
	c := convert.New()
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

// Edge case: missing attributes yield zero values (no error)
func TestConvert_MissingAttributes(t *testing.T) {
	c := convert.New()
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
