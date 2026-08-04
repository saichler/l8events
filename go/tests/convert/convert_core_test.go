package convert_test

import (
	"github.com/saichler/l8events/go/convert"
	evt "github.com/saichler/l8types/go/types/l8events"
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

func TestConvert_NilRecord(t *testing.T) {
	c := convert.New()
	_, err := c.Convert(nil)
	if err == nil {
		t.Fatal("expected error for nil record")
	}
}

func TestConvert_Unspecified(t *testing.T) {
	c := convert.New()
	_, err := c.Convert(&evt.EventRecord{})
	if err == nil {
		t.Fatal("expected error for UNSPECIFIED category")
	}
}

func TestConvert_Custom(t *testing.T) {
	c := convert.New()
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
	c := convert.NewEmpty()
	r := makeRecord(evt.EventCategory_EVENT_CATEGORY_AUDIT, nil)
	_, err := c.Convert(r)
	if err == nil {
		t.Fatal("expected error for unregistered category")
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
	c := convert.New()
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
