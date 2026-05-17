/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8events"
)

var _ ifs.IEvents = (*Events)(nil)

// Events implements ifs.IEvents, routing events to the Events service via VNic unicast.
type Events struct {
	vnic ifs.IVNic
}

func NewEvents(vnic ifs.IVNic) *Events {
	return &Events{vnic: vnic}
}

func (e *Events) PostEvent(category l8events.EventCategory, eventType string,
	severity l8events.Severity, sourceId, sourceName, sourceType, message string,
	attributes map[string]string) {

	event := &l8events.EventRecord{
		Category:   category,
		EventType:  eventType,
		Severity:   severity,
		SourceId:   sourceId,
		SourceName: sourceName,
		SourceType: sourceType,
		Message:    message,
		OccurredAt: time.Now().Unix(),
		Attributes: attributes,
	}
	e.post(event)
}

func (e *Events) PostAuditEvent(evt *l8events.AuditEvent) {
	e.post(evt)
}

func (e *Events) PostSystemEvent(evt *l8events.SystemEvent) {
	e.post(evt)
}

func (e *Events) PostMonitoringEvent(evt *l8events.MonitoringEvent) {
	e.post(evt)
}

func (e *Events) PostSecurityEvent(evt *l8events.SecurityEvent) {
	e.post(evt)
}

func (e *Events) PostIntegrationEvent(evt *l8events.IntegrationEvent) {
	e.post(evt)
}

func (e *Events) PostNetworkEvent(evt *l8events.NetworkEvent) {
	e.post(evt)
}

func (e *Events) PostKubernetesEvent(evt *l8events.KubernetesEvent) {
	e.post(evt)
}

func (e *Events) PostPerformanceEvent(evt *l8events.PerformanceEvent) {
	e.post(evt)
}

func (e *Events) PostSyslogEvent(evt *l8events.SyslogEvent) {
	e.post(evt)
}

func (e *Events) PostTrapEvent(evt *l8events.TrapEvent) {
	e.post(evt)
}

func (e *Events) PostComputeEvent(evt *l8events.ComputeEvent) {
	e.post(evt)
}

func (e *Events) PostStorageEvent(evt *l8events.StorageEvent) {
	e.post(evt)
}

func (e *Events) PostPowerEvent(evt *l8events.PowerEvent) {
	e.post(evt)
}

func (e *Events) PostGpuEvent(evt *l8events.GpuEvent) {
	e.post(evt)
}

func (e *Events) PostTopologyEvent(evt *l8events.TopologyEvent) {
	e.post(evt)
}

func (e *Events) PostAutomationEvent(evt *l8events.AutomationEvent) {
	e.post(evt)
}

func (e *Events) post(payload interface{}) {
	err := e.vnic.Unicast("", EventsServiceName, EventsServiceArea, ifs.POST, payload)
	if err != nil {
		e.vnic.Resources().Logger().Warning("PostEvent: " + err.Error())
	}
}
