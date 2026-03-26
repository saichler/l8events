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
	"fmt"
	"time"

	evt "github.com/saichler/l8events/go/types/l8events"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

// PostEvent creates and posts an EventRecord to the Events service.
// Returns silently if the Events service is not activated.
func PostEvent(vnic ifs.IVNic, category evt.EventCategory, eventType string,
	severity evt.Severity, sourceId, sourceName, sourceType, message string,
	attributes map[string]string) {

	h, ok := vnic.Resources().Services().ServiceHandler(EventsServiceName, EventsServiceArea)
	if !ok {
		vnic.Resources().Logger().Warning("PostEvent: Events service not activated, event discarded")
		return
	}

	event := &evt.EventRecord{
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

	resp := h.Post(object.New(nil, event), vnic)
	if resp.Error() != nil {
		vnic.Resources().Logger().Error(fmt.Sprintf("PostEvent failed: %v", resp.Error()))
	}
}

// PostAuditEvent is a convenience for audit trail events (category=AUDIT).
func PostAuditEvent(vnic ifs.IVNic, eventType string, severity evt.Severity,
	userId, action, target, message string) {

	attrs := map[string]string{
		"userId": userId,
		"action": action,
		"target": target,
	}
	PostEvent(vnic, evt.EventCategory_EVENT_CATEGORY_AUDIT, eventType,
		severity, userId, userId, "user", message, attrs)
}

// PostSecurityEvent is a convenience for security events (category=SECURITY).
func PostSecurityEvent(vnic ifs.IVNic, eventType string, severity evt.Severity,
	userId, action, message string) {

	attrs := map[string]string{
		"userId": userId,
		"action": action,
	}
	PostEvent(vnic, evt.EventCategory_EVENT_CATEGORY_SECURITY, eventType,
		severity, userId, userId, "user", message, attrs)
}
