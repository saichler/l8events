package convert

import (
	evt "github.com/saichler/l8events/go/types/l8events"
	"google.golang.org/protobuf/proto"
)

// networkParser converts EventRecord → NetworkEvent.
type networkParser struct{}

func (p *networkParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	devType, err := i32(a, "deviceType")
	if err != nil {
		return nil, err
	}
	e := &evt.NetworkEvent{
		SubCategory: evt.NetworkEventType(sc),
		DeviceType:  devType,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.DeviceName = str(a, "deviceName")
	e.DeviceIp = str(a, "deviceIp")
	e.ComponentId = str(a, "componentId")
	e.ComponentName = str(a, "componentName")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}

// kubernetesParser converts EventRecord → KubernetesEvent.
type kubernetesParser struct{}

func (p *kubernetesParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	readyRep, err := i32(a, "readyReplicas")
	if err != nil {
		return nil, err
	}
	desiredRep, err := i32(a, "desiredReplicas")
	if err != nil {
		return nil, err
	}
	e := &evt.KubernetesEvent{
		SubCategory:     evt.KubernetesEventType(sc),
		ReadyReplicas:   readyRep,
		DesiredReplicas: desiredRep,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.ClusterId = str(a, "clusterId")
	e.Namespace = str(a, "namespace")
	e.ResourceName = str(a, "resourceName")
	e.ResourceKind = str(a, "resourceKind")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	e.Reason = str(a, "reason")
	e.ContainerName = str(a, "containerName")
	return e, nil
}

// computeParser converts EventRecord → ComputeEvent.
type computeParser struct{}

func (p *computeParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	cpuCount, err := i32(a, "cpuCount")
	if err != nil {
		return nil, err
	}
	memMb, err := i64(a, "memoryMb")
	if err != nil {
		return nil, err
	}
	e := &evt.ComputeEvent{
		SubCategory: evt.ComputeEventType(sc),
		CpuCount:    cpuCount,
		MemoryMb:    memMb,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.HostName = str(a, "hostName")
	e.HostIp = str(a, "hostIp")
	e.VmName = str(a, "vmName")
	e.VmId = str(a, "vmId")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}

// storageParser converts EventRecord → StorageEvent.
type storageParser struct{}

func (p *storageParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	capBytes, err := i64(a, "capacityBytes")
	if err != nil {
		return nil, err
	}
	usedBytes, err := i64(a, "usedBytes")
	if err != nil {
		return nil, err
	}
	usagePct, err := f64(a, "usagePercent")
	if err != nil {
		return nil, err
	}
	e := &evt.StorageEvent{
		SubCategory:   evt.StorageEventType(sc),
		CapacityBytes: capBytes,
		UsedBytes:     usedBytes,
		UsagePercent:  usagePct,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.ArrayName = str(a, "arrayName")
	e.VolumeName = str(a, "volumeName")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}

// powerParser converts EventRecord → PowerEvent.
type powerParser struct{}

func (p *powerParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	voltage, err := f64(a, "voltage")
	if err != nil {
		return nil, err
	}
	amps, err := f64(a, "currentAmps")
	if err != nil {
		return nil, err
	}
	loadPct, err := f64(a, "loadPercent")
	if err != nil {
		return nil, err
	}
	wattage, err := f64(a, "wattage")
	if err != nil {
		return nil, err
	}
	battPct, err := f64(a, "batteryPercent")
	if err != nil {
		return nil, err
	}
	runtime, err := i32(a, "runtimeMinutes")
	if err != nil {
		return nil, err
	}
	e := &evt.PowerEvent{
		SubCategory:    evt.PowerEventType(sc),
		Voltage:        voltage,
		CurrentAmps:    amps,
		LoadPercent:    loadPct,
		Wattage:        wattage,
		BatteryPercent: battPct,
		RuntimeMinutes: runtime,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.DeviceName = str(a, "deviceName")
	e.ComponentName = str(a, "componentName")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}

// gpuParser converts EventRecord → GpuEvent.
type gpuParser struct{}

func (p *gpuParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	gpuIdx, err := i32(a, "gpuIndex")
	if err != nil {
		return nil, err
	}
	temp, err := f64(a, "temperatureCelsius")
	if err != nil {
		return nil, err
	}
	util, err := f64(a, "utilizationPercent")
	if err != nil {
		return nil, err
	}
	memUsed, err := i64(a, "memoryUsedBytes")
	if err != nil {
		return nil, err
	}
	memTotal, err := i64(a, "memoryTotalBytes")
	if err != nil {
		return nil, err
	}
	powerDraw, err := f64(a, "powerDrawWatts")
	if err != nil {
		return nil, err
	}
	eccErrs, err := i64(a, "eccErrors")
	if err != nil {
		return nil, err
	}
	e := &evt.GpuEvent{
		SubCategory:        evt.GpuEventType(sc),
		GpuIndex:           gpuIdx,
		TemperatureCelsius: temp,
		UtilizationPercent: util,
		MemoryUsedBytes:    memUsed,
		MemoryTotalBytes:   memTotal,
		PowerDrawWatts:     powerDraw,
		EccErrors:          eccErrs,
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.DeviceName = str(a, "deviceName")
	e.HostName = str(a, "hostName")
	e.GpuModel = str(a, "gpuModel")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}

// topologyParser converts EventRecord → TopologyEvent.
type topologyParser struct{}

func (p *topologyParser) Parse(r *evt.EventRecord) (proto.Message, error) {
	a := r.Attributes
	sc, err := subCategory(a)
	if err != nil {
		return nil, err
	}
	e := &evt.TopologyEvent{
		SubCategory: evt.TopologyEventType(sc),
	}
	setCommon(r, &e.EventId, &e.PropertyId, &e.SourceId, &e.SourceType, &e.Message)
	e.LocalDeviceId = str(a, "localDeviceId")
	e.LocalDeviceName = str(a, "localDeviceName")
	e.LocalInterface = str(a, "localInterface")
	e.RemoteDeviceId = str(a, "remoteDeviceId")
	e.RemoteDeviceName = str(a, "remoteDeviceName")
	e.RemoteInterface = str(a, "remoteInterface")
	e.DiscoveryProtocol = str(a, "discoveryProtocol")
	e.PreviousState = str(a, "previousState")
	e.CurrentState = str(a, "currentState")
	return e, nil
}
